package handler

import (
	"github.com/gofiber/fiber/v3"

	"salus/backend/internal/response"
	"salus/backend/internal/service"
)

type TecnicosHandler struct{ Svc *service.Tecnicos }

type tecnicoReq struct {
	Nombre   string `json:"nombre" validate:"required,max=120,nombre"`
	Puesto   string `json:"puesto" validate:"required,oneof=AYUDANTE ENFERMERIA ANESTESIOLOGIA LOGISTICA OTRO"`
	Telefono string `json:"telefono" validate:"omitempty,telefono"`
}

func (r tecnicoReq) datos() service.DatosTecnico {
	return service.DatosTecnico{Nombre: r.Nombre, Puesto: r.Puesto, Telefono: r.Telefono}
}

// Listar: GET /api/v1/tecnicos?page=&pageSize=&q=&sort=
func (h TecnicosHandler) Listar(c fiber.Ctx) error {
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

// Obtener: GET /api/v1/tecnicos/{id}
func (h TecnicosHandler) Obtener(c fiber.Ctx) error {
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

// Crear: POST /api/v1/tecnicos
func (h TecnicosHandler) Crear(c fiber.Ctx) error {
	var req tecnicoReq
	if err := bind(c, &req); err != nil {
		return err
	}
	p, err := h.Svc.Crear(c.Context(), req.datos())
	if err != nil {
		return errorCatalogo(err)
	}
	return response.Created(c, p)
}

// Actualizar: PUT /api/v1/tecnicos/{id}
func (h TecnicosHandler) Actualizar(c fiber.Ctx) error {
	id, err := idParam(c)
	if err != nil {
		return err
	}
	var req tecnicoReq
	if err := bind(c, &req); err != nil {
		return err
	}
	p, err := h.Svc.Actualizar(c.Context(), id, req.datos())
	if err != nil {
		return errorCatalogo(err)
	}
	return response.OK(c, p)
}

// Baja: DELETE /api/v1/tecnicos/{id} (baja lógica).
func (h TecnicosHandler) Baja(c fiber.Ctx) error {
	id, err := idParam(c)
	if err != nil {
		return err
	}
	if err := h.Svc.Baja(c.Context(), id); err != nil {
		return errorCatalogo(err)
	}
	return response.NoContent(c)
}
