package service

import (
	"context"
	"database/sql"
	"fmt"
	"math"
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

// Crear genera la clave (EQ-001), la reserva, da de alta el modelo y su primera
// tarifa (vigente desde hoy) en UNA transacción.
func (s *Equipos) Crear(ctx context.Context, d DatosEquipo) (*model.Equipo, error) {
	e := d.modelo(model.NuevoID())
	ahora, hoy := time.Now().UTC(), model.Hoy()
	e.CreadoEn, e.ActualizadoEn = ahora, ahora
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
		out, err = repo.PorID(ctx, e.ID, hoy)
		return err
	})
	return out, err
}

// Actualizar cambia el modelo y, solo si el costo cambió, agrega la tarifa
// vigente desde hoy (las anteriores quedan como historial). UNA transacción.
func (s *Equipos) Actualizar(ctx context.Context, id model.ID, d DatosEquipo) (*model.Equipo, error) {
	e := d.modelo(id)
	ahora, hoy := time.Now().UTC(), model.Hoy()
	e.ActualizadoEn = ahora
	var out *model.Equipo
	err := enTx(ctx, s.db, func(tx *sql.Tx) error {
		repo := repository.NewEquipos(tx)
		actual, err := repo.PorID(ctx, id, hoy)
		if err != nil {
			return err
		}
		if err := repo.Actualizar(ctx, e); err != nil {
			return err
		}
		if math.Abs(actual.CostoHora-d.CostoHora) >= 0.005 {
			if err := repo.FijarTarifa(ctx, id, d.CostoHora, hoy, ahora); err != nil {
				return err
			}
		}
		out, err = repo.PorID(ctx, id, hoy)
		return err
	})
	return out, err
}

// Baja lógica: la clave sigue reservada.
func (s *Equipos) Baja(ctx context.Context, id model.ID) error {
	return repository.NewEquipos(s.db).Baja(ctx, id, time.Now().UTC())
}
