package repository

import (
	"context"
	"database/sql"
	"errors"

	"salus/backend/internal/model"
)

type Direcciones struct{ db DBTX }

func NewDirecciones(db DBTX) *Direcciones { return &Direcciones{db: db} }

// PorCP devuelve los asentamientos de un código postal con su municipio y
// estado (una sola consulta con JOIN).
func (r *Direcciones) PorCP(ctx context.Context, cp string) ([]model.Asentamiento, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT A.ID, A.CP, A.NOMBRE, A.TIPO, A.ZONA, A.CIUDAD, M.NOMBRE, E.ABREVIATURA, E.NOMBRE
		FROM ASENTAMIENTO A
		JOIN MUNICIPIO M ON M.ID = A.MUNICIPIO_ID
		JOIN ESTADO E ON E.ID = M.ESTADO_ID
		WHERE A.CP = ?
		ORDER BY A.NOMBRE`, cp)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []model.Asentamiento{}
	for rows.Next() {
		var a model.Asentamiento
		var zona, ciudad sql.NullString
		if err := rows.Scan(&a.ID, &a.CP, &a.Nombre, &a.Tipo, &zona, &ciudad, &a.Municipio, &a.Estado, &a.EstadoNombre); err != nil {
			return nil, err
		}
		a.Zona, a.Ciudad = zona.String, ciudad.String
		out = append(out, a)
	}
	return out, rows.Err()
}

// Columnas de dirección para un SELECT con LEFT JOIN al catálogo (alias
// A, M y E); tabla es el alias de la entidad dueña (D, H…).
func columnasDireccion(tabla string) string {
	return tabla + `.CALLE, ` + tabla + `.NUM_EXTERIOR, ` + tabla + `.NUM_INTERIOR, ` + tabla + `.ASENTAMIENTO_ID,
		A.CP, A.NOMBRE, M.NOMBRE, E.ABREVIATURA, E.NOMBRE`
}

func joinDireccion(tabla string) string {
	return `
		LEFT JOIN ASENTAMIENTO A ON A.ID = ` + tabla + `.ASENTAMIENTO_ID
		LEFT JOIN MUNICIPIO M ON M.ID = A.MUNICIPIO_ID
		LEFT JOIN ESTADO E ON E.ID = M.ESTADO_ID`
}

// direccionScan recibe las columnas de columnasDireccion.
type direccionScan struct {
	calle, numExt, numInt, cp, colonia, municipio, estado, estadoNombre sql.NullString
	asentamiento                                                        model.ID
}

func (d *direccionScan) destinos() []any {
	return []any{&d.calle, &d.numExt, &d.numInt, &d.asentamiento, &d.cp, &d.colonia, &d.municipio, &d.estado, &d.estadoNombre}
}

func (d *direccionScan) modelo() model.Direccion {
	out := model.Direccion{Calle: d.calle.String, NumExterior: d.numExt.String, NumInterior: d.numInt.String,
		CP: d.cp.String, Colonia: d.colonia.String, Municipio: d.municipio.String,
		Estado: d.estado.String, EstadoNombre: d.estadoNombre.String}
	if !d.asentamiento.IsZero() {
		id := d.asentamiento
		out.AsentamientoID = &id
	}
	return out
}

// valoresDireccion: calle, números y asentamiento en el orden de las columnas.
func valoresDireccion(a model.Direccion) []any {
	var asen any
	if a.AsentamientoID != nil {
		asen = *a.AsentamientoID
	}
	return []any{nulo(a.Calle), nulo(a.NumExterior), nulo(a.NumInterior), asen}
}

// EstadoDeAsentamiento devuelve la abreviatura del estado de una colonia.
func (r *Direcciones) EstadoDeAsentamiento(ctx context.Context, id model.ID) (string, error) {
	var abr string
	err := r.db.QueryRowContext(ctx, `
		SELECT E.ABREVIATURA FROM ASENTAMIENTO A
		JOIN MUNICIPIO M ON M.ID = A.MUNICIPIO_ID
		JOIN ESTADO E ON E.ID = M.ESTADO_ID
		WHERE A.ID = ?`, id).Scan(&abr)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrReferencia{Campo: "asentamientoId"}
	}
	return abr, err
}
