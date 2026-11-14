package handler

import (
	"github.com/gofiber/fiber/v3"

	"salus/backend/internal/model"
	"salus/backend/internal/response"
	"salus/backend/internal/service"
)

type DireccionesHandler struct{ Svc *service.Direcciones }

// codigoPostalResp: las colonias de un C.P. con su municipio y estado.
type codigoPostalResp struct {
	CP            string               `json:"cp"`
	Asentamientos []model.Asentamiento `json:"asentamientos"`
}

// PorCP: GET /api/v1/codigos-postales/{cp}. 404 si el C.P. no está en el catálogo.
func (h DireccionesHandler) PorCP(c fiber.Ctx) error {
	cp := c.Params("cp")
	if !reCP.MatchString(cp) {
		return response.Validation(response.FieldError{Field: "cp", Message: "El código postal tiene 5 dígitos"})
	}
	as, err := h.Svc.PorCP(c.Context(), cp)
	if err != nil {
		return err
	}
	if len(as) == 0 {
		return response.NotFound("No encontramos ese código postal en el catálogo")
	}
	return response.OK(c, codigoPostalResp{CP: cp, Asentamientos: as})
}
