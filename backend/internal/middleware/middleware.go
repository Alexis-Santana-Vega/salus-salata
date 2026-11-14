// Package middleware agrupa lo que corre antes de cada endpoint.
package middleware

import (
	"log/slog"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/helmet"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"

	"salus/backend/internal/config"
	"salus/backend/internal/response"
)

// Register aplica, en orden, los middlewares globales.
func Register(app *fiber.App, cfg config.Config, log *slog.Logger) {
	app.Use(requestid.New()) // X-Request-ID para rastrear cada petición en el log

	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c fiber.Ctx, e any) {
			log.Error("panic recuperado", "requestId", requestid.FromContext(c), "panic", e)
		},
	}))

	app.Use(AccessLog(log))

	// Cabeceras de seguridad. La API solo devuelve JSON, así que la CSP es
	// estricta; /docs define la suya propia.
	app.Use(helmet.New(helmet.Config{
		Next:                  func(c fiber.Ctx) bool { return strings.HasPrefix(c.Path(), "/docs") },
		ContentSecurityPolicy: "default-src 'none'; frame-ancestors 'none'",
		HSTSMaxAge:            hstsMaxAge(cfg),
	}))

	if len(cfg.CORSOrigins) > 0 {
		app.Use(cors.New(cors.Config{
			AllowOrigins:     cfg.CORSOrigins,
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Authorization", "Content-Type", "X-Request-ID"},
			AllowCredentials: true,
			MaxAge:           600,
		}))
	}

	// Límite general por IP contra abusos. El login tendrá uno propio, más
	// estricto, en el paso 0.4.
	app.Use(limiter.New(limiter.Config{
		Max:        300,
		Expiration: time.Minute,
		LimitReached: func(fiber.Ctx) error {
			return response.TooManyRequests("Demasiadas peticiones. Espera un momento.")
		},
	}))
}

func hstsMaxAge(cfg config.Config) int {
	if cfg.IsProduction() {
		return 31536000 // 1 año; solo tiene efecto detrás de HTTPS
	}
	return 0
}

// AccessLog registra una línea por petición. Nunca registra cuerpos,
// cabeceras de autorización ni cookies.
func AccessLog(log *slog.Logger) fiber.Handler {
	return func(c fiber.Ctx) error {
		inicio := time.Now()
		err := c.Next()

		status := c.Response().StatusCode()
		if err != nil {
			// El ErrorHandler todavía no escribió la respuesta: deducimos el código.
			if e, ok := err.(*response.Error); ok {
				status = e.Status
			} else if fe, ok := err.(*fiber.Error); ok {
				status = fe.Code
			} else {
				status = fiber.StatusInternalServerError
			}
		}

		nivel := slog.LevelInfo
		switch {
		case status >= 500:
			nivel = slog.LevelError
		case status >= 400:
			nivel = slog.LevelWarn
		}
		log.Log(c.Context(), nivel, "http",
			"requestId", requestid.FromContext(c),
			"method", c.Method(),
			"path", c.Path(),
			"status", status,
			"ms", time.Since(inicio).Milliseconds(),
			"ip", c.IP(),
		)
		return err
	}
}
