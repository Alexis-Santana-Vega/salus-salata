package service

import (
	"context"
	"database/sql"
	"time"

	"salus/backend/internal/model"
	"salus/backend/internal/repository"
)

type Tecnicos struct{ db *sql.DB }

func NewTecnicos(db *sql.DB) *Tecnicos { return &Tecnicos{db: db} }

// DatosTecnico es lo que se captura en alta y edición.
type DatosTecnico struct{ Nombre, Puesto, Telefono string }

func (d DatosTecnico) limpio() DatosTecnico {
	return DatosTecnico{Nombre: limpiar(d.Nombre), Puesto: d.Puesto, Telefono: limpiar(d.Telefono)}
}

func (s *Tecnicos) Listar(ctx context.Context, q repository.Consulta) ([]model.Tecnico, int, error) {
	return repository.NewTecnicos(s.db).Listar(ctx, q)
}

func (s *Tecnicos) Obtener(ctx context.Context, id model.ID) (*model.Tecnico, error) {
	return repository.NewTecnicos(s.db).PorID(ctx, id)
}

func (s *Tecnicos) Crear(ctx context.Context, d DatosTecnico) (*model.Tecnico, error) {
	d, ahora := d.limpio(), time.Now().UTC()
	t := model.Tecnico{ID: model.NuevoID(), Nombre: d.Nombre, Puesto: d.Puesto, Telefono: d.Telefono, CreadoEn: ahora, ActualizadoEn: ahora}
	if err := repository.NewTecnicos(s.db).Crear(ctx, t); err != nil {
		return nil, err
	}
	return &t, nil
}

// Actualizar y releer en la misma transacción: se devuelve justo lo guardado.
func (s *Tecnicos) Actualizar(ctx context.Context, id model.ID, d DatosTecnico) (*model.Tecnico, error) {
	d = d.limpio()
	var t *model.Tecnico
	err := enTx(ctx, s.db, func(tx *sql.Tx) error {
		repo := repository.NewTecnicos(tx)
		if err := repo.Actualizar(ctx, model.Tecnico{ID: id, Nombre: d.Nombre, Puesto: d.Puesto, Telefono: d.Telefono, ActualizadoEn: time.Now().UTC()}); err != nil {
			return err
		}
		var err error
		t, err = repo.PorID(ctx, id)
		return err
	})
	return t, err
}

func (s *Tecnicos) Baja(ctx context.Context, id model.ID) error {
	return repository.NewTecnicos(s.db).Baja(ctx, id, time.Now().UTC())
}
