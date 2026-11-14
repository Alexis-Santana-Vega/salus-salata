package response

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"
)

// Códigos de error estables: el frontend decide con ellos, no con el texto.
const (
	CodeBadRequest         = "BAD_REQUEST"
	CodeValidation         = "VALIDATION_ERROR"
	CodeUnauthorized       = "UNAUTHORIZED"
	CodeInvalidCredentials = "INVALID_CREDENTIALS"
	CodePasswordChange     = "PASSWORD_CHANGE_REQUIRED"
	CodeForbidden          = "FORBIDDEN"
	CodeNotFound           = "NOT_FOUND"
	CodeMethodNotAllowed   = "METHOD_NOT_ALLOWED"
	CodeConflict           = "CONFLICT"
	CodePayloadTooLarge    = "PAYLOAD_TOO_LARGE"
	CodeUnsupportedMedia   = "UNSUPPORTED_MEDIA_TYPE"
	CodeTooManyRequests    = "TOO_MANY_REQUESTS"
	CodeInternal           = "INTERNAL_ERROR"
	CodeServiceUnavailable = "SERVICE_UNAVAILABLE"
)

// FieldError describe un problema en un campo de la petición.
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Error es el único tipo de error que viaja al cliente.
type Error struct {
	Status  int          `json:"-"`
	Code    string       `json:"code"`
	Message string       `json:"message"`
	Details []FieldError `json:"details,omitempty"`
}

func (e *Error) Error() string { return e.Code + ": " + e.Message }

func New(status int, code, msg string) *Error {
	return &Error{Status: status, Code: code, Message: msg}
}

func BadRequest(msg string) *Error { return New(http.StatusBadRequest, CodeBadRequest, msg) }
func Unauthorized(msg string) *Error {
	return New(http.StatusUnauthorized, CodeUnauthorized, msg)
}
func Forbidden(msg string) *Error { return New(http.StatusForbidden, CodeForbidden, msg) }
func NotFound(msg string) *Error  { return New(http.StatusNotFound, CodeNotFound, msg) }
func Conflict(msg string) *Error  { return New(http.StatusConflict, CodeConflict, msg) }
func TooManyRequests(msg string) *Error {
	return New(http.StatusTooManyRequests, CodeTooManyRequests, msg)
}
func UnsupportedMediaType() *Error {
	return New(http.StatusUnsupportedMediaType, CodeUnsupportedMedia, "El contenido debe enviarse como JSON")
}
func ServiceUnavailable(msg string) *Error {
	return New(http.StatusServiceUnavailable, CodeServiceUnavailable, msg)
}

// Validation responde 422 con el detalle por campo.
func Validation(details ...FieldError) *Error {
	e := New(http.StatusUnprocessableEntity, CodeValidation, "Revisa los datos enviados")
	e.Details = details
	return e
}

type errorBody struct {
	Error struct {
		*Error
		RequestID string `json:"requestId,omitempty"`
	} `json:"error"`
}

// ErrorHandler convierte cualquier error en el formato único. Los errores
// inesperados se registran completos en el log, pero al cliente solo le llega
// un mensaje genérico con el requestId para rastrearlo.
func ErrorHandler(log *slog.Logger) fiber.ErrorHandler {
	return func(c fiber.Ctx, err error) error {
		var appErr *Error
		var fe *fiber.Error

		// Traduce el error a *Error: propio, de Fiber o inesperado (500).
		switch {
		case errors.As(err, &appErr):
		case errors.As(err, &fe):
			appErr = fromFiber(fe)
		default:
			appErr = New(http.StatusInternalServerError, CodeInternal, "Ocurrió un error inesperado. Intenta de nuevo.")
		}

		// Solo los 5xx van al log como error: son fallas nuestras, no del cliente.
		rid := requestid.FromContext(c)
		if appErr.Status >= http.StatusInternalServerError {
			log.Error("error en la petición",
				"requestId", rid, "method", c.Method(), "path", c.Path(), "error", err)
		}

		var body errorBody
		body.Error.Error = appErr
		body.Error.RequestID = rid
		return c.Status(appErr.Status).JSON(body)
	}
}

// fromFiber traduce los errores internos de Fiber al formato único, sin
// exponer su mensaje técnico.
func fromFiber(fe *fiber.Error) *Error {
	switch fe.Code {
	case http.StatusNotFound:
		return NotFound("El recurso no existe")
	case http.StatusMethodNotAllowed:
		return New(fe.Code, CodeMethodNotAllowed, "Método no permitido")
	case http.StatusRequestEntityTooLarge:
		return New(fe.Code, CodePayloadTooLarge, "La petición es demasiado grande")
	case http.StatusTooManyRequests:
		return TooManyRequests("Demasiadas peticiones. Espera un momento.")
	case http.StatusBadRequest:
		return BadRequest("Petición inválida")
	}
	if fe.Code >= http.StatusInternalServerError {
		return New(http.StatusInternalServerError, CodeInternal, "Ocurrió un error inesperado. Intenta de nuevo.")
	}
	return New(fe.Code, CodeBadRequest, fe.Message)
}
