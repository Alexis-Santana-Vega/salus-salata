package handler

import (
	"errors"
	"strings"
	"unicode/utf8"

	"github.com/gofiber/fiber/v3"

	"salus/backend/internal/model"
	"salus/backend/internal/repository"
	"salus/backend/internal/response"
)

// Piezas comunes de los catálogos (CRUD con baja lógica). Cada catálogo
// nuevo solo declara su petición, su handler y sus rutas.

const maxBusqueda = 100

// consulta lee ?page, ?pageSize, ?q y ?sort con las mismas reglas en todos
// los listados. El campo de ?sort lo valida el repositorio (lista blanca).
func consulta(c fiber.Ctx) (repository.Consulta, response.Meta, error) {
	page, size, offset, err := paginacion(c)
	if err != nil {
		return repository.Consulta{}, response.Meta{}, err
	}
	q := strings.Join(strings.Fields(c.Query("q")), " ")
	if utf8.RuneCountInString(q) > maxBusqueda {
		return repository.Consulta{}, response.Meta{}, response.Validation(
			response.FieldError{Field: "q", Message: "Usa como máximo 100 caracteres"})
	}
	return repository.Consulta{Q: q, Orden: c.Query("sort"), Offset: offset, Limite: size},
		response.Meta{Page: page, PageSize: size}, nil
}

// idParam lee el {id} de la ruta. Un id mal formado es un recurso que no existe.
func idParam(c fiber.Ctx) (model.ID, error) {
	id, err := model.ParseID(c.Params("id"))
	if err != nil {
		return model.ID{}, response.NotFound("El registro no existe o ya fue dado de baja")
	}
	return id, nil
}

// Cómo se nombra cada campo único en el mensaje de duplicado.
var nombreCampoUnico = map[string]string{
	"nombre":       "este nombre",
	"clave":        "esta clave",
	"numero":       "este número",
	"codigoBarras": "este código de barras",
	"numSerie":     "este número de serie",
}

// errorCatalogo traduce los errores del repositorio al formato de la API.
func errorCatalogo(err error) error {
	var dup repository.ErrDuplicado
	switch {
	case errors.Is(err, repository.ErrNoEncontrado):
		return response.NotFound("El registro no existe o ya fue dado de baja")
	case errors.Is(err, repository.ErrOrdenInvalido):
		return response.Validation(response.FieldError{Field: "sort", Message: "No se puede ordenar por ese campo"})
	case errors.As(err, &dup):
		texto, ok := nombreCampoUnico[dup.Campo]
		if !ok {
			texto = "ese valor"
		}
		e := response.Conflict("Ya existe un registro con esos datos")
		e.Details = []response.FieldError{{Field: dup.Campo, Message: "Ya existe un registro con " + texto}}
		return e
	}
	return err
}
