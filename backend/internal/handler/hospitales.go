package handler

import (
	"github.com/gofiber/fiber/v3"

	"salus/backend/internal/model"
	"salus/backend/internal/response"
	"salus/backend/internal/service"
)

type HospitalesHandler struct{ Svc *service.Hospitales }

// La dirección del hospital es obligatoria (calle y colonia del catálogo).
type direccionHospitalReq struct {
	Calle          string    `json:"calle" validate:"required,max=120,textogeneral"`
	NumExterior    string    `json:"numExterior" validate:"omitempty,max=20,numdomicilio"`
	NumInterior    string    `json:"numInterior" validate:"omitempty,max=20,numdomicilio"`
	AsentamientoID *model.ID `json:"asentamientoId" validate:"required"`
}

type contactoReq struct {
	Nombre   string `json:"nombre" validate:"required,max=120,nombre"`
	Puesto   string `json:"puesto" validate:"omitempty,max=60,textogeneral"`
	Telefono string `json:"telefono" validate:"omitempty,telefono"`
}

type hospitalReq struct {
	Nombre    string               `json:"nombre" validate:"required,max=150,textogeneral"`
	Telefono  string               `json:"telefono" validate:"omitempty,telefono"`
	Direccion direccionHospitalReq `json:"direccion"`
	Latitud   *float64             `json:"latitud" validate:"omitempty,gte=-90,lte=90"`
	Longitud  *float64             `json:"longitud" validate:"omitempty,gte=-180,lte=180"`
	LigaMapa  string               `json:"ligaMapa" validate:"omitempty,max=500,url,startswith=https://"`
	Contactos []contactoReq        `json:"contactos" validate:"max=20,dive"`
}

func (r hospitalReq) datos() service.DatosHospital {
	cs := make([]model.Contacto, len(r.Contactos))
	for i, c := range r.Contactos {
		cs[i] = model.Contacto(c)
	}
	d := r.Direccion
	return service.DatosHospital{Nombre: r.Nombre, Telefono: r.Telefono, LigaMapa: r.LigaMapa,
		Direccion: model.Direccion{Calle: d.Calle, NumExterior: d.NumExterior, NumInterior: d.NumInterior, AsentamientoID: d.AsentamientoID},
		Latitud:   r.Latitud, Longitud: r.Longitud, Contactos: cs}
}

// Listar: GET /api/v1/hospitales?page=&pageSize=&q=&sort=
func (h HospitalesHandler) Listar(c fiber.Ctx) error {
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

// Obtener: GET /api/v1/hospitales/{id}
func (h HospitalesHandler) Obtener(c fiber.Ctx) error {
	id, err := idParam(c)
	if err != nil {
		return err
	}
	x, err := h.Svc.Obtener(c.Context(), id)
	if err != nil {
		return errorCatalogo(err)
	}
	return response.OK(c, x)
}

// Crear: POST /api/v1/hospitales
func (h HospitalesHandler) Crear(c fiber.Ctx) error {
	var req hospitalReq
	if err := bind(c, &req); err != nil {
		return err
	}
	x, err := h.Svc.Crear(c.Context(), req.datos())
	if err != nil {
		return errorCatalogo(err)
	}
	return response.Created(c, x)
}

// Actualizar: PUT /api/v1/hospitales/{id}
func (h HospitalesHandler) Actualizar(c fiber.Ctx) error {
	id, err := idParam(c)
	if err != nil {
		return err
	}
	var req hospitalReq
	if err := bind(c, &req); err != nil {
		return err
	}
	x, err := h.Svc.Actualizar(c.Context(), id, req.datos())
	if err != nil {
		return errorCatalogo(err)
	}
	return response.OK(c, x)
}

// Baja: DELETE /api/v1/hospitales/{id} (baja lógica).
func (h HospitalesHandler) Baja(c fiber.Ctx) error {
	id, err := idParam(c)
	if err != nil {
		return err
	}
	if err := h.Svc.Baja(c.Context(), id); err != nil {
		return errorCatalogo(err)
	}
	return response.NoContent(c)
}
