package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"salus/backend/internal/model"
)

type Doctores struct{ db DBTX }

func NewDoctores(db DBTX) *Doctores { return &Doctores{db: db} }

var (
	ordenDoctores   = map[string]string{"nombre": "NOMBRE", "numero": "NUMERO", "correo": "CORREO", "municipio": "MUNICIPIO", "creadoEn": "CREADO_EN"}
	indicesDoctores = map[string]string{"UQ_DOCTOR_NUMERO": "numero", "UQ_DOCTOR_CORREO": "correo"}
)

const selectDoctor = `
SELECT ID, NUMERO, NOMBRE, CORREO, TEL_PERSONAL, TEL_CONSULTORIO,
       CALLE, NUM_EXTERIOR, NUM_INTERIOR, COLONIA, CP, MUNICIPIO, ESTADO,
       CREADO_EN, ACTUALIZADO_EN
FROM DOCTOR`

func scanDoctor(s interface{ Scan(...any) error }) (*model.Doctor, error) {
	var d model.Doctor
	var t [11]sql.NullString // opcionales: NULL = ""
	err := s.Scan(&d.ID, &t[0], &d.Nombre, &d.Correo, &t[1], &t[2],
		&t[3], &t[4], &t[5], &t[6], &t[7], &t[8], &t[9], &d.CreadoEn, &d.ActualizadoEn)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNoEncontrado
	}
	if err != nil {
		return nil, err
	}
	d.Numero, d.TelPersonal, d.TelConsultorio = t[0].String, t[1].String, t[2].String
	d.Direccion = model.Direccion{Calle: t[3].String, NumExterior: t[4].String, NumInterior: t[5].String,
		Colonia: t[6].String, CP: t[7].String, Municipio: t[8].String, Estado: t[9].String}
	return &d, nil
}

// Listar devuelve una página de doctores activos; busca por nombre, número o correo.
func (r *Doctores) Listar(ctx context.Context, q Consulta) ([]model.Doctor, int, error) {
	orden, err := orderBy(q.Orden, ordenDoctores, "nombre")
	if err != nil {
		return nil, 0, err
	}
	where, args := " WHERE ACTIVO = TRUE", []any{}
	if q.Q != "" {
		where += " AND (NOMBRE CONTAINING ? OR NUMERO CONTAINING ? OR CORREO CONTAINING ?)"
		args = append(args, q.Q, q.Q, q.Q)
	}

	var total int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM DOCTOR`+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.db.QueryContext(ctx,
		selectDoctor+where+orden+` OFFSET ? ROWS FETCH NEXT ? ROWS ONLY`,
		append(args, q.Offset, q.Limite)...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	out := make([]model.Doctor, 0, q.Limite)
	for rows.Next() {
		d, err := scanDoctor(rows)
		if err != nil {
			return nil, 0, err
		}
		out = append(out, *d)
	}
	return out, total, rows.Err()
}

func (r *Doctores) PorID(ctx context.Context, id model.ID) (*model.Doctor, error) {
	return scanDoctor(r.db.QueryRowContext(ctx, selectDoctor+` WHERE ID = ? AND ACTIVO = TRUE`, id))
}

// campos en el orden de las columnas editables.
func camposDoctor(d model.Doctor) []any {
	a := d.Direccion
	return []any{nulo(d.Numero), d.Nombre, d.Correo, nulo(d.TelPersonal), nulo(d.TelConsultorio),
		nulo(a.Calle), nulo(a.NumExterior), nulo(a.NumInterior), nulo(a.Colonia), nulo(a.CP), nulo(a.Municipio), nulo(a.Estado)}
}

func (r *Doctores) Crear(ctx context.Context, d model.Doctor) error {
	args := append([]any{d.ID}, camposDoctor(d)...)
	args = append(args, d.CreadoEn, d.ActualizadoEn)
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO DOCTOR (ID, NUMERO, NOMBRE, CORREO, TEL_PERSONAL, TEL_CONSULTORIO,
		                    CALLE, NUM_EXTERIOR, NUM_INTERIOR, COLONIA, CP, MUNICIPIO, ESTADO,
		                    CREADO_EN, ACTUALIZADO_EN)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, args...)
	return duplicado(err, indicesDoctores)
}

func (r *Doctores) Actualizar(ctx context.Context, d model.Doctor) error {
	args := append(camposDoctor(d), d.ActualizadoEn, d.ID)
	res, err := r.db.ExecContext(ctx, `
		UPDATE DOCTOR SET NUMERO = ?, NOMBRE = ?, CORREO = ?, TEL_PERSONAL = ?, TEL_CONSULTORIO = ?,
		       CALLE = ?, NUM_EXTERIOR = ?, NUM_INTERIOR = ?, COLONIA = ?, CP = ?, MUNICIPIO = ?, ESTADO = ?,
		       ACTUALIZADO_EN = ?
		WHERE ID = ? AND ACTIVO = TRUE`, args...)
	if err != nil {
		return duplicado(err, indicesDoctores)
	}
	return unaFila(res)
}

func (r *Doctores) Baja(ctx context.Context, id model.ID, ahora time.Time) error {
	res, err := r.db.ExecContext(ctx, `
		UPDATE DOCTOR SET ACTIVO = FALSE, ACTUALIZADO_EN = ? WHERE ID = ? AND ACTIVO = TRUE`, ahora, id)
	if err != nil {
		return err
	}
	return unaFila(res)
}
