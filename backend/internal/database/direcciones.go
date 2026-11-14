package database

import (
	"bytes"
	"compress/gzip"
	"context"
	"database/sql"
	_ "embed"
	"encoding/csv"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

// Catálogo de municipios (INEGI) y asentamientos con código postal (SEPOMEX),
// generado con cmd/catalogo-direcciones. Solo trae los estados elegidos.
//
//go:embed catalogos/direcciones.csv.gz
var direccionesGz []byte

type municipioCat struct{ ent, mun, nombre string }
type asentamientoCat struct{ ent, mun, clave, cp, nombre, tipo, zona, ciudad string }

// CargarDirecciones carga los estados del archivo que aún no estén en la base.
// Cada estado va en su propia transacción (municipios y asentamientos juntos o
// nada). Para sumar un estado basta regenerar el archivo: los ya cargados no se
// tocan. Corre después de Migrate.
func CargarDirecciones(ctx context.Context, db *sql.DB, log *slog.Logger) error {
	municipios, asentamientos, err := leerDirecciones()
	if err != nil {
		return err
	}
	porEstado := map[string][]municipioCat{}
	var orden []string
	for _, m := range municipios {
		if _, ok := porEstado[m.ent]; !ok {
			orden = append(orden, m.ent)
		}
		porEstado[m.ent] = append(porEstado[m.ent], m)
	}
	asenPorEstado := map[string][]asentamientoCat{}
	for _, a := range asentamientos {
		asenPorEstado[a.ent] = append(asenPorEstado[a.ent], a)
	}

	for _, ent := range orden {
		inicio := time.Now()
		cargado, err := cargarEstado(ctx, db, ent, porEstado[ent], asenPorEstado[ent])
		if err != nil {
			return fmt.Errorf("catálogo de direcciones, estado %s: %w", ent, err)
		}
		if cargado {
			log.Info("catálogo de direcciones cargado", "estado", ent,
				"municipios", len(porEstado[ent]), "asentamientos", len(asenPorEstado[ent]), "ms", time.Since(inicio).Milliseconds())
		}
	}
	return nil
}

func cargarEstado(ctx context.Context, db *sql.DB, ent string, ms []municipioCat, as []asentamientoCat) (bool, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return false, err
	}
	defer func() { _ = tx.Rollback() }()

	// ¿Ya está? (se revisa dentro de la transacción).
	var estadoID []byte
	var hay int
	if err := tx.QueryRowContext(ctx, `
		SELECT E.ID, (SELECT COUNT(*) FROM MUNICIPIO M WHERE M.ESTADO_ID = E.ID)
		FROM ESTADO E WHERE E.CLAVE = ?`, ent).Scan(&estadoID, &hay); err != nil {
		return false, fmt.Errorf("estado %s no existe en ESTADO: %w", ent, err)
	}
	if hay > 0 {
		return false, nil
	}

	insMun, err := tx.PrepareContext(ctx, `INSERT INTO MUNICIPIO (ID, ESTADO_ID, CLAVE, NOMBRE) VALUES (?, ?, ?, ?)`)
	if err != nil {
		return false, err
	}
	defer insMun.Close()
	ids := map[string][]byte{}
	for _, m := range ms {
		id := nuevoID()
		if _, err := insMun.ExecContext(ctx, id, estadoID, m.mun, m.nombre); err != nil {
			return false, fmt.Errorf("municipio %s: %w", m.mun, err)
		}
		ids[m.mun] = id
	}

	insAsen, err := tx.PrepareContext(ctx, `
		INSERT INTO ASENTAMIENTO (ID, MUNICIPIO_ID, CLAVE, CP, NOMBRE, TIPO, ZONA, CIUDAD)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return false, err
	}
	defer insAsen.Close()
	for _, a := range as {
		if _, err := insAsen.ExecContext(ctx, nuevoID(), ids[a.mun], a.clave, a.cp, a.nombre, a.tipo, nulo(a.zona), nulo(a.ciudad)); err != nil {
			return false, fmt.Errorf("asentamiento %s-%s: %w", a.mun, a.clave, err)
		}
	}
	return true, tx.Commit()
}

func leerDirecciones() ([]municipioCat, []asentamientoCat, error) {
	gz, err := gzip.NewReader(bytes.NewReader(direccionesGz))
	if err != nil {
		return nil, nil, err
	}
	r := csv.NewReader(gz)
	r.FieldsPerRecord = -1
	var ms []municipioCat
	var as []asentamientoCat
	for {
		f, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, err
		}
		switch {
		case f[0] == "M" && len(f) == 4:
			ms = append(ms, municipioCat{f[1], f[2], f[3]})
		case f[0] == "A" && len(f) == 9:
			as = append(as, asentamientoCat{f[1], f[2], f[3], f[4], f[5], f[6], f[7], f[8]})
		default:
			return nil, nil, fmt.Errorf("línea inválida en direcciones.csv.gz: %v", f)
		}
	}
	return ms, as, nil
}

// nuevoID: UUID v7 en 16 bytes (database no depende de model).
func nuevoID() []byte {
	u := uuid.Must(uuid.NewV7())
	return u[:]
}

func nulo(s string) any {
	if s == "" {
		return nil
	}
	return s
}
