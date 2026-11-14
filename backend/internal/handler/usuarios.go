package handler

import (
	"github.com/gofiber/fiber/v3"

	"salus/backend/internal/response"
	"salus/backend/internal/service"
)

type UsuariosHandler struct{ Svc *service.Usuarios }

// Listar: GET /api/v1/usuarios?page=&pageSize= (solo Sysadmin).
func (h UsuariosHandler) Listar(c fiber.Ctx) error {
	page, size, offset, err := paginacion(c)
	if err != nil {
		return err
	}
	usuarios, total, err := h.Svc.Listar(c.Context(), offset, size)
	if err != nil {
		return err
	}
	return response.List(c, usuarios, response.Meta{Page: page, PageSize: size, Total: total})
}
