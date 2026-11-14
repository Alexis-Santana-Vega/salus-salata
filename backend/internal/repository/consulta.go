package repository

import (
	"database/sql"
	"errors"
	"regexp"
	"strings"
)

// Consulta son los parámetros comunes de todos los listados.
type Consulta struct {
	Q      string // texto a buscar (vacío = sin filtro)
	Orden  string // campo de la API, con "-" adelante si es descendente
	Offset int
	Limite int
}

// ErrOrdenInvalido: el campo pedido en ?sort no se puede ordenar.
var ErrOrdenInvalido = errors.New("orden inválido")

// orderBy traduce el campo de la API a su columna con la lista blanca del
// recurso. Siempre desempata por el id (columna id, p. ej. "D.ID" con JOIN)
// para que la paginación sea estable.
func orderBy(orden string, columnas map[string]string, defecto, id string) (string, error) {
	if orden == "" {
		orden = defecto
	}
	dir := "ASC"
	if campo, ok := strings.CutPrefix(orden, "-"); ok {
		orden, dir = campo, "DESC"
	}
	col, ok := columnas[orden]
	if !ok {
		return "", ErrOrdenInvalido
	}
	return " ORDER BY " + col + " " + dir + ", " + id, nil
}

// ErrDuplicado: un valor único ya existe. Campo es el nombre en la API.
type ErrDuplicado struct{ Campo string }

func (e ErrDuplicado) Error() string { return "duplicado: " + e.Campo }

// Firebird: attempt to store duplicate value … in unique index "UQ_X"
// o violation of PRIMARY or UNIQUE KEY constraint "UQ_X" (el driver a veces
// manda el nombre sin comillas).
var reIndiceUnico = regexp.MustCompile(`(?:unique index|UNIQUE KEY constraint) "?([A-Za-z0-9_$]+)"?`)

// duplicado convierte la violación de un índice único en ErrDuplicado usando
// el mapa índice → campo del recurso. Otros errores pasan igual.
func duplicado(err error, indices map[string]string) error {
	if err == nil {
		return nil
	}
	if m := reIndiceUnico.FindStringSubmatch(err.Error()); m != nil {
		if campo, ok := indices[m[1]]; ok {
			return ErrDuplicado{Campo: campo}
		}
	}
	return err
}

// unaFila confirma que un UPDATE tocó un registro; si no, no existe.
func unaFila(res sql.Result) error {
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNoEncontrado
	}
	return nil
}

// nulo guarda un texto opcional vacío como NULL.
func nulo(s string) any {
	if s == "" {
		return nil
	}
	return s
}

// ErrReferencia: el registro apunta a algo que no existe (llave foránea).
type ErrReferencia struct{ Campo string }

func (e ErrReferencia) Error() string { return "referencia inválida: " + e.Campo }

var reLlaveForanea = regexp.MustCompile(`FOREIGN KEY constraint "?([A-Za-z0-9_$]+)"?`)

// integridad traduce violaciones de índices únicos (ErrDuplicado) y de llaves
// foráneas (ErrReferencia) con los mapas restricción → campo del recurso.
func integridad(err error, unicos, foraneas map[string]string) error {
	if err == nil {
		return nil
	}
	if m := reLlaveForanea.FindStringSubmatch(err.Error()); m != nil {
		if campo, ok := foraneas[m[1]]; ok {
			return ErrReferencia{Campo: campo}
		}
	}
	return duplicado(err, unicos)
}
