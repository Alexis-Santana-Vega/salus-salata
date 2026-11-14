package handler

import (
	"github.com/gofiber/fiber/v3"

	"salus/backend/internal/model"
	"salus/backend/internal/response"
	"salus/backend/internal/service"
)

type DoctoresHandler struct{ Svc *service.Doctores }

// direccionReq valida el bloque de dirección (lo reusa hospital).
type direccionReq struct {
	Calle       string `json:"calle" validate:"max=120,texto"`
	NumExterior string `json:"numExterior" validate:"max=20,texto"`
	NumInterior string `json:"numInterior" validate:"max=20,texto"`
	Colonia     string `json:"colonia" validate:"max=120,texto"`
	CP          string `json:"cp" validate:"omitempty,cp"`
	Municipio   string `json:"municipio" validate:"max=120,texto"`
	Estado      string `json:"estado" validate:"omitempty,estado"`
}

func (r direccionReq) modelo() model.Direccion {
	return model.Direccion(r)
}

type doctorReq struct {
	Numero         string       `json:"numero" validate:"max=20,texto"`
	Nombre         string       `json:"nombre" validate:"required,max=120,texto"`
	Correo         string       `json:"correo" validate:"required,max=254,email"`
	TelPersonal    string       `json:"telPersonal" validate:"omitempty,telefono"`
	TelConsultorio string       `json:"telConsultorio" validate:"omitempty,telefono"`
	Direccion      direccionReq `json:"direccion"`
}

func (r doctorReq) datos() service.DatosDoctor {
	return service.DatosDoctor{Numero: r.Numero, Nombre: r.Nombre, Correo: r.Correo,
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
