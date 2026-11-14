// Package middleware agrupa lo que corre antes de cada endpoint.
package middleware

import (
	"errors"
	"log/slog"
	"net"
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

// Límites por IP (ventana deslizante de 1 minuto). El login tendrá uno propio,
// más estricto, en el paso 0.4.
const (
	limiteGeneral   = 120 // cualquier petición
	limiteEscritura = 30  // POST, PUT, PATCH, DELETE
)

// Register aplica los middlewares globales. El orden importa: primero se
// identifica y protege la petición, al final se limita.
func Register(app *fiber.App, cfg config.Config, log *slog.Logger) {
	// Id único por petición (cabecera X-Request-ID) para rastrearla en el log.
	app.Use(requestid.New())

	// Un panic no tumba el servidor: se registra y se responde 500 genérico.
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c fiber.Ctx, e any) {
			log.Error("panic recuperado", "requestId", requestid.FromContext(c), "panic", e)
		},
	}))

	app.Use(AccessLog(log))

	// Cabeceras de seguridad para el navegador. La API solo devuelve JSON, así
	// que la CSP bloquea todo; /docs trae su propia CSP.
	app.Use(helmet.New(helmet.Config{
		Next:                  func(c fiber.Ctx) bool { return strings.HasPrefix(c.Path(), "/docs") },
		ContentSecurityPolicy: "default-src 'none'; frame-ancestors 'none'",
		HSTSMaxAge:            hstsMaxAge(cfg),
	}))

	// CORS solo si hay orígenes configurados (desarrollo). En producción el
	// front va por el mismo origen y esto no se monta.
	if len(cfg.CORSOrigins) > 0 {
		app.Use(cors.New(cors.Config{
			AllowOrigins:     cfg.CORSOrigins,
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
			AllowHeaders:     []string{"Authorization", "Content-Type", "X-Request-ID"},
			ExposeHeaders:    []string{"X-Request-ID", "Retry-After", "X-RateLimit-Limit", "X-RateLimit-Remaining"},
			AllowCredentials: true, // necesario para la cookie del refresh (paso 0.4)
			MaxAge:           600,
		}))
	}

	app.Use(RequireUserAgent())
	app.Use(RequireJSON())

	// Límite general y límite extra para escrituras, por IP.
	app.Use(newLimiter(limiteGeneral, "g:", func(c fiber.Ctx) bool { return esHealthcheckLocal(c) }))
	app.Use(newLimiter(limiteEscritura, "w:", func(c fiber.Ctx) bool { return !esEscritura(c.Method()) }))
}

// newLimiter crea un límite por IP con ventana deslizante (más difícil de
// burlar que una ventana fija, que se reinicia de golpe cada minuto).
func newLimiter(max int, prefijo string, saltar func(fiber.Ctx) bool) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:               max,
		Expiration:        time.Minute,
		LimiterMiddleware: limiter.SlidingWindow{},
		Next:              saltar,
		KeyGenerator:      func(c fiber.Ctx) string { return prefijo + c.IP() },
		LimitReached: func(fiber.Ctx) error {
			return response.TooManyRequests("Demasiadas peticiones. Espera un momento.")
		},
	})
}

// RequireUserAgent rechaza peticiones a la API sin User-Agent: los navegadores
// siempre lo envían; su ausencia delata scripts muy básicos.
func RequireUserAgent() fiber.Handler {
	return func(c fiber.Ctx) error {
		if strings.HasPrefix(c.Path(), "/api/") && strings.TrimSpace(c.Get(fiber.HeaderUserAgent)) == "" {
			return response.Forbidden("Petición no permitida")
		}
		return c.Next()
	}
}

// RequireJSON exige Content-Type JSON en escrituras con cuerpo. Evita envíos
// de formularios desde otros sitios y entradas en formatos inesperados.
func RequireJSON() fiber.Handler {
	return func(c fiber.Ctx) error {
		if esEscritura(c.Method()) && len(c.Body()) > 0 &&
			!strings.HasPrefix(c.Get(fiber.HeaderContentType), fiber.MIMEApplicationJSON) {
			return response.UnsupportedMediaType()
		}
		return c.Next()
	}
}

func esEscritura(m string) bool {
	return m == fiber.MethodPost || m == fiber.MethodPut || m == fiber.MethodPatch || m == fiber.MethodDelete
}

// esHealthcheckLocal deja pasar sin límite solo el chequeo de Docker
// (desde la misma máquina), no el de clientes externos.
func esHealthcheckLocal(c fiber.Ctx) bool {
	ip := net.ParseIP(c.IP())
	return c.Path() == "/api/v1/health" && ip != nil && ip.IsLoopback()
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

		// Si hubo error, el ErrorHandler aún no escribió la respuesta:
		// el código se deduce del propio error.
		status := c.Response().StatusCode()
		if err != nil {
			var appErr *response.Error
			var fe *fiber.Error
			switch {
			case errors.As(err, &appErr):
				status = appErr.Status
			case errors.As(err, &fe):
				status = fe.Code
			default:
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
