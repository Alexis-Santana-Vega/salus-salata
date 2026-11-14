package handler

import (
	"github.com/gofiber/fiber/v3"

	"salus/backend/internal/response"
	"salus/backend/internal/service"
)

type InsumosHandler struct{ Svc *service.Insumos }

type insumoReq struct {
	Descripcion  string   `json:"descripcion" validate:"required,max=200,textogeneral"`
	Unidad       string   `json:"unidad" validate:"required,oneof=PIEZA CAJA PAQUETE"`
	Precio       *float64 `json:"precio" validate:"required,gte=0,lte=9999999999,dosdecimales"`
	CodigoBarras string   `json:"codigoBarras" validate:"omitempty,max=64,codigo"`
	StockActual  *int     `json:"stockActual" validate:"required,gte=0,lte=1000000"`
	StockMinimo  *int     `json:"stockMinimo" validate:"required,gte=0,lte=1000000"`
}

func (r insumoReq) datos() service.DatosInsumo {
	return service.DatosInsumo{Descripcion: r.Descripcion, Unidad: r.Unidad, Precio: *r.Precio,
		CodigoBarras: r.CodigoBarras, StockActual: *r.StockActual, StockMinimo: *r.StockMinimo}
}

// Listar: GET /api/v1/insumos?page=&pageSize=&q=&sort=
func (h InsumosHandler) Listar(c fiber.Ctx) error {
	q, meta, err := consulta(c)
	if err != nil {
		return err
	}
	items, total, err := h.Svc.Listar(c.Context(), q)
	if err != nil {
		return errorCatalogo(err)
	}
	meta.Total = total
	return response.List(c, items, meta)
}

// Obtener: GET /api/v1/insumos/{id}
func (h InsumosHandler) Obtener(c fiber.Ctx) error {
	id, err := idParam(c)
	if err != nil {
		return err
	}
	p, err := h.Svc.Obtener(c.Context(), id)
	if err != nil {
		return errorCatalogo(err)
	}
	return response.OK(c, p)
}

// Crear: POST /api/v1/insumos
func (h InsumosHandler) Crear(c fiber.Ctx) error {
	var req insumoReq
	if err := bind(c, &req); err != nil {
		return err
	}
	p, err := h.Svc.Crear(c.Context(), req.datos())
	if err != nil {
		return errorCatalogo(err)
	}
	return response.Created(c, p)
}

// Actualizar: PUT /api/v1/insumos/{id}
func (h InsumosHandler) Actualizar(c fiber.Ctx) error {
	id, err := idParam(c)
	if err != nil {
		return err
	}
	var req insumoReq
	if err := bind(c, &req); err != nil {
		return err
	}
	p, err := h.Svc.Actualizar(c.Context(), id, req.datos())
	if err != nil {
		return errorCatalogo(err)
	}
	return response.OK(c, p)
}

// Baja: DELETE /api/v1/insumos/{id} (baja lógica).
func (h InsumosHandler) Baja(c fiber.Ctx) error {
	id, err := idParam(c)
	if err != nil {
		return err
	}
	if err := h.Svc.Baja(c.Context(), id); err != nil {
		return errorCatalogo(err)
	}
	return response.NoContent(c)
}
