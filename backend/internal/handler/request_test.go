package handler

import (
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestReglas(t *testing.T) {
	type req struct {
		Nombre   string `json:"nombre" validate:"required,max=120,nombre"`
		Puesto   string `json:"puesto" validate:"required,oneof=AYUDANTE OTRO"`
		Telefono string `json:"telefono" validate:"omitempty,telefono"`
	}
	casos := []struct {
		r  req
		ok bool
	}{
		{req{"Rodrigo", "AYUDANTE", ""}, true},
		{req{"Dra. María-José O'Neil Peña", "AYUDANTE", "+52 (477) 123-4567"}, true},
		{req{"Rodrigo", "AYUDANTE", "abc"}, false},
		{req{"Rodrigo", "AYUDANTE", "12345"}, false},
		{req{"Rodrigo", "JEFE", ""}, false},
		{req{"<b>x</b>", "AYUDANTE", ""}, false},
		{req{"Rodrigo2", "AYUDANTE", ""}, false},
		{req{"x'; DROP TABLE USUARIO;--", "AYUDANTE", ""}, false},
		{req{"Ana 😀", "AYUDANTE", ""}, false},
	}
	for _, c := range casos {
		if err := validate.Struct(c.r); (err == nil) != c.ok {
			t.Errorf("%+v: esperado ok=%v, error=%v", c.r, c.ok, err)
		}
	}
}

func TestFormatos(t *testing.T) {
	casos := map[string]map[string]bool{
		"textogeneral": {"Artroscopia / Morcelación": true, "Av. Juárez #120 (esq. 5 de Mayo)": true, "Holep; DROP": false, "a{b}": false, "x=1": false, `di"ce`: false, " empieza con espacio": false},
		"numdomicilio": {"120-B": true, "S/N": true, "Int. 4": true, "12;": false},
		"codigo":       {"DOC-001": true, "GTO-001": true, "DOC 001": false, "DÓC": false, "-01": false},
		"cp":           {"42000": true, "4200": false, "1.234": false},
	}
	for tag, valores := range casos {
		for v, ok := range valores {
			if reglas[tag].re.MatchString(v) != ok {
				t.Errorf("%s(%q): esperado %v", tag, v, ok)
			}
		}
	}
}

func TestCampoConIndice(t *testing.T) {
	type contacto struct {
		Nombre string `json:"nombre" validate:"required"`
	}
	type req struct {
		Contactos []contacto `json:"contactos" validate:"dive"`
	}
	err := validate.Struct(req{Contactos: []contacto{{"A"}, {""}}})
	if err == nil {
		t.Fatal("debió fallar")
	}
	for _, fe := range err.(validator.ValidationErrors) {
		if got := campo(fe); got != "contactos[1].nombre" {
			t.Fatalf("campo = %q", got)
		}
	}
}
