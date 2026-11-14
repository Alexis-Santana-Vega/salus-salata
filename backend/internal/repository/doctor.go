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
	ordenDoctores    = map[string]string{"nombre": "D.NOMBRE", "numero": "D.NUMERO", "correo": "D.CORREO", "municipio": "M.NOMBRE", "creadoEn": "D.CREADO_EN"}
	indicesDoctores  = map[string]string{"UQ_DOCTOR_NUMERO": "numero", "UQ_DOCTOR_CORREO": "correo"}
	foraneasDoctores = map[string]string{"FK_DOCTOR_ASENTAMIENTO": "asentamientoId"}
)

// Dirección con LEFT JOIN al catálogo: una sola consulta (sin N+1).
var selectDoctor = `
SELECT D.ID, D.NUMERO, D.NOMBRE, D.CORREO, D.TEL_PERSONAL, D.TEL_CONSULTORIO,
       ` + columnasDireccion("D") + `, D.CREADO_EN, D.ACTUALIZADO_EN
FROM DOCTOR D` + joinDireccion("D")

func scanDoctor(s interface{ Scan(...any) error }) (*model.Doctor, error) {
	var d model.Doctor
	var numero, telP, telC sql.NullString
	var dir direccionScan
	dest := append([]any{&d.ID, &numero, &d.Nombre, &d.Correo, &telP, &telC}, dir.destinos()...)
	dest = append(dest, &d.CreadoEn, &d.ActualizadoEn)
	err := s.Scan(dest...)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNoEncontrado
	}
	if err != nil {
		return nil, err
	}
	d.Numero, d.TelPersonal, d.TelConsultorio = numero.String, telP.String, telC.String
	d.Direccion = dir.modelo()
	return &d, nil
}

// Listar devuelve una página de doctores activos; busca por nombre, número o correo.
func (r *Doctores) Listar(ctx context.Context, q Consulta) ([]model.Doctor, int, error) {
	orden, err := orderBy(q.Orden, ordenDoctores, "nombre", "D.ID")
	if err != nil {
		return nil, 0, err
	}
	where, args := " WHERE D.ACTIVO = TRUE", []any{}
	if q.Q != "" {
		where += " AND (D.NOMBRE CONTAINING ? OR D.NUMERO CONTAINING ? OR D.CORREO CONTAINING ?)"
		args = append(args, q.Q, q.Q, q.Q)
	}

	var total int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM DOCTOR D`+where, args...).Scan(&total); err != nil {
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
	return scanDoctor(r.db.QueryRowContext(ctx, selectDoctor+` WHERE D.ID = ? AND D.ACTIVO = TRUE`, id))
}

// campos editables en el orden de las columnas.
func camposDoctor(d model.Doctor) []any {
	return append([]any{nulo(d.Numero), d.Nombre, d.Correo, nulo(d.TelPersonal), nulo(d.TelConsultorio)},
		valoresDireccion(d.Direccion)...)
}

func (r *Doctores) Crear(ctx context.Context, d model.Doctor) error {
	args := append([]any{d.ID}, camposDoctor(d)...)
	args = append(args, d.CreadoEn, d.ActualizadoEn)
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO DOCTOR (ID, NUMERO, NOMBRE, CORREO, TEL_PERSONAL, TEL_CONSULTORIO,
		                    CALLE, NUM_EXTERIOR, NUM_INTERIOR, ASENTAMIENTO_ID, CREADO_EN, ACTUALIZADO_EN)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, args...)
	return integridad(err, indicesDoctores, foraneasDoctores)
}

func (r *Doctores) Actualizar(ctx context.Context, d model.Doctor) error {
	args := append(camposDoctor(d), d.ActualizadoEn, d.ID)
	res, err := r.db.ExecContext(ctx, `
		UPDATE DOCTOR SET NUMERO = ?, NOMBRE = ?, CORREO = ?, TEL_PERSONAL = ?, TEL_CONSULTORIO = ?,
		       CALLE = ?, NUM_EXTERIOR = ?, NUM_INTERIOR = ?, ASENTAMIENTO_ID = ?, ACTUALIZADO_EN = ?
		WHERE ID = ? AND ACTIVO = TRUE`, args...)
	if err != nil {
		return integridad(err, indicesDoctores, foraneasDoctores)
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
