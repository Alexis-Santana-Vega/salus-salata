package handler

import (
	"math"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Reglas de formato de toda la API. Son listas BLANCAS (lo que sí se permite):
// todo lo demás se rechaza, incluidos < > { } [ ] \ ` $ | ^ ~ = * % " ; y
// caracteres de control o emojis. El front aplica las mismas (utils/rules.ts).
// La defensa real contra inyección sigue siendo el SQL parametrizado y el
// escape de Vue; estas reglas evitan datos raros desde la captura.
var (
	// Nombres de personas: letras (con acentos y ñ), espacios, punto, apóstrofo y guion.
	reNombre = regexp.MustCompile(`^[\p{L}\p{M}][\p{L}\p{M} .'-]*$`)
	// Texto general (procedimientos, calles, puestos, nombres de hospital):
	// letras, números, espacios y . , # / ( ) ' & + - °
	reTextoGeneral = regexp.MustCompile(`^[\p{L}\p{M}\p{N}][\p{L}\p{M}\p{N} .,#/()'&+°-]*$`)
	// Número exterior o interior: letras, números, espacios y # / . - (p. ej. "120-B", "S/N").
	reNumDomicilio = regexp.MustCompile(`^[\p{L}\p{N}][\p{L}\p{N} #/.-]*$`)
	// Códigos y claves: letras sin acento, números y guion (DOC-001, GTO-001).
	reCodigo = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9-]*$`)
	// Teléfono: + opcional y de 7 a 20 caracteres entre dígitos, espacios, guiones y paréntesis.
	reTelefono = regexp.MustCompile(`^\+?[0-9][0-9 ()-]{6,19}$`)
	// Código postal: exactamente 5 dígitos.
	reCP = regexp.MustCompile(`^[0-9]{5}$`)
	// Número de serie: letras sin acento, números y . / - (p. ej. "LX300-88412").
	reSerie = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9./-]*$`)
)

// reglas: etiqueta → expresión y mensaje para el usuario.
var reglas = map[string]struct {
	re      *regexp.Regexp
	mensaje string
}{
	"nombre":       {reNombre, "Usa solo letras, espacios, punto, apóstrofo o guion"},
	"textogeneral": {reTextoGeneral, "Usa solo letras, números, espacios y . , # / ( ) ' & + - °"},
	"numdomicilio": {reNumDomicilio, "Usa solo letras, números, espacios y # / . -"},
	"codigo":       {reCodigo, "Usa solo letras sin acento, números y guion"},
	"telefono":     {reTelefono, "Escribe un teléfono válido (solo números, espacios, guiones o +)"},
	"cp":           {reCP, "El código postal tiene 5 dígitos"},
	"serie":        {reSerie, "Usa solo letras sin acento, números y . / -"},
}

func nuevoValidador() *validator.Validate {
	v := validator.New(validator.WithRequiredStructEnabled())
	// Los errores usan el nombre del campo JSON, no el de Go.
	v.RegisterTagNameFunc(func(f reflect.StructField) string {
		nombre := strings.SplitN(f.Tag.Get("json"), ",", 2)[0]
		if nombre == "-" {
			return ""
		}
		return nombre
	})
	// "dosdecimales": importes con a lo más 2 decimales (centavos).
	_ = v.RegisterValidation("dosdecimales", func(fl validator.FieldLevel) bool {
		c := fl.Field().Float() * 100
		return math.Abs(c-math.Round(c)) < 1e-6
	})
	for tag, r := range reglas {
		// Se valida sin los espacios de los extremos: el servicio los quita al guardar.
		_ = v.RegisterValidation(tag, func(fl validator.FieldLevel) bool {
			return r.re.MatchString(strings.TrimSpace(fl.Field().String()))
		})
	}
	return v
}

// campo es el nombre del campo para el front: el nombre JSON del campo y,
// dentro de una lista, su ruta con índice ("contactos[0].nombre"). Los
// campos de objetos anidados sin lista conservan solo su nombre ("calle").
func campo(fe validator.FieldError) string {
	ns := fe.Namespace()
	if i := strings.Index(ns, "."); i >= 0 {
		ns = ns[i+1:] // quita el nombre del struct raíz
	}
	if j := strings.LastIndex(ns, "]"); j >= 0 {
		inicio := strings.LastIndex(ns[:strings.Index(ns, "[")], ".") + 1
		return ns[inicio:]
	}
	return fe.Field()
}
