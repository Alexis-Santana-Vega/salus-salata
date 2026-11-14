package middleware

import (
	"slices"
	"strings"

	"github.com/gofiber/fiber/v3"

	"salus/backend/internal/model"
	"salus/backend/internal/response"
	"salus/backend/internal/security"
)

const claveSesion = "sesion"

// SesionActual es lo que los handlers saben del usuario autenticado.
type SesionActual struct {
	UsuarioID   model.ID
	Rol         string
	DebeCambiar bool
}

// Sesion devuelve la sesión guardada por RequireAuth.
func Sesion(c fiber.Ctx) SesionActual {
	s, _ := c.Locals(claveSesion).(SesionActual)
	return s
}

// RequireAuth exige "Authorization: Bearer <token>" válido. Si el usuario
// debe cambiar su contraseña, solo deja pasar las rutas marcadas con
// permitirCambioPendiente (ver, cambiar contraseña y salir).
func RequireAuth(tokens *security.Tokens, permitirCambioPendiente bool) fiber.Handler {
	return func(c fiber.Ctx) error {
		raw, ok := strings.CutPrefix(c.Get(fiber.HeaderAuthorization), "Bearer ")
		if !ok || raw == "" {
			return response.Unauthorized("Inicia sesión para continuar")
		}
		claims, err := tokens.Validar(raw)
		if err != nil {
			return response.Unauthorized("Tu sesión terminó. Inicia sesión de nuevo.")
		}
		id, err := claims.UsuarioID()
		if err != nil {
			return response.Unauthorized("Tu sesión terminó. Inicia sesión de nuevo.")
		}
		if claims.DebeCambiar && !permitirCambioPendiente {
			return response.New(fiber.StatusForbidden, response.CodePasswordChange,
				"Cambia tu contraseña para continuar")
		}
		c.Locals(claveSesion, SesionActual{UsuarioID: id, Rol: claims.Rol, DebeCambiar: claims.DebeCambiar})
		return c.Next()
	}
}

// RequireRole deja pasar solo a los roles indicados. Va después de RequireAuth.
func RequireRole(roles ...string) fiber.Handler {
	return func(c fiber.Ctx) error {
		if !slices.Contains(roles, Sesion(c).Rol) {
			return response.Forbidden("No tienes permiso para esta acción")
		}
		return c.Next()
	}
}
