// Package handler contiene los manejadores HTTP: leen la petición, llaman al
// servicio y responden con el formato único.
package handler

import (
	"context"
	"database/sql"
	"time"

	"github.com/gofiber/fiber/v3"

	"salus/backend/internal/response"
)

type HealthHandler struct {
	DB      *sql.DB
	Version string
}

type healthData struct {
	Status   string `json:"status"`
	Database string `json:"database"`
	Version  string `json:"version"`
	Time     string `json:"time"`
}

// Get responde 200 si la API y la base están bien, 503 si la base no responde.
func (h HealthHandler) Get(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
	defer cancel()

	if err := h.DB.PingContext(ctx); err != nil {
		return response.ServiceUnavailable("La base de datos no está disponible")
	}
	return response.OK(c, healthData{
		Status:   "ok",
		Database: "ok",
		Version:  h.Version,
		Time:     time.Now().Format(time.RFC3339),
	})
}
