package repository

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"time"

	"salus/backend/internal/model"
)

type Insumos struct{ db DBTX }

func NewInsumos(db DBTX) *Insumos { return &Insumos{db: db} }

var (
	ordenInsumos   = map[string]string{"clave": "CLAVE", "descripcion": "DESCRIPCION", "precio": "PRECIO", "stockActual": "STOCK_ACTUAL", "stockMinimo": "STOCK_MINIMO", "creadoEn": "CREADO_EN"}
	indicesInsumos = map[string]string{"PK_CLAVE_CATALOGO": "clave", "UQ_INSUMO_CLAVE": "clave", "UQ_INSUMO_CODIGO_BARRAS": "codigoBarras"}
)

const selectInsumo = `
SELECT ID, CLAVE, DESCRIPCION, UNIDAD, PRECIO, CODIGO_BARRAS, STOCK_ACTUAL, STOCK_MINIMO, CREADO_EN, ACTUALIZADO_EN
FROM INSUMO`

func scanInsumo(s interface{ Scan(...any) error }) (*model.Insumo, error) {
	var i model.Insumo
	var precio string
	var codigo sql.NullString
	err := s.Scan(&i.ID, &i.Clave, &i.Descripcion, &i.Unidad, &precio, &codigo, &i.StockActual, &i.StockMinimo, &i.CreadoEn, &i.ActualizadoEn)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNoEncontrado
	}
	if err != nil {
		return nil, err
	}
	i.Precio, _ = strconv.ParseFloat(precio, 64) // NUMERIC llega como texto
	i.CodigoBarras = codigo.String
	i.BajoMinimo = i.StockActual < i.StockMinimo
	return &i, nil
}

// Listar devuelve una página de insumos activos; busca por clave, descripción o código de barras.
func (r *Insumos) Listar(ctx context.Context, q Consulta) ([]model.Insumo, int, error) {
	orden, err := orderBy(q.Orden, ordenInsumos, "clave", "ID")
	if err != nil {
		return nil, 0, err
	}
	where, args := " WHERE ACTIVO = TRUE", []any{}
	if q.Q != "" {
		where += " AND (CLAVE CONTAINING ? OR DESCRIPCION CONTAINING ? OR CODIGO_BARRAS CONTAINING ?)"
		args = append(args, q.Q, q.Q, q.Q)
	}
	var total int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM INSUMO`+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.db.QueryContext(ctx, selectInsumo+where+orden+` OFFSET ? ROWS FETCH NEXT ? ROWS ONLY`,
		append(args, q.Offset, q.Limite)...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	out := make([]model.Insumo, 0, q.Limite)
	for rows.Next() {
		i, err := scanInsumo(rows)
		if err != nil {
			return nil, 0, err
		}
		out = append(out, *i)
	}
	return out, total, rows.Err()
}

func (r *Insumos) PorID(ctx context.Context, id model.ID) (*model.Insumo, error) {
	return scanInsumo(r.db.QueryRowContext(ctx, selectInsumo+` WHERE ID = ? AND ACTIVO = TRUE`, id))
}

// ClaveDe devuelve la clave actual de un insumo activo.
func (r *Insumos) ClaveDe(ctx context.Context, id model.ID) (string, error) {
	var clave string
	err := r.db.QueryRowContext(ctx, `SELECT CLAVE FROM INSUMO WHERE ID = ? AND ACTIVO = TRUE`, id).Scan(&clave)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrNoEncontrado
	}
	return clave, err
}

// ReservarClave registra la clave en CLAVE_CATALOGO (falla si ya la usa un
// equipo o un insumo, aunque esté dado de baja).
func (r *Insumos) ReservarClave(ctx context.Context, clave string) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO CLAVE_CATALOGO (CLAVE, TIPO) VALUES (?, 'INSUMO')`, clave)
	return integridad(err, indicesInsumos, nil)
}

// CambiarClave corrige la clave: se cambia en CLAVE_CATALOGO y baja en
// cascada al insumo. La anterior queda libre.
func (r *Insumos) CambiarClave(ctx context.Context, anterior, nueva string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE CLAVE_CATALOGO SET CLAVE = ? WHERE CLAVE = ? AND TIPO = 'INSUMO'`, nueva, anterior)
	return integridad(err, indicesInsumos, nil)
}

func precioSQL(p float64) string { return strconv.FormatFloat(p, 'f', 2, 64) }

func (r *Insumos) Crear(ctx context.Context, i model.Insumo) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO INSUMO (ID, CLAVE, DESCRIPCION, UNIDAD, PRECIO, CODIGO_BARRAS, STOCK_ACTUAL, STOCK_MINIMO, CREADO_EN, ACTUALIZADO_EN)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		i.ID, i.Clave, i.Descripcion, i.Unidad, precioSQL(i.Precio), nulo(i.CodigoBarras), i.StockActual, i.StockMinimo, i.CreadoEn, i.ActualizadoEn)
	return integridad(err, indicesInsumos, nil)
}

// Actualizar cambia todo menos la clave (esa va por CambiarClave).
func (r *Insumos) Actualizar(ctx context.Context, i model.Insumo) error {
	res, err := r.db.ExecContext(ctx, `
		UPDATE INSUMO SET DESCRIPCION = ?, UNIDAD = ?, PRECIO = ?, CODIGO_BARRAS = ?, STOCK_ACTUAL = ?, STOCK_MINIMO = ?, ACTUALIZADO_EN = ?
		WHERE ID = ? AND ACTIVO = TRUE`,
		i.Descripcion, i.Unidad, precioSQL(i.Precio), nulo(i.CodigoBarras), i.StockActual, i.StockMinimo, i.ActualizadoEn, i.ID)
	if err != nil {
		return integridad(err, indicesInsumos, nil)
	}
	return unaFila(res)
}

func (r *Insumos) Baja(ctx context.Context, id model.ID, ahora time.Time) error {
	res, err := r.db.ExecContext(ctx, `
		UPDATE INSUMO SET ACTIVO = FALSE, ACTUALIZADO_EN = ? WHERE ID = ? AND ACTIVO = TRUE`, ahora, id)
	if err != nil {
		return err
	}
	return unaFila(res)
}
