package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"salus/backend/internal/model"
	"salus/backend/internal/repository"
)

type Equipos struct{ db *sql.DB }

func NewEquipos(db *sql.DB) *Equipos { return &Equipos{db: db} }

// DatosEquipo es lo que se captura en alta y edición.
type DatosEquipo struct {
	Nombre, MarcaModelo string
	CicloMax            int
	RequierePurga       bool
	PurgaCada           *int
	CostoHora           float64
	Unidades            []DatosUnidad // estado completo: las guardadas que falten se dan de baja
}

// DatosUnidad: sin ID es una unidad nueva; con ID, una ya guardada.
type DatosUnidad struct {
	ID          *model.ID
	NumSerie    string
	FechaAlta   time.Time
	UsosPrevios int
}

func (d DatosEquipo) modelo(id model.ID) model.Equipo {
	e := model.Equipo{ID: id, Nombre: limpiar(d.Nombre), MarcaModelo: limpiar(d.MarcaModelo),
		CicloMax: d.CicloMax, RequierePurga: d.RequierePurga}
	if d.RequierePurga {
		e.PurgaCada = d.PurgaCada // sin purga no se guarda "cada cuántos"
	}
	return e
}

func (s *Equipos) Listar(ctx context.Context, q repository.Consulta) ([]model.Equipo, int, error) {
	return repository.NewEquipos(s.db).Listar(ctx, q, model.Hoy())
}

func (s *Equipos) Obtener(ctx context.Context, id model.ID) (*model.Equipo, error) {
	return repository.NewEquipos(s.db).PorID(ctx, id, model.Hoy())
}

// Crear genera la clave (EQ-001), la reserva, da de alta el modelo, su primera
// tarifa (vigente desde hoy) y sus unidades en UNA transacción.
func (s *Equipos) Crear(ctx context.Context, d DatosEquipo) (*model.Equipo, error) {
	e := d.modelo(model.NuevoID())
	ahora, hoy := time.Now().UTC(), model.Hoy()
	e.CreadoEn, e.ActualizadoEn = ahora, ahora
	if err := validarUnidades(d.Unidades, hoy); err != nil {
		return nil, err
	}
	var out *model.Equipo
	err := enTx(ctx, s.db, func(tx *sql.Tx) error {
		repo := repository.NewEquipos(tx)
		n, err := repository.Siguiente(ctx, tx, "EQUIPO")
		if err != nil {
			return err
		}
		e.Clave = fmt.Sprintf("EQ-%03d", n)
		if err := repo.ReservarClave(ctx, e.Clave); err != nil {
			return err
		}
		if err := repo.Crear(ctx, e); err != nil {
			return err
		}
		if err := repo.FijarTarifa(ctx, e.ID, d.CostoHora, hoy, ahora); err != nil {
			return err
		}
		if err := guardarUnidades(ctx, tx, e, d.Unidades, nil, ahora); err != nil {
			return err
		}
		out, err = repo.PorID(ctx, e.ID, hoy)
		return err
	})
	return out, err
}

// Actualizar cambia el modelo, agrega la tarifa de hoy solo si el costo cambió
// (las anteriores quedan como historial) y deja las unidades como vienen.
// UNA transacción.
func (s *Equipos) Actualizar(ctx context.Context, id model.ID, d DatosEquipo) (*model.Equipo, error) {
	e := d.modelo(id)
	ahora, hoy := time.Now().UTC(), model.Hoy()
	e.ActualizadoEn = ahora
	if err := validarUnidades(d.Unidades, hoy); err != nil {
		return nil, err
	}
	var out *model.Equipo
	err := enTx(ctx, s.db, func(tx *sql.Tx) error {
		repo := repository.NewEquipos(tx)
		actual, err := repo.PorID(ctx, id, hoy)
		if err != nil {
			return err
		}
		e.Clave = actual.Clave
		if err := repo.Actualizar(ctx, e); err != nil {
			return err
		}
		if math.Abs(actual.CostoHora-d.CostoHora) >= 0.005 {
			if err := repo.FijarTarifa(ctx, id, d.CostoHora, hoy, ahora); err != nil {
				return err
			}
		}
		// Sin "unidades" en la petición (nil) las unidades no se tocan.
		if d.Unidades != nil {
			if err := guardarUnidades(ctx, tx, e, d.Unidades, actual.Unidades, ahora); err != nil {
				return err
			}
		}
		out, err = repo.PorID(ctx, id, hoy)
		return err
	})
	return out, err
}

