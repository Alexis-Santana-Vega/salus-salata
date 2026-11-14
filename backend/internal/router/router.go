// Package router declara todas las rutas de la API en un solo lugar.
package router

import (
	"database/sql"

	"github.com/gofiber/fiber/v3"

	"salus/backend/internal/config"
	"salus/backend/internal/docs"
	"salus/backend/internal/handler"
	"salus/backend/internal/response"
)

type Deps struct {
	Config  config.Config
	DB      *sql.DB
	Version string
}

func Setup(app *fiber.App, d Deps) {
	if d.Config.DocsEnabled {
		docs.Register(app)
	}

	v1 := app.Group("/api/v1")

	health := handler.HealthHandler{DB: d.DB, Version: d.Version}
	v1.Get("/health", health.Get)

	// Cualquier otra ruta responde 404 con el formato único.
	app.Use(func(c fiber.Ctx) error {
		return response.NotFound("La ruta " + c.Method() + " " + c.Path() + " no existe")
	})
}
