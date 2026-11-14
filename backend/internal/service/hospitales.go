package service

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"salus/backend/internal/model"
	"salus/backend/internal/repository"
)

// ErrCampo es una regla de negocio que no se cumple en un campo concreto.
type ErrCampo struct{ Campo, Mensaje string }

func (e ErrCampo) Error() string { return e.Campo + ": " + e.Mensaje }

type Hospitales struct{ db *sql.DB }

func NewHospitales(db *sql.DB) *Hospitales { return &Hospitales{db: db} }

// DatosHospital es lo que se captura en alta y edición.
type DatosHospital struct {
	Nombre, Telefono, LigaMapa string
	Direccion                  model.Direccion
	Latitud, Longitud          *float64
	Contactos                  []model.Contacto
}

func (d DatosHospital) modelo(id model.ID) model.Hospital {
	cs := make([]model.Contacto, len(d.Contactos))
	for i, c := range d.Contactos {
		cs[i] = model.Contacto{Nombre: limpiar(c.Nombre), Puesto: limpiar(c.Puesto), Telefono: limpiar(c.Telefono)}
	}
	return model.Hospital{ID: id, Nombre: limpiar(d.Nombre),
		Telefono: limpiar(d.Telefono), LigaMapa: strings.TrimSpace(d.LigaMapa), Direccion: limpiarDireccion(d.Direccion),
		Latitud: d.Latitud, Longitud: d.Longitud, Contactos: cs}
}

func (s *Hospitales) Listar(ctx context.Context, q repository.Consulta) ([]model.Hospital, int, error) {
	return repository.NewHospitales(s.db).Listar(ctx, q)
}

func (s *Hospitales) Obtener(ctx context.Context, id model.ID) (*model.Hospital, error) {
	return repository.NewHospitales(s.db).PorID(ctx, id)
}

func (s *Hospitales) Crear(ctx context.Context, d DatosHospital) (*model.Hospital, error) {
	return s.guardar(ctx, d.modelo(model.NuevoID()), true)
}

func (s *Hospitales) Actualizar(ctx context.Context, id model.ID, d DatosHospital) (*model.Hospital, error) {
	return s.guardar(ctx, d.modelo(id), false)
}

// guardar escribe hospital y contactos en UNA transacción: todo o nada. En el
// alta genera la clave: abreviatura del estado de la colonia y su consecutivo
// (HGO-001). Después la clave no cambia.
func (s *Hospitales) guardar(ctx context.Context, h model.Hospital, nuevo bool) (*model.Hospital, error) {
	ahora := time.Now().UTC()
	h.ActualizadoEn = ahora
	if nuevo {
		h.CreadoEn = ahora
	}
	var out *model.Hospital
	err := enTx(ctx, s.db, func(tx *sql.Tx) error {
		repo := repository.NewHospitales(tx)
		var err error
		if nuevo {
			estado, err := repository.NewDirecciones(tx).EstadoDeAsentamiento(ctx, *h.Direccion.AsentamientoID)
			if err != nil {
				return err
			}
			n, err := repository.Siguiente(ctx, tx, "HOSPITAL:"+estado)
			if err != nil {
				return err
			}
			h.Clave = fmt.Sprintf("%s-%03d", estado, n)
			err = repo.Crear(ctx, h)
			if err != nil {
				return err
			}
		} else if err = repo.Actualizar(ctx, h); err != nil {
			return err
		}
		if err := repo.ReemplazarContactos(ctx, h.ID, h.Contactos); err != nil {
			return err
		}
		out, err = repo.PorID(ctx, h.ID)
		return err
	})
	return out, err
}

func (s *Hospitales) Baja(ctx context.Context, id model.ID) error {
	return repository.NewHospitales(s.db).Baja(ctx, id, time.Now().UTC())
}