// Baja lógica del equipo y de sus unidades: las claves siguen reservadas.
func (s *Equipos) Baja(ctx context.Context, id model.ID) error {
	ahora := time.Now().UTC()
	return enTx(ctx, s.db, func(tx *sql.Tx) error {
		repo := repository.NewEquipos(tx)
		if err := repo.Baja(ctx, id, ahora); err != nil {
			return err
		}
		return repo.BajaUnidades(ctx, id, nil, ahora)
	})
}

func campoUnidad(i int, campo string) string { return fmt.Sprintf("unidades[%d].%s", i, campo) }

// validarUnidades revisa lo que no depende de la base: fecha no futura y
// series sin repetir dentro de la misma captura.
func validarUnidades(us []DatosUnidad, hoy time.Time) error {
	series := map[string]bool{}
	for i, u := range us {
		if u.FechaAlta.After(hoy) {
			return ErrCampo{campoUnidad(i, "fechaAlta"), "La fecha de alta no puede ser futura"}
		}
		serie := strings.ToUpper(strings.TrimSpace(u.NumSerie))
		if serie == "" {
			continue
		}
		if series[serie] {
			return ErrCampo{campoUnidad(i, "numSerie"), "Este número de serie ya está en otra unidad"}
		}
		series[serie] = true
	}
	return nil
}

// guardarUnidades deja las unidades del equipo como vienen: da de baja las
// guardadas que ya no están, actualiza las que traen ID y crea las nuevas con
// la clave siguiente del equipo (EQ-001-01). Va dentro de la transacción.
func guardarUnidades(ctx context.Context, tx *sql.Tx, e model.Equipo, us []DatosUnidad, actuales []model.Unidad, ahora time.Time) error {
	repo := repository.NewEquipos(tx)
	pedidas := map[model.ID]bool{}
	for i, u := range us {
		if u.ID != nil {
			if pedidas[*u.ID] {
				return ErrCampo{campoUnidad(i, "id"), "La unidad viene repetida"}
			}
			pedidas[*u.ID] = true
		}
	}
	// Primero las bajas, para que sus series queden libres.
	var bajas []model.ID
	existe := map[model.ID]bool{}
	for _, a := range actuales {
		existe[a.ID] = true
		if !pedidas[a.ID] {
			bajas = append(bajas, a.ID)
		}
	}
	if len(bajas) > 0 {
		if err := repo.BajaUnidades(ctx, e.ID, bajas, ahora); err != nil {
			return err
		}
	}
	for i, d := range us {
		u := model.Unidad{NumSerie: strings.ToUpper(strings.TrimSpace(d.NumSerie)), UsosPrevios: d.UsosPrevios}
		var err error
		if d.ID != nil {
			if !existe[*d.ID] {
				return ErrCampo{campoUnidad(i, "id"), "La unidad ya no existe; recarga la página"}
			}
			u.ID = *d.ID
			err = repo.ActualizarUnidad(ctx, e.ID, u, d.FechaAlta, ahora)
		} else {
			var n int
			if n, err = repository.Siguiente(ctx, tx, "UNIDAD:"+e.Clave); err != nil {
				return err
			}
			u.ID, u.Clave = model.NuevoID(), fmt.Sprintf("%s-%02d", e.Clave, n)
			err = repo.CrearUnidad(ctx, e.ID, u, d.FechaAlta, ahora)
		}
		var dup repository.ErrDuplicado
		if errors.As(err, &dup) && dup.Campo == "numSerie" {
			return ErrCampo{campoUnidad(i, "numSerie"), "Ya existe una unidad con este número de serie"}
		}
		if err != nil {
			return err
		}
	}
	return nil
}
