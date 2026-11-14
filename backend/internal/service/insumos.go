package service

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"salus/backend/internal/model"
	"salus/backend/internal/repository"
)

type Insumos struct{ db *sql.DB }

func NewInsumos(db *sql.DB) *Insumos { return &Insumos{db: db} }

// DatosInsumo es lo que se captura en alta y edición.
type DatosInsumo struct {
	Clave, Descripcion, Unidad, CodigoBarras string
	Precio                                   float64
	StockActual, StockMinimo                 int
}

func (d DatosInsumo) modelo(id model.ID) model.Insumo {
	return model.Insumo{ID: id, Clave: strings.ToUpper(limpiar(d.Clave)), Descripcion: limpiar(d.Descripcion),
		Unidad: d.Unidad, Precio: d.Precio, CodigoBarras: strings.ToUpper(limpiar(d.CodigoBarras)),
		StockActual: d.StockActual, StockMinimo: d.StockMinimo}
}

func (s *Insumos) Listar(ctx context.Context, q repository.Consulta) ([]model.Insumo, int, error) {
	return repository.NewInsumos(s.db).Listar(ctx, q)
}

func (s *Insumos) Obtener(ctx context.Context, id model.ID) (*model.Insumo, error) {
	return repository.NewInsumos(s.db).PorID(ctx, id)
}

// Crear reserva la clave y da de alta el insumo en UNA transacción.
func (s *Insumos) Crear(ctx context.Context, d DatosInsumo) (*model.Insumo, error) {
	i := d.modelo(model.NuevoID())
	i.CreadoEn = time.Now().UTC()
	i.ActualizadoEn = i.CreadoEn
	var out *model.Insumo
	err := enTx(ctx, s.db, func(tx *sql.Tx) error {
		repo := repository.NewInsumos(tx)
		if err := repo.ReservarClave(ctx, i.Clave); err != nil {
			return err
		}
		if err := repo.Crear(ctx, i); err != nil {
			return err
		}
		var err error
		out, err = repo.PorID(ctx, i.ID)
		return err
	})
	return out, err
}

// Actualizar corrige la clave (si cambió) y el resto de los datos en UNA transacción.
func (s *Insumos) Actualizar(ctx context.Context, id model.ID, d DatosInsumo) (*model.Insumo, error) {
	i := d.modelo(id)
	i.ActualizadoEn = time.Now().UTC()
	var out *model.Insumo
	err := enTx(ctx, s.db, func(tx *sql.Tx) error {
		repo := repository.NewInsumos(tx)
		actual, err := repo.ClaveDe(ctx, id)
		if err != nil {
			return err
		}
		if actual != i.Clave {
			if err := repo.CambiarClave(ctx, actual, i.Clave); err != nil {
				return err
			}
		}
		if err := repo.Actualizar(ctx, i); err != nil {
			return err
		}
		out, err = repo.PorID(ctx, id)
		return err
	})
	return out, err
}

// Baja lógica: la clave sigue reservada (no se reusa).
func (s *Insumos) Baja(ctx context.Context, id model.ID) error {
	return repository.NewInsumos(s.db).Baja(ctx, id, time.Now().UTC())
}
