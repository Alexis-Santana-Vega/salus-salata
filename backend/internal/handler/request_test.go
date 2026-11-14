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

func TestCP(t *testing.T) {
	for cp, ok := range map[string]bool{"42000": true, "4200": false, "1.234": false, "42a00": false, "420000": false} {
		if reCP.MatchString(cp) != ok {
			t.Errorf("cp %q: esperado %v", cp, ok)
		}
	}
}
