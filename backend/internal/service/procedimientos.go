package service

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"salus/backend/internal/model"
	"salus/backend/internal/repository"
)

type Procedimientos struct{ db *sql.DB }

func NewProcedimientos(db *sql.DB) *Procedimientos { return &Procedimientos{db: db} }

// DatosProcedimiento es lo que se puede capturar en alta y edición.
type DatosProcedimiento struct{ Nombre string }

func (s *Procedimientos) Listar(ctx context.Context, q repository.Consulta) ([]model.Procedimiento, int, error) {
	return repository.NewProcedimientos(s.db).Listar(ctx, q)
}

func (s *Procedimientos) Obtener(ctx context.Context, id model.ID) (*model.Procedimiento, error) {
	return repository.NewProcedimientos(s.db).PorID(ctx, id)
}

func (s *Procedimientos) Crear(ctx context.Context, d DatosProcedimiento) (*model.Procedimiento, error) {
	ahora := time.Now().UTC()
	p := model.Procedimiento{ID: model.NuevoID(), Nombre: limpiar(d.Nombre), CreadoEn: ahora, ActualizadoEn: ahora}
	if err := repository.NewProcedimientos(s.db).Crear(ctx, p); err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *Procedimientos) Actualizar(ctx context.Context, id model.ID, d DatosProcedimiento) (*model.Procedimiento, error) {
	// Actualizar y releer en la misma transacción: se devuelve justo lo guardado.
	var p *model.Procedimiento
	err := enTx(ctx, s.db, func(tx *sql.Tx) error {
		repo := repository.NewProcedimientos(tx)
		if err := repo.Actualizar(ctx, model.Procedimiento{ID: id, Nombre: limpiar(d.Nombre), ActualizadoEn: time.Now().UTC()}); err != nil {
			return err
		}
		var err error
		p, err = repo.PorID(ctx, id)
		return err
	})
	return p, err
}

func (s *Procedimientos) Baja(ctx context.Context, id model.ID) error {
	return repository.NewProcedimientos(s.db).Baja(ctx, id, time.Now().UTC())
}

// limpiar quita espacios de los extremos y deja uno solo entre palabras.
func limpiar(s string) string { return strings.Join(strings.Fields(s), " ") }
