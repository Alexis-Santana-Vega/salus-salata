package repository

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"

	"salus/backend/internal/model"
)

type Hospitales struct{ db DBTX }

func NewHospitales(db DBTX) *Hospitales { return &Hospitales{db: db} }

var (
	ordenHospitales    = map[string]string{"clave": "H.CLAVE", "nombre": "H.NOMBRE", "municipio": "M.NOMBRE", "creadoEn": "H.CREADO_EN"}
	indicesHospitales  = map[string]string{"UQ_HOSPITAL_CLAVE": "clave"}
	foraneasHospitales = map[string]string{"FK_HOSPITAL_ASENTAMIENTO": "asentamientoId"}
)

var selectHospital = `
SELECT H.ID, H.CLAVE, H.NOMBRE, H.TELEFONO, ` + columnasDireccion("H") + `,
       H.LATITUD, H.LONGITUD, H.LIGA_MAPA, H.CREADO_EN, H.ACTUALIZADO_EN
FROM HOSPITAL H` + joinDireccion("H")

func scanHospital(s interface{ Scan(...any) error }) (*model.Hospital, error) {
	var h model.Hospital
	var tel, lat, lon, liga sql.NullString
	var dir direccionScan
	dest := append([]any{&h.ID, &h.Clave, &h.Nombre, &tel}, dir.destinos()...)
	dest = append(dest, &lat, &lon, &liga, &h.CreadoEn, &h.ActualizadoEn)
	err := s.Scan(dest...)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNoEncontrado
	}
	if err != nil {
		return nil, err
	}
	h.Telefono, h.LigaMapa, h.Direccion = tel.String, liga.String, dir.modelo()
	h.Latitud, h.Longitud = decimal(lat), decimal(lon)
	h.Contactos = []model.Contacto{}
	return &h, nil
}

// decimal: el driver entrega NUMERIC como texto.
func decimal(s sql.NullString) *float64 {
	if !s.Valid {
		return nil
	}
	f, err := strconv.ParseFloat(s.String, 64)
	if err != nil {
		return nil
	}
	return &f
}

