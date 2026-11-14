package repository

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"time"

	"salus/backend/internal/model"
)

type Equipos struct{ db DBTX }

func NewEquipos(db DBTX) *Equipos { return &Equipos{db: db} }

var (
	ordenEquipos   = map[string]string{"clave": "E.CLAVE", "nombre": "E.NOMBRE", "marcaModelo": "E.MARCA_MODELO", "cicloMax": "E.CICLO_MAX", "creadoEn": "E.CREADO_EN"}
	indicesEquipos = map[string]string{"PK_CLAVE_CATALOGO": "clave", "UQ_EQUIPO_CLAVE": "clave"}
)

// La tarifa vigente sale con una subconsulta en la misma consulta (sin N+1).
const selectEquipo = `
SELECT E.ID, E.CLAVE, E.NOMBRE, E.MARCA_MODELO, E.CICLO_MAX, E.REQUIERE_PURGA, E.PURGA_CADA,
       (SELECT FIRST 1 T.COSTO_HORA FROM TARIFA_EQUIPO T
         WHERE T.EQUIPO_ID = E.ID AND T.VIGENTE_DESDE <= ? ORDER BY T.VIGENTE_DESDE DESC),
       E.CREADO_EN, E.ACTUALIZADO_EN
FROM EQUIPO E`

func scanEquipo(s interface{ Scan(...any) error }) (*model.Equipo, error) {
	var e model.Equipo
	var purga sql.NullInt64
	var costo sql.NullString
	err := s.Scan(&e.ID, &e.Clave, &e.Nombre, &e.MarcaModelo, &e.CicloMax, &e.RequierePurga, &purga, &costo, &e.CreadoEn, &e.ActualizadoEn)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNoEncontrado
	}
	if err != nil {
		return nil, err
	}
	if purga.Valid {
		p := int(purga.Int64)
		e.PurgaCada = &p
	}
	e.CostoHora, _ = strconv.ParseFloat(costo.String, 64)
	e.Tarifas = []model.Tarifa{}
	return &e, nil
}

// Listar devuelve una página de equipos activos con su costo vigente hoy.
func (r *Equipos) Listar(ctx context.Context, q Consulta, hoy time.Time) ([]model.Equipo, int, error) {
	orden, err := orderBy(q.Orden, ordenEquipos, "clave", "E.ID")
	if err != nil {
		return nil, 0, err
	}
	where, args := " WHERE E.ACTIVO = TRUE", []any{}
	if q.Q != "" {
		where += " AND (E.CLAVE CONTAINING ? OR E.NOMBRE CONTAINING ? OR E.MARCA_MODELO CONTAINING ?)"
		args = append(args, q.Q, q.Q, q.Q)
	}
	var total int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM EQUIPO E`+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	params := append(append([]any{hoy}, args...), q.Offset, q.Limite)
	rows, err := r.db.QueryContext(ctx, selectEquipo+where+orden+` OFFSET ? ROWS FETCH NEXT ? ROWS ONLY`, params...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	out := make([]model.Equipo, 0, q.Limite)
	for rows.Next() {
		e, err := scanEquipo(rows)
		if err != nil {
			return nil, 0, err
		}
		out = append(out, *e)
	}
	return out, total, rows.Err()
}

// PorID devuelve el equipo con su historial de tarifas.
func (r *Equipos) PorID(ctx context.Context, id model.ID, hoy time.Time) (*model.Equipo, error) {
	e, err := scanEquipo(r.db.QueryRowContext(ctx, selectEquipo+` WHERE E.ID = ? AND E.ACTIVO = TRUE`, hoy, id))
	if err != nil {
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx, `
		SELECT COSTO_HORA, VIGENTE_DESDE FROM TARIFA_EQUIPO WHERE EQUIPO_ID = ? ORDER BY VIGENTE_DESDE DESC`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var costo string
		var desde time.Time
		if err := rows.Scan(&costo, &desde); err != nil {
			return nil, err
		}
		c, _ := strconv.ParseFloat(costo, 64)
		e.Tarifas = append(e.Tarifas, model.Tarifa{CostoHora: c, VigenteDesde: desde.Format("2006-01-02")})
	}
	return e, rows.Err()
}

// ReservarClave registra la clave del equipo en CLAVE_CATALOGO.
func (r *Equipos) ReservarClave(ctx context.Context, clave string) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO CLAVE_CATALOGO (CLAVE, TIPO) VALUES (?, 'EQUIPO')`, clave)
	return integridad(err, indicesEquipos, nil)
}

func purgaSQL(e model.Equipo) any {
	if !e.RequierePurga || e.PurgaCada == nil {
		return nil
	}
	return *e.PurgaCada
}

func (r *Equipos) Crear(ctx context.Context, e model.Equipo) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO EQUIPO (ID, CLAVE, NOMBRE, MARCA_MODELO, CICLO_MAX, REQUIERE_PURGA, PURGA_CADA, CREADO_EN, ACTUALIZADO_EN)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		e.ID, e.Clave, e.Nombre, e.MarcaModelo, e.CicloMax, e.RequierePurga, purgaSQL(e), e.CreadoEn, e.ActualizadoEn)
	return integridad(err, indicesEquipos, nil)
}

// Actualizar cambia los datos del modelo (la clave no cambia).
func (r *Equipos) Actualizar(ctx context.Context, e model.Equipo) error {
	res, err := r.db.ExecContext(ctx, `
		UPDATE EQUIPO SET NOMBRE = ?, MARCA_MODELO = ?, CICLO_MAX = ?, REQUIERE_PURGA = ?, PURGA_CADA = ?, ACTUALIZADO_EN = ?
		WHERE ID = ? AND ACTIVO = TRUE`,
		e.Nombre, e.MarcaModelo, e.CicloMax, e.RequierePurga, purgaSQL(e), e.ActualizadoEn, e.ID)
	if err != nil {
		return err
	}
	return unaFila(res)
}

// FijarTarifa deja el costo vigente desde el día dado: si ya hay tarifa de
// ese día la corrige; si no, agrega una nueva (las anteriores no se tocan).
func (r *Equipos) FijarTarifa(ctx context.Context, equipo model.ID, costo float64, desde time.Time, ahora time.Time) error {
	res, err := r.db.ExecContext(ctx,
		`UPDATE TARIFA_EQUIPO SET COSTO_HORA = ? WHERE EQUIPO_ID = ? AND VIGENTE_DESDE = ?`, precioSQL(costo), equipo, desde)
	if err != nil {
		return err
	}
	if n, err := res.RowsAffected(); err != nil || n > 0 {
		return err
	}
	_, err = r.db.ExecContext(ctx, `
		INSERT INTO TARIFA_EQUIPO (ID, EQUIPO_ID, COSTO_HORA, VIGENTE_DESDE, CREADO_EN) VALUES (?, ?, ?, ?, ?)`,
		model.NuevoID(), equipo, precioSQL(costo), desde, ahora)
	return err
}

func (r *Equipos) Baja(ctx context.Context, id model.ID, ahora time.Time) error {
	res, err := r.db.ExecContext(ctx, `
		UPDATE EQUIPO SET ACTIVO = FALSE, ACTUALIZADO_EN = ? WHERE ID = ? AND ACTIVO = TRUE`, ahora, id)
	if err != nil {
		return err
	}
	return unaFila(res)
}
