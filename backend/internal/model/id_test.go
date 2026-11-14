package model

import (
	"encoding/json"
	"testing"
)

func TestID(t *testing.T) {
	id := NuevoID()

	// Ida y vuelta por la BD: 16 bytes.
	v, err := id.Value()
	if err != nil || len(v.([]byte)) != 16 {
		t.Fatalf("Value: %v %v", v, err)
	}
	var leido ID
	if err := leido.Scan(v); err != nil || leido != id {
		t.Fatalf("Scan: %v %v", leido, err)
	}

	// Ida y vuelta por JSON: texto.
	b, _ := json.Marshal(struct{ ID ID }{id})
	var dst struct{ ID ID }
	if err := json.Unmarshal(b, &dst); err != nil || dst.ID != id {
		t.Fatalf("JSON: %s %v", b, err)
	}

	// Cero = NULL.
	if v, _ := (ID{}).Value(); v != nil {
		t.Fatal("el id cero debe guardarse como NULL")
	}
	if _, err := ParseID("no-es-uuid"); err == nil {
		t.Fatal("debió rechazar un id inválido")
	}
}
