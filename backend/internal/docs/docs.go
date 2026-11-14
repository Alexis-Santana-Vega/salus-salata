// Package docs sirve la documentación de la API con Scalar.
// Solo se monta cuando DOCS_ENABLED=true (nunca en producción).
package docs

import (
	_ "embed"

	"github.com/gofiber/fiber/v3"
)

//go:embed openapi.yaml
var spec []byte

// Versión fija de Scalar (no "latest") para que la página no cambie sola.
const scalarScript = "https://cdn.jsdelivr.net/npm/@scalar/api-reference@1.75.0"

var page = []byte(`<!doctype html>
<html lang="es">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>SALUS API · Documentación</title>
</head>
<body>
  <script id="api-reference" data-url="/docs/openapi.yaml"></script>
  <script src="` + scalarScript + `"></script>
</body>
</html>`)

// Register monta /docs (página) y /docs/openapi.yaml (especificación).
func Register(app *fiber.App) {
	app.Get("/docs", func(c fiber.Ctx) error {
		c.Set("Content-Security-Policy",
			"default-src 'self'; script-src 'self' https://cdn.jsdelivr.net; "+
				"style-src 'self' 'unsafe-inline' https://cdn.jsdelivr.net https://fonts.googleapis.com; "+
				"font-src 'self' data: https://fonts.scalar.com https://fonts.gstatic.com https://cdn.jsdelivr.net; "+
				"img-src 'self' data:; connect-src 'self'; frame-ancestors 'none'")
		c.Type("html", "utf-8")
		return c.Send(page)
	})
	app.Get("/docs/openapi.yaml", func(c fiber.Ctx) error {
		c.Type("yaml", "utf-8")
		return c.Send(spec)
	})
}
