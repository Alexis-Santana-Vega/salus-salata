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
	v1.Put("/auth/contrasena", sesionPendiente, auth.CambiarContrasena)

	// --- Usuarios (solo Sysadmin) ---
	usuarios := handler.UsuariosHandler{Svc: d.Usuarios}
	v1.Get("/usuarios", sesion, middleware.RequireRole(model.RolSysadmin), usuarios.Listar)

	// Cualquier otra ruta responde 404 con el formato único.
	app.Use(func(c fiber.Ctx) error {
		return response.NotFound("La ruta " + c.Method() + " " + c.Path() + " no existe")
	})
}
