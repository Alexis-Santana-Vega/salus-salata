package handler

import "testing"

func TestReglas(t *testing.T) {
	type req struct {
		Nombre   string `json:"nombre" validate:"required,max=120,texto"`
		Puesto   string `json:"puesto" validate:"required,oneof=AYUDANTE OTRO"`
		Telefono string `json:"telefono" validate:"omitempty,telefono"`
	}
	casos := []struct {
		r  req
		ok bool
	}{
		{req{"Rodrigo", "AYUDANTE", ""}, true},
		{req{"Rodrigo", "AYUDANTE", "+52 (477) 123-4567"}, true},
		{req{"Rodrigo", "AYUDANTE", "477 123 4567"}, true},
		{req{"Rodrigo", "AYUDANTE", "abc"}, false},
		{req{"Rodrigo", "AYUDANTE", "12345"}, false},
		{req{"Rodrigo", "JEFE", ""}, false},
		{req{"<b>x</b>", "AYUDANTE", ""}, false},
	}
	for _, c := range casos {
		if err := validate.Struct(c.r); (err == nil) != c.ok {
			t.Errorf("%+v: esperado ok=%v, error=%v", c.r, c.ok, err)
		}
	}
}

func TestDireccion(t *testing.T) {
	casos := []struct {
		d  direccionReq
		ok bool
	}{
		{direccionReq{}, true},
		{direccionReq{CP: "37000", Estado: "GTO"}, true},
		{direccionReq{CP: "1.234"}, false},
		{direccionReq{CP: "3700"}, false},
		{direccionReq{Estado: "XX"}, false},
	}
	for _, c := range casos {
		if err := validate.Struct(c.d); (err == nil) != c.ok {
			t.Errorf("%+v: esperado ok=%v, error=%v", c.d, c.ok, err)
		}
	}
}
