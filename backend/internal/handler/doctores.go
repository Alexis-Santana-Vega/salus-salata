package handler

import (
	"github.com/gofiber/fiber/v3"

	"salus/backend/internal/model"
	"salus/backend/internal/response"
	"salus/backend/internal/service"
)

type DoctoresHandler struct{ Svc *service.Doctores }

// direccionReq: lo que se captura de una dirección (lo reusa hospital). La
// colonia se elige del catálogo por su id; lo demás sale del catálogo.
type direccionReq struct {
	Calle          string    `json:"calle" validate:"omitempty,max=120,textogeneral"`
	NumExterior    string    `json:"numExterior" validate:"omitempty,max=20,numdomicilio"`
	NumInterior    string    `json:"numInterior" validate:"omitempty,max=20,numdomicilio"`
	AsentamientoID *model.ID `json:"asentamientoId"`
}

func (r direccionReq) modelo() model.Direccion {
	return model.Direccion{Calle: r.Calle, NumExterior: r.NumExterior, NumInterior: r.NumInterior, AsentamientoID: r.AsentamientoID}
}

type doctorReq struct {
	Nombre         string       `json:"nombre" validate:"required,max=120,nombre"`
	Correo         string       `json:"correo" validate:"required,max=254,email"`
	TelPersonal    string       `json:"telPersonal" validate:"omitempty,telefono"`
	TelConsultorio string       `json:"telConsultorio" validate:"omitempty,telefono"`
	Direccion      direccionReq `json:"direccion"`
}

func (r doctorReq) datos() service.DatosDoctor {
	return service.DatosDoctor{Nombre: r.Nombre, Correo: r.Correo,
		TelPersonal: r.TelPersonal, TelConsultorio: r.TelConsultorio, Direccion: r.Direccion.modelo()}
}

// Listar: GET /api/v1/doctores?page=&pageSize=&q=&sort=
func (h DoctoresHandler) Listar(c fiber.Ctx) error {
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

// Obtener: GET /api/v1/doctores/{id}
func (h DoctoresHandler) Obtener(c fiber.Ctx) error {
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

// Crear: POST /api/v1/doctores
func (h DoctoresHandler) Crear(c fiber.Ctx) error {
	var req doctorReq
	if err := bind(c, &req); err != nil {
		return err
	}
	p, err := h.Svc.Crear(c.Context(), req.datos())
	if err != nil {
		return errorCatalogo(err)
	}
	return response.Created(c, p)
}

// Actualizar: PUT /api/v1/doctores/{id}
func (h DoctoresHandler) Actualizar(c fiber.Ctx) error {
	id, err := idParam(c)
	if err != nil {
		return err
	}
	var req doctorReq
	if err := bind(c, &req); err != nil {
		return err
	}
	p, err := h.Svc.Actualizar(c.Context(), id, req.datos())
	if err != nil {
		return errorCatalogo(err)
	}
	return response.OK(c, p)
}

// Baja: DELETE /api/v1/doctores/{id} (baja lógica).
func (h DoctoresHandler) Baja(c fiber.Ctx) error {
	id, err := idParam(c)
	if err != nil {
		return err
	}
	if err := h.Svc.Baja(c.Context(), id); err != nil {
		return errorCatalogo(err)
	}
	return response.NoContent(c)
}
