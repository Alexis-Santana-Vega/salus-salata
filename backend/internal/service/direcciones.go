package service

import (
	"context"
	"database/sql"

	"salus/backend/internal/model"
	"salus/backend/internal/repository"
)

type Direcciones struct{ db *sql.DB }

func NewDirecciones(db *sql.DB) *Direcciones { return &Direcciones{db: db} }

// PorCP devuelve las colonias de un código postal; lista vacía si no está en el catálogo.
func (s *Direcciones) PorCP(ctx context.Context, cp string) ([]model.Asentamiento, error) {
	return repository.NewDirecciones(s.db).PorCP(ctx, cp)
}
