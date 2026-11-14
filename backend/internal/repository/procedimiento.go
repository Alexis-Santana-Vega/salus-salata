package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"salus/backend/internal/model"
)

type Procedimientos struct{ db DBTX }

func NewProcedimientos(db DBTX) *Procedimientos { return &Procedimientos{db: db} }

// Campos ordenables (API → columna) e índices únicos (índice → campo).
var (
	ordenProcedimientos   = map[string]string{"nombre": "NOMBRE", "creadoEn": "CREADO_EN"}
	indicesProcedimientos = map[string]string{"UQ_PROCEDIMIENTO_NOMBRE": "nombre"}
)

const selectProcedimiento = `SELECT ID, NOMBRE, CREADO_EN, ACTUALIZADO_EN FROM PROCEDIMIENTO`

func scanProcedimiento(s interface{ Scan(...any) error }) (*model.Procedimiento, error) {
	var p model.Procedimiento
	err := s.Scan(&p.ID, &p.Nombre, &p.CreadoEn, &p.ActualizadoEn)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNoEncontrado
	}
	return &p, err
}

// Listar devuelve una página de procedimientos activos y el total.
// La búsqueda usa la colación del campo: sin acentos ni mayúsculas.
func (r *Procedimientos) Listar(ctx context.Context, q Consulta) ([]model.Procedimiento, int, error) {
	orden, err := orderBy(q.Orden, ordenProcedimientos, "nombre", "ID")
	if err != nil {
		return nil, 0, err
	}
	where, args := " WHERE ACTIVO = TRUE", []any{}
	if q.Q != "" {
		where += " AND NOMBRE CONTAINING ?"
		args = append(args, q.Q)
	}

	var total int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM PROCEDIMIENTO`+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.db.QueryContext(ctx,
		selectProcedimiento+where+orden+` OFFSET ? ROWS FETCH NEXT ? ROWS ONLY`,
		append(args, q.Offset, q.Limite)...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	out := make([]model.Procedimiento, 0, q.Limite)
	for rows.Next() {
		p, err := scanProcedimiento(rows)
		if err != nil {
			return nil, 0, err
		}
		out = append(out, *p)
	}
	return out, total, rows.Err()
}

// PorID busca un procedimiento activo.
func (r *Procedimientos) PorID(ctx context.Context, id model.ID) (*model.Procedimiento, error) {
	return scanProcedimiento(r.db.QueryRowContext(ctx, selectProcedimiento+` WHERE ID = ? AND ACTIVO = TRUE`, id))
}

func (r *Procedimientos) Crear(ctx context.Context, p model.Procedimiento) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO PROCEDIMIENTO (ID, NOMBRE, CREADO_EN, ACTUALIZADO_EN) VALUES (?, ?, ?, ?)`,
		p.ID, p.Nombre, p.CreadoEn, p.ActualizadoEn)
	return duplicado(err, indicesProcedimientos)
}

// Actualizar cambia un procedimiento activo; ErrNoEncontrado si no existe.
func (r *Procedimientos) Actualizar(ctx context.Context, p model.Procedimiento) error {
	res, err := r.db.ExecContext(ctx, `
		UPDATE PROCEDIMIENTO SET NOMBRE = ?, ACTUALIZADO_EN = ? WHERE ID = ? AND ACTIVO = TRUE`,
		p.Nombre, p.ActualizadoEn, p.ID)
	if err != nil {
		return duplicado(err, indicesProcedimientos)
	}
	return unaFila(res)
}

// Baja hace la baja lógica: el registro se conserva para el historial.
func (r *Procedimientos) Baja(ctx context.Context, id model.ID, ahora time.Time) error {
	res, err := r.db.ExecContext(ctx, `
		UPDATE PROCEDIMIENTO SET ACTIVO = FALSE, ACTUALIZADO_EN = ? WHERE ID = ? AND ACTIVO = TRUE`, ahora, id)
	if err != nil {
		return err
	}
	return unaFila(res)
}
