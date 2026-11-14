package database

import (
	"bufio"
	"context"
	"crypto/sha256"
	"database/sql"
	"embed"
	"encoding/hex"
	"fmt"
	"io/fs"
	"log/slog"
	"regexp"
	"sort"
	"strings"
)

// Las migraciones viven dentro del binario. Reglas:
//   - Nombre: NNNN_descripcion.sql (se aplican en orden y una sola vez).
//   - Sentencias terminadas en ";" al final de la línea. Para triggers y
//     procedimientos se usa SET TERM, igual que en isql.
//   - Firebird no deja usar en una transacción una tabla creada en ella:
//     los datos de una tabla nueva van en la siguiente migración.
//   - Una migración aplicada NO se edita: se crea otra nueva. Si el archivo
//     cambia, la API se niega a arrancar (se compara su SHA-256).
//
//go:embed migrations/*.sql
var migrationsFS embed.FS

var (
	nombreMigracion = regexp.MustCompile(`^(\d{4})_[a-z0-9_]+\.sql$`)
	setTerm         = regexp.MustCompile(`(?i)^SET\s+TERM\s+(\S+)\s+(\S+)$`)
)

type migracion struct {
	version, nombre, checksum string
	sentencias                []string
}

// Migrate aplica las migraciones pendientes, cada una en su propia transacción.
func Migrate(ctx context.Context, db *sql.DB, log *slog.Logger) error {
	if err := asegurarTabla(ctx, db); err != nil {
		return err
	}
	pendientes, err := cargar()
	if err != nil {
		return err
	}
	aplicadas, err := leerAplicadas(ctx, db)
	if err != nil {
		return err
	}

	// Las ya aplicadas solo se verifican; las nuevas se ejecutan en orden.
	for _, m := range pendientes {
		if sum, ok := aplicadas[m.version]; ok {
			if sum != m.checksum {
				return fmt.Errorf("la migración %s ya aplicada fue modificada; crea una nueva en lugar de editarla", m.nombre)
			}
			continue
		}
		if err := aplicar(ctx, db, m); err != nil {
			return fmt.Errorf("migración %s: %w", m.nombre, err)
		}
		log.Info("migración aplicada", "migracion", m.nombre)
	}
	return nil
}

// asegurarTabla crea SCHEMA_MIGRATIONS la primera vez (Firebird no tiene
// CREATE TABLE IF NOT EXISTS, por eso se consulta el catálogo).
func asegurarTabla(ctx context.Context, db *sql.DB) error {
	var n int
	err := db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM RDB$RELATIONS WHERE RDB$RELATION_NAME = 'SCHEMA_MIGRATIONS'`).Scan(&n)
	if err != nil || n > 0 {
		return err
	}
	_, err = db.ExecContext(ctx, `CREATE TABLE SCHEMA_MIGRATIONS (
		VERSION     VARCHAR(4)   NOT NULL PRIMARY KEY,
		NOMBRE      VARCHAR(200) NOT NULL,
		CHECKSUM    CHAR(64)     NOT NULL,
		APLICADA_EN TIMESTAMP    DEFAULT LOCALTIMESTAMP NOT NULL
	)`)
	return err
}

// leerAplicadas devuelve versión → checksum de lo ya aplicado. Firebird
// rellena CHAR con espacios, por eso el TrimSpace.
func leerAplicadas(ctx context.Context, db *sql.DB) (map[string]string, error) {
	rows, err := db.QueryContext(ctx, `SELECT VERSION, CHECKSUM FROM SCHEMA_MIGRATIONS`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := map[string]string{}
	for rows.Next() {
		var v, sum string
		if err := rows.Scan(&v, &sum); err != nil {
			return nil, err
		}
		out[strings.TrimSpace(v)] = strings.TrimSpace(sum)
	}
	return out, rows.Err()
}

// aplicar ejecuta una migración completa o nada: si una sentencia falla, se
// revierte todo y no se registra.
func aplicar(ctx context.Context, db *sql.DB, m migracion) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }() // no hace nada si ya hubo Commit

	for _, s := range m.sentencias {
		if _, err := tx.ExecContext(ctx, s); err != nil {
			return fmt.Errorf("%w\nSQL: %s", err, s)
		}
	}
	if _, err := tx.ExecContext(ctx,
		`INSERT INTO SCHEMA_MIGRATIONS (VERSION, NOMBRE, CHECKSUM) VALUES (?, ?, ?)`,
		m.version, m.nombre, m.checksum); err != nil {
		return err
	}
	return tx.Commit()
}

// cargar lee los .sql embebidos, valida nombres y versiones únicas, y calcula
// el checksum de cada archivo.
func cargar() ([]migracion, error) {
	archivos, err := fs.Glob(migrationsFS, "migrations/*.sql")
	if err != nil {
		return nil, err
	}
	sort.Strings(archivos)

	var out []migracion
	vistas := map[string]string{}
	for _, ruta := range archivos {
		nombre := strings.TrimPrefix(ruta, "migrations/")
		match := nombreMigracion.FindStringSubmatch(nombre)
		if match == nil {
			return nil, fmt.Errorf("nombre de migración inválido %q (usa NNNN_descripcion.sql)", nombre)
		}
		if previa, ok := vistas[match[1]]; ok {
			return nil, fmt.Errorf("versión repetida %s: %s y %s", match[1], previa, nombre)
		}
		vistas[match[1]] = nombre

		contenido, err := migrationsFS.ReadFile(ruta)
		if err != nil {
			return nil, err
		}
		sentencias, err := dividir(string(contenido))
		if err != nil {
			return nil, fmt.Errorf("%s: %w", nombre, err)
		}
		sum := sha256.Sum256(contenido)
		out = append(out, migracion{
			version:    match[1],
			nombre:     nombre,
			checksum:   hex.EncodeToString(sum[:]),
			sentencias: sentencias,
		})
	}
	return out, nil
}

// dividir separa un script en sentencias respetando SET TERM.
func dividir(script string) ([]string, error) {
	var (
		out  []string
		buf  strings.Builder
		term = ";"
	)
	sc := bufio.NewScanner(strings.NewReader(strings.ReplaceAll(script, "\r\n", "\n")))
	for sc.Scan() {
		linea := strings.TrimRight(sc.Text(), " \t")
		recortada := strings.TrimSpace(linea)

		// Líneas vacías y comentarios entre sentencias se ignoran.
		if buf.Len() == 0 && (recortada == "" || strings.HasPrefix(recortada, "--")) {
			continue
		}
		if buf.Len() == 0 {
			// "SET TERM <nuevo> <actual>", p. ej. "SET TERM ^ ;" y luego "SET TERM ; ^".
			if m := setTerm.FindStringSubmatch(recortada); m != nil && m[2] == term {
				term = m[1]
				continue
			}
		}
		// Acumula líneas hasta encontrar el terminador actual.
		buf.WriteString(linea)
		buf.WriteByte('\n')
		if strings.HasSuffix(recortada, term) {
			s := strings.TrimSpace(buf.String())
			s = strings.TrimSpace(strings.TrimSuffix(s, term))
			if s != "" {
				out = append(out, s)
			}
			buf.Reset()
		}
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	if strings.TrimSpace(buf.String()) != "" {
		return nil, fmt.Errorf("sentencia sin terminar (falta %q al final)", term)
	}
	if term != ";" {
		return nil, fmt.Errorf("SET TERM sin cerrar: termina el script con SET TERM ; %s", term)
	}
	return out, nil
}
