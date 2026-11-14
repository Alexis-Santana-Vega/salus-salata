package handler

import (
	"github.com/gofiber/fiber/v3"

	"salus/backend/internal/response"
	"salus/backend/internal/service"
)

type EquiposHandler struct{ Svc *service.Equipos }

type equipoReq struct {
	Nombre        string   `json:"nombre" validate:"required,max=120,textogeneral"`
	MarcaModelo   string   `json:"marcaModelo" validate:"required,max=120,textogeneral"`
	CicloMax      *int     `json:"cicloMax" validate:"required,gte=1,lte=100000"`
	RequierePurga bool     `json:"requierePurga"`
	PurgaCada     *int     `json:"purgaCada" validate:"required_if=RequierePurga true,omitempty,gte=1,lte=100000"`
	CostoHora     *float64 `json:"costoHora" validate:"required,gte=0,lte=9999999999,dosdecimales"`
}

func (r equipoReq) datos() service.DatosEquipo {
	return service.DatosEquipo{Nombre: r.Nombre, MarcaModelo: r.MarcaModelo, CicloMax: *r.CicloMax,
		RequierePurga: r.RequierePurga, PurgaCada: r.PurgaCada, CostoHora: *r.CostoHora}
}

// Listar: GET /api/v1/equipos?page=&pageSize=&q=&sort=
func (h EquiposHandler) Listar(c fiber.Ctx) error {
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

// Obtener: GET /api/v1/equipos/{id}
func (h EquiposHandler) Obtener(c fiber.Ctx) error {
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

// Crear: POST /api/v1/equipos
func (h EquiposHandler) Crear(c fiber.Ctx) error {
	var req equipoReq
	if err := bind(c, &req); err != nil {
		return err
	}
	p, err := h.Svc.Crear(c.Context(), req.datos())
	if err != nil {
		return errorCatalogo(err)
	}
	return response.Created(c, p)
}

// Actualizar: PUT /api/v1/equipos/{id}
func (h EquiposHandler) Actualizar(c fiber.Ctx) error {
	id, err := idParam(c)
	if err != nil {
		return err
	}
	var req equipoReq
	if err := bind(c, &req); err != nil {
		return err
	}
	p, err := h.Svc.Actualizar(c.Context(), id, req.datos())
	if err != nil {
		return errorCatalogo(err)
	}
	return response.OK(c, p)
}

// Baja: DELETE /api/v1/equipos/{id} (baja lógica).
func (h EquiposHandler) Baja(c fiber.Ctx) error {
	id, err := idParam(c)
	if err != nil {
		return err
	}
	if err := h.Svc.Baja(c.Context(), id); err != nil {
		return errorCatalogo(err)
	}
	return response.NoContent(c)
}
