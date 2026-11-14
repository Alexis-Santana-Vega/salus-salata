package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"salus/backend/internal/model"
)

type Tecnicos struct{ db DBTX }

func NewTecnicos(db DBTX) *Tecnicos { return &Tecnicos{db: db} }

var ordenTecnicos = map[string]string{"nombre": "NOMBRE", "puesto": "PUESTO", "creadoEn": "CREADO_EN"}

const selectTecnico = `SELECT ID, NOMBRE, PUESTO, TELEFONO, CREADO_EN, ACTUALIZADO_EN FROM TECNICO`

func scanTecnico(s interface{ Scan(...any) error }) (*model.Tecnico, error) {
	var t model.Tecnico
	var tel sql.NullString
	err := s.Scan(&t.ID, &t.Nombre, &t.Puesto, &tel, &t.CreadoEn, &t.ActualizadoEn)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNoEncontrado
	}
	t.Telefono = tel.String
	return &t, err
}

// Listar devuelve una página de técnicos activos; busca por nombre o teléfono.
func (r *Tecnicos) Listar(ctx context.Context, q Consulta) ([]model.Tecnico, int, error) {
	orden, err := orderBy(q.Orden, ordenTecnicos, "nombre")
	if err != nil {
		return nil, 0, err
	}
	where, args := " WHERE ACTIVO = TRUE", []any{}
	if q.Q != "" {
		where += " AND (NOMBRE CONTAINING ? OR TELEFONO CONTAINING ?)"
		args = append(args, q.Q, q.Q)
	}

	var total int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM TECNICO`+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.db.QueryContext(ctx,
		selectTecnico+where+orden+` OFFSET ? ROWS FETCH NEXT ? ROWS ONLY`,
		append(args, q.Offset, q.Limite)...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	out := make([]model.Tecnico, 0, q.Limite)
	for rows.Next() {
		t, err := scanTecnico(rows)
		if err != nil {
			return nil, 0, err
		}
		out = append(out, *t)
	}
	return out, total, rows.Err()
}

func (r *Tecnicos) PorID(ctx context.Context, id model.ID) (*model.Tecnico, error) {
	return scanTecnico(r.db.QueryRowContext(ctx, selectTecnico+` WHERE ID = ? AND ACTIVO = TRUE`, id))
}

func (r *Tecnicos) Crear(ctx context.Context, t model.Tecnico) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO TECNICO (ID, NOMBRE, PUESTO, TELEFONO, CREADO_EN, ACTUALIZADO_EN) VALUES (?, ?, ?, ?, ?, ?)`,
		t.ID, t.Nombre, t.Puesto, nulo(t.Telefono), t.CreadoEn, t.ActualizadoEn)
	return err
}

func (r *Tecnicos) Actualizar(ctx context.Context, t model.Tecnico) error {
	res, err := r.db.ExecContext(ctx, `
		UPDATE TECNICO SET NOMBRE = ?, PUESTO = ?, TELEFONO = ?, ACTUALIZADO_EN = ?
		WHERE ID = ? AND ACTIVO = TRUE`,
		t.Nombre, t.Puesto, nulo(t.Telefono), t.ActualizadoEn, t.ID)
	if err != nil {
		return err
	}
	return unaFila(res)
}

func (r *Tecnicos) Baja(ctx context.Context, id model.ID, ahora time.Time) error {
	res, err := r.db.ExecContext(ctx, `
		UPDATE TECNICO SET ACTIVO = FALSE, ACTUALIZADO_EN = ? WHERE ID = ? AND ACTIVO = TRUE`, ahora, id)
	if err != nil {
		return err
	}
	return unaFila(res)
}