// Listar devuelve una página de hospitales activos con sus contactos (los
// contactos de toda la página en una sola consulta con IN: sin N+1).
func (r *Hospitales) Listar(ctx context.Context, q Consulta) ([]model.Hospital, int, error) {
	orden, err := orderBy(q.Orden, ordenHospitales, "nombre", "H.ID")
	if err != nil {
		return nil, 0, err
	}
	where, args := " WHERE H.ACTIVO = TRUE", []any{}
	if q.Q != "" {
		where += " AND (H.NOMBRE CONTAINING ? OR H.CLAVE CONTAINING ? OR M.NOMBRE CONTAINING ?)"
		args = append(args, q.Q, q.Q, q.Q)
	}

	var total int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM HOSPITAL H`+joinDireccion("H")+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.db.QueryContext(ctx,
		selectHospital+where+orden+` OFFSET ? ROWS FETCH NEXT ? ROWS ONLY`,
		append(args, q.Offset, q.Limite)...)
	if err != nil {
		return nil, 0, err
	}
	out := make([]model.Hospital, 0, q.Limite)
	for rows.Next() {
		h, err := scanHospital(rows)
		if err != nil {
			rows.Close()
			return nil, 0, err
		}
		out = append(out, *h)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return out, total, r.cargarContactos(ctx, out)
}

// cargarContactos llena los contactos de varios hospitales con un solo IN.
func (r *Hospitales) cargarContactos(ctx context.Context, hs []model.Hospital) error {
	if len(hs) == 0 {
		return nil
	}
	pos := make(map[model.ID]int, len(hs))
	args := make([]any, len(hs))
	for i, h := range hs {
		pos[h.ID] = i
		args[i] = h.ID
	}
	rows, err := r.db.QueryContext(ctx, `
		SELECT HOSPITAL_ID, NOMBRE, PUESTO, TELEFONO FROM CONTACTO_HOSPITAL
		WHERE HOSPITAL_ID IN (?`+strings.Repeat(", ?", len(hs)-1)+`)
		ORDER BY HOSPITAL_ID, ORDEN`, args...)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var id model.ID
		var c model.Contacto
		var puesto, tel sql.NullString
		if err := rows.Scan(&id, &c.Nombre, &puesto, &tel); err != nil {
			return err
		}
		c.Puesto, c.Telefono = puesto.String, tel.String
		i := pos[id]
		hs[i].Contactos = append(hs[i].Contactos, c)
	}
	return rows.Err()
}

func (r *Hospitales) PorID(ctx context.Context, id model.ID) (*model.Hospital, error) {
	h, err := scanHospital(r.db.QueryRowContext(ctx, selectHospital+` WHERE H.ID = ? AND H.ACTIVO = TRUE`, id))
	if err != nil {
		return nil, err
	}
	lista := []model.Hospital{*h}
	if err := r.cargarContactos(ctx, lista); err != nil {
		return nil, err
	}
	return &lista[0], nil
}

func camposHospital(h model.Hospital) []any {
	args := []any{h.Clave, h.Nombre, nulo(h.Telefono)}
	args = append(args, valoresDireccion(h.Direccion)...)
	return append(args, numero(h.Latitud), numero(h.Longitud), nulo(h.LigaMapa))
}

func numero(f *float64) any {
	if f == nil {
		return nil
	}
	return *f
}

func (r *Hospitales) Crear(ctx context.Context, h model.Hospital) error {
	args := append([]any{h.ID}, camposHospital(h)...)
	args = append(args, h.CreadoEn, h.ActualizadoEn)
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO HOSPITAL (ID, CLAVE, NOMBRE, TELEFONO, CALLE, NUM_EXTERIOR, NUM_INTERIOR, ASENTAMIENTO_ID,
		                      LATITUD, LONGITUD, LIGA_MAPA, CREADO_EN, ACTUALIZADO_EN)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, args...)
	return integridad(err, indicesHospitales, foraneasHospitales)
}

func (r *Hospitales) Actualizar(ctx context.Context, h model.Hospital) error {
	args := append(camposHospital(h)[1:], h.ActualizadoEn, h.ID) // la clave no cambia
	res, err := r.db.ExecContext(ctx, `
		UPDATE HOSPITAL SET NOMBRE = ?, TELEFONO = ?, CALLE = ?, NUM_EXTERIOR = ?, NUM_INTERIOR = ?,
		       ASENTAMIENTO_ID = ?, LATITUD = ?, LONGITUD = ?, LIGA_MAPA = ?, ACTUALIZADO_EN = ?
		WHERE ID = ? AND ACTIVO = TRUE`, args...)
	if err != nil {
		return integridad(err, indicesHospitales, foraneasHospitales)
	}
	return unaFila(res)
}

// ReemplazarContactos deja exactamente los contactos dados (va dentro de la
// transacción del alta o edición del hospital).
func (r *Hospitales) ReemplazarContactos(ctx context.Context, hospital model.ID, cs []model.Contacto) error {
	if _, err := r.db.ExecContext(ctx, `DELETE FROM CONTACTO_HOSPITAL WHERE HOSPITAL_ID = ?`, hospital); err != nil {
		return err
	}
	for i, c := range cs {
		if _, err := r.db.ExecContext(ctx, `
			INSERT INTO CONTACTO_HOSPITAL (ID, HOSPITAL_ID, ORDEN, NOMBRE, PUESTO, TELEFONO) VALUES (?, ?, ?, ?, ?, ?)`,
			model.NuevoID(), hospital, i+1, c.Nombre, nulo(c.Puesto), nulo(c.Telefono)); err != nil {
			return err
		}
	}
	return nil
}

func (r *Hospitales) Baja(ctx context.Context, id model.ID, ahora time.Time) error {
	res, err := r.db.ExecContext(ctx, `
		UPDATE HOSPITAL SET ACTIVO = FALSE, ACTUALIZADO_EN = ? WHERE ID = ? AND ACTIVO = TRUE`, ahora, id)
	if err != nil {
		return err
	}
	return unaFila(res)
}
