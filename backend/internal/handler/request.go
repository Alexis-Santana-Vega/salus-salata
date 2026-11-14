package handler

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"

	"salus/backend/internal/response"
)

// validate es seguro para uso concurrente y cachea las reglas de cada struct.
var validate = nuevoValidador()

// bind lee el cuerpo JSON en dst y lo valida con las etiquetas `validate`.
// Los campos desconocidos se ignoran: así el cliente puede mandar campos
// nuevos sin romper versiones viejas del servidor.
func bind(c fiber.Ctx, dst any) error {
	if err := json.Unmarshal(c.Body(), dst); err != nil {
		return response.BadRequest("El cuerpo de la petición no es un JSON válido")
	}
	if err := validate.Struct(dst); err != nil {
		var ves validator.ValidationErrors
		if !errors.As(err, &ves) {
			return err
		}
		detalles := make([]response.FieldError, 0, len(ves))
		for _, fe := range ves {
			detalles = append(detalles, response.FieldError{Field: campo(fe), Message: mensaje(fe)})
		}
		return response.Validation(detalles...)
	}
	return nil
}

// mensaje traduce cada regla a un texto para el usuario.
func mensaje(fe validator.FieldError) string {
	if r, ok := reglas[fe.Tag()]; ok {
		return r.mensaje
	}
	switch fe.Tag() {
	case "required":
		return "Este campo es obligatorio"
	case "required_if":
		return "Este campo es obligatorio"
	case "email":
		return "Escribe un correo válido"
	case "min":
		return "Usa al menos " + fe.Param() + " caracteres"
	case "max":
		return "Usa como máximo " + fe.Param() + " caracteres"
	case "oneof":
		return "Elige una opción de la lista"
	case "gte":
		return "Debe ser mayor o igual a " + fe.Param()
	case "lte":
		return "Debe ser menor o igual a " + fe.Param()
	case "startswith":
		return "Debe empezar con " + fe.Param()
	case "dosdecimales":
		return "Usa a lo más 2 decimales"
	case "url":
		return "Escribe una liga completa (https://…)"
	case "latitude":
		return "Latitud entre -90 y 90"
	case "longitude":
		return "Longitud entre -180 y 180"
	case "dive", "min_items", "maxitems":
		return "Revisa la lista"
	}
	return "Valor no válido"
}

// Paginación estándar de todos los listados: ?page=1&pageSize=20 (máx. 100).
const (
	pageSizeDefault = 20
	pageSizeMax     = 100
)

// paginacion lee page y pageSize; devuelve offset y límite para la consulta.
func paginacion(c fiber.Ctx) (page, size, offset int, err error) {
	page, size = 1, pageSizeDefault
	var detalles []response.FieldError
	if v := c.Query("page"); v != "" {
		if page, err = strconv.Atoi(v); err != nil || page < 1 {
			detalles = append(detalles, response.FieldError{Field: "page", Message: "Debe ser un número mayor o igual a 1"})
		}
	}
	if v := c.Query("pageSize"); v != "" {
		if size, err = strconv.Atoi(v); err != nil || size < 1 || size > pageSizeMax {
			detalles = append(detalles, response.FieldError{Field: "pageSize", Message: "Debe estar entre 1 y " + strconv.Itoa(pageSizeMax)})
		}
	}
	if len(detalles) > 0 {
		return 0, 0, 0, response.Validation(detalles...)
	}
	return page, size, (page - 1) * size, nil
}
