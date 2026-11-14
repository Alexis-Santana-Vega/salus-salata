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

func TestLoadRechazaSYSDBA(t *testing.T) {
	t.Setenv("FIREBIRD_DB_ALIAS", "salus")
	t.Setenv("FIREBIRD_USER", "sysdba")
	t.Setenv("FIREBIRD_PASSWORD", "x")
	if _, err := Load(); err == nil {
		t.Fatal("esperaba error al usar SYSDBA")
	}
}
