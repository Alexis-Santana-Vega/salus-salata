package config

import "testing"

func TestValidateAlias(t *testing.T) {
	casos := map[string]bool{
		"salus":                            true,
		"":                                 false,
		"/var/lib/firebird/data/salus.fdb": false,
		`C:\datos\salus.fdb`:               false,
		"salus.fdb":                        false,
	}
	for alias, valido := range casos {
		if err := validateAlias(alias); (err == nil) != valido {
			t.Errorf("alias %q: válido=%v, error=%v", alias, valido, err)
		}
	}
}

func TestValidateOrigin(t *testing.T) {
	casos := []struct {
		origen string
		prod   bool
		valido bool
	}{
		{"http://localhost:3000", false, true},
		{"https://salus.mx", true, true},
		{"http://salus.mx", true, false},
		{"*", false, false},
		{"https://*.salus.mx", true, false},
		{"http://localhost:3000/app", false, false},
		{"localhost:3000", false, false},
	}
	for _, c := range casos {
		if err := validateOrigin(c.origen, c.prod); (err == nil) != c.valido {
			t.Errorf("origen %q (prod=%v): válido=%v, error=%v", c.origen, c.prod, c.valido, err)
		}
	}
}

func TestLoadRechazaSYSDBA(t *testing.T) {
	t.Setenv("FIREBIRD_DB_ALIAS", "salus")
	t.Setenv("FIREBIRD_USER", "sysdba")
	t.Setenv("FIREBIRD_PASSWORD", "x")
	if _, err := Load(); err == nil {
		t.Fatal("esperaba error al usar SYSDBA")
	}
}
