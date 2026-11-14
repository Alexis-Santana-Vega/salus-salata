package service

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"salus/backend/internal/model"
	"salus/backend/internal/repository"
)

type Doctores struct{ db *sql.DB }

func NewDoctores(db *sql.DB) *Doctores { return &Doctores{db: db} }

// DatosDoctor es lo que se captura en alta y edición.
type DatosDoctor struct {
	Numero, Nombre, Correo, TelPersonal, TelConsultorio string
	Direccion                                           model.Direccion
}

// normalizar limpia espacios; número en mayúsculas y correo en minúsculas
// para que "doc-1" y "DOC-1" (o dos correos con distinta capitalización) choquen.
func (d DatosDoctor) normalizar() DatosDoctor {
	d.Numero = strings.ToUpper(limpiar(d.Numero))
	d.Nombre = limpiar(d.Nombre)
	d.Correo = Normalizar(d.Correo)
	d.TelPersonal, d.TelConsultorio = limpiar(d.TelPersonal), limpiar(d.TelConsultorio)
	d.Direccion = limpiarDireccion(d.Direccion)
	return d
}

func limpiarDireccion(a model.Direccion) model.Direccion {
	return model.Direccion{Calle: limpiar(a.Calle), NumExterior: limpiar(a.NumExterior), NumInterior: limpiar(a.NumInterior),
		Colonia: limpiar(a.Colonia), CP: limpiar(a.CP), Municipio: limpiar(a.Municipio), Estado: a.Estado}
}

func (d DatosDoctor) modelo(id model.ID) model.Doctor {
	return model.Doctor{ID: id, Numero: d.Numero, Nombre: d.Nombre, Correo: d.Correo,
		TelPersonal: d.TelPersonal, TelConsultorio: d.TelConsultorio, Direccion: d.Direccion}
}

func (s *Doctores) Listar(ctx context.Context, q repository.Consulta) ([]model.Doctor, int, error) {
	return repository.NewDoctores(s.db).Listar(ctx, q)
}

func (s *Doctores) Obtener(ctx context.Context, id model.ID) (*model.Doctor, error) {
	return repository.NewDoctores(s.db).PorID(ctx, id)
}

func (s *Doctores) Crear(ctx context.Context, d DatosDoctor) (*model.Doctor, error) {
	doc := d.normalizar().modelo(model.NuevoID())
	doc.CreadoEn = time.Now().UTC()
	doc.ActualizadoEn = doc.CreadoEn
	if err := repository.NewDoctores(s.db).Crear(ctx, doc); err != nil {
		return nil, err
	}
	return &doc, nil
}

// Actualizar y releer en la misma transacción: se devuelve justo lo guardado.
func (s *Doctores) Actualizar(ctx context.Context, id model.ID, d DatosDoctor) (*model.Doctor, error) {
	doc := d.normalizar().modelo(id)
	doc.ActualizadoEn = time.Now().UTC()
	var out *model.Doctor
	err := enTx(ctx, s.db, func(tx *sql.Tx) error {
		repo := repository.NewDoctores(tx)
		if err := repo.Actualizar(ctx, doc); err != nil {
			return err
		}
		var err error
		out, err = repo.PorID(ctx, id)
		return err
	})
	return out, err
}

func (s *Doctores) Baja(ctx context.Context, id model.ID) error {
	return repository.NewDoctores(s.db).Baja(ctx, id, time.Now().UTC())
}
