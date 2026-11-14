package repository

import (
	"errors"
	"testing"
)

func TestOrderBy(t *testing.T) {
	cols := map[string]string{"nombre": "NOMBRE"}
	casos := []struct{ orden, want string }{
		{"", " ORDER BY NOMBRE ASC, ID"},
		{"nombre", " ORDER BY NOMBRE ASC, ID"},
		{"-nombre", " ORDER BY NOMBRE DESC, ID"},
	}
	for _, c := range casos {
		got, err := orderBy(c.orden, cols, "nombre", "ID")
		if err != nil || got != c.want {
			t.Errorf("orderBy(%q) = %q, %v", c.orden, got, err)
		}
	}
	// Nada fuera de la lista blanca llega al SQL.
	for _, malo := range []string{"ID; DROP TABLE X", "otro", "--nombre"} {
		if _, err := orderBy(malo, cols, "nombre", "ID"); !errors.Is(err, ErrOrdenInvalido) {
			t.Errorf("orderBy(%q) debió rechazarse", malo)
		}
	}
}

func TestDuplicado(t *testing.T) {
	idx := map[string]string{"UQ_PROCEDIMIENTO_NOMBRE": "nombre"}
	err := duplicado(errors.New(`attempt to store duplicate value (visible to active transactions) in unique index "UQ_PROCEDIMIENTO_NOMBRE"`), idx)
	var dup ErrDuplicado
	if !errors.As(err, &dup) || dup.Campo != "nombre" {
		t.Fatalf("no detectó el duplicado: %v", err)
	}
	// Así lo manda el driver de Firebird: sin comillas.
	err = duplicado(errors.New("attempt to store duplicate value (visible to active transactions) in unique index UQ_PROCEDIMIENTO_NOMBRE\nProblematic key"), idx)
	if !errors.As(err, &dup) || dup.Campo != "nombre" {
		t.Fatalf("no detectó el duplicado sin comillas: %v", err)
	}
	otro := errors.New("otra cosa")
	if duplicado(otro, idx) != otro {
		t.Fatal("un error ajeno debe pasar igual")
	}
}

func TestIntegridad(t *testing.T) {
	err := integridad(errors.New(`violation of FOREIGN KEY constraint "FK_DOCTOR_ASENTAMIENTO" on table "DOCTOR"`),
		nil, map[string]string{"FK_DOCTOR_ASENTAMIENTO": "asentamientoId"})
	var ref ErrReferencia
	if !errors.As(err, &ref) || ref.Campo != "asentamientoId" {
		t.Fatalf("no detectó la llave foránea: %v", err)
	}
}
