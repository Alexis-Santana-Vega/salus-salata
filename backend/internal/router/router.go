// Package router declara todas las rutas de la API en un solo lugar.
package router

import (
	"database/sql"

	"github.com/gofiber/fiber/v3"

	"salus/backend/internal/config"
	"salus/backend/internal/docs"
	"salus/backend/internal/handler"
	"salus/backend/internal/middleware"
	"salus/backend/internal/model"
	"salus/backend/internal/response"
	"salus/backend/internal/security"
	"salus/backend/internal/service"
)

type Deps struct {
	Config   config.Config
	DB       *sql.DB
	Version  string
	Tokens   *security.Tokens
	Auth     *service.Auth
	Usuarios *service.Usuarios

	Procedimientos *service.Procedimientos
	Tecnicos       *service.Tecnicos
	Doctores       *service.Doctores
}

func Setup(app *fiber.App, d Deps) {
	if d.Config.DocsEnabled {
		docs.Register(app)
	}

	v1 := app.Group("/api/v1")

	health := handler.HealthHandler{DB: d.DB, Version: d.Version}
	v1.Get("/health", health.Get)

	// --- Autenticación ---
	auth := handler.AuthHandler{Svc: d.Auth, CookieSegura: d.Config.IsProduction()}
	sesion := middleware.RequireAuth(d.Tokens, false)         // sesión normal
	sesionPendiente := middleware.RequireAuth(d.Tokens, true) // admite "debe cambiar contraseña"

	v1.Post("/auth/login", middleware.LoginLimiter(), auth.Login)
	v1.Post("/auth/refresh", auth.Refresh)
	v1.Post("/auth/logout", auth.Logout)
	v1.Get("/auth/me", sesionPendiente, auth.Me)
	v1.Put("/auth/contrasena", middleware.LoginLimiter(), sesionPendiente, auth.CambiarContrasena)

	// --- Usuarios (solo Sysadmin) ---
	usuarios := handler.UsuariosHandler{Svc: d.Usuarios}
	v1.Get("/usuarios", sesion, middleware.RequireRole(model.RolSysadmin), usuarios.Listar)

	// --- Catálogos (Sysadmin y Admin) ---
	// Cada catálogo es un grupo con su propia sesión y rol: así una ruta
	// inexistente sigue respondiendo 404 y no 401.
	admin := []any{sesion, middleware.RequireRole(model.RolSysadmin, model.RolAdmin)}
	crud(v1, "/procedimientos", admin, handler.ProcedimientosHandler{Svc: d.Procedimientos})
	crud(v1, "/tecnicos", admin, handler.TecnicosHandler{Svc: d.Tecnicos})
	crud(v1, "/doctores", admin, handler.DoctoresHandler{Svc: d.Doctores})

	// Cualquier otra ruta responde 404 con el formato único.
	app.Use(func(c fiber.Ctx) error {
		return response.NotFound("La ruta " + c.Method() + " " + c.Path() + " no existe")
	})
}

// Recurso es lo que implementa cada catálogo con baja lógica.
type recurso interface {
	Listar(fiber.Ctx) error
	Obtener(fiber.Ctx) error
	Crear(fiber.Ctx) error
	Actualizar(fiber.Ctx) error
	Baja(fiber.Ctx) error
}

// crud registra las cinco rutas estándar de un catálogo bajo su grupo.
func crud(r fiber.Router, ruta string, mw []any, h recurso) {
	g := r.Group(ruta, mw...)
	g.Get("", h.Listar)
	g.Post("", h.Crear)
	g.Get("/:id", h.Obtener)
	g.Put("/:id", h.Actualizar)
	g.Delete("/:id", h.Baja)
}
