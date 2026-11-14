package handler

import (
	"github.com/gofiber/fiber/v3"

	"salus/backend/internal/response"
	"salus/backend/internal/service"
)

type ProcedimientosHandler struct{ Svc *service.Procedimientos }

type procedimientoReq struct {
	Nombre string `json:"nombre" validate:"required,max=120,texto"`
}

func (r procedimientoReq) datos() service.DatosProcedimiento {
	return service.DatosProcedimiento{Nombre: r.Nombre}
}

// Listar: GET /api/v1/procedimientos?page=&pageSize=&q=&sort=
func (h ProcedimientosHandler) Listar(c fiber.Ctx) error {
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

// Obtener: GET /api/v1/procedimientos/{id}
func (h ProcedimientosHandler) Obtener(c fiber.Ctx) error {
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

// Crear: POST /api/v1/procedimientos
func (h ProcedimientosHandler) Crear(c fiber.Ctx) error {
	var req procedimientoReq
	if err := bind(c, &req); err != nil {
		return err
	}
	p, err := h.Svc.Crear(c.Context(), req.datos())
	if err != nil {
		return errorCatalogo(err)
	}
	return response.Created(c, p)
}

// Actualizar: PUT /api/v1/procedimientos/{id}
func (h ProcedimientosHandler) Actualizar(c fiber.Ctx) error {
	id, err := idParam(c)
	if err != nil {
		return err
	}
	var req procedimientoReq
	if err := bind(c, &req); err != nil {
		return err
	}
	p, err := h.Svc.Actualizar(c.Context(), id, req.datos())
	if err != nil {
		return errorCatalogo(err)
	}
	return response.OK(c, p)
}

// Baja: DELETE /api/v1/procedimientos/{id} (baja lógica).
func (h ProcedimientosHandler) Baja(c fiber.Ctx) error {
	id, err := idParam(c)
	if err != nil {
		return err
	}
	if err := h.Svc.Baja(c.Context(), id); err != nil {
		return errorCatalogo(err)
	}
	return response.NoContent(c)
}
