package security

import (
	"strings"
	"testing"
	"time"
)

func TestPassword(t *testing.T) {
	h, err := HashPassword("Contraseña segura 123")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(h, "$argon2id$v=19$m=19456,t=2,p=1$") {
		t.Fatalf("formato inesperado: %s", h)
	}
	if ok, _ := VerifyPassword("Contraseña segura 123", h); !ok {
		t.Fatal("la contraseña correcta no validó")
	}
	if ok, _ := VerifyPassword("otra", h); ok {
		t.Fatal("una contraseña incorrecta validó")
	}
	h2, _ := HashPassword("Contraseña segura 123")
	if h == h2 {
		t.Fatal("dos hashes iguales: la sal no es aleatoria")
	}
	if NeedsRehash(h) {
		t.Fatal("no debería pedir rehash con parámetros actuales")
	}
}

func TestTokens(t *testing.T) {
	tk := NewTokens(strings.Repeat("k", 32), time.Minute)
	raw, err := tk.Firmar(7, "ADMIN", false, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	c, err := tk.Validar(raw)
	if err != nil {
		t.Fatal(err)
	}
	if id, _ := c.UsuarioID(); id != 7 || c.Rol != "ADMIN" {
		t.Fatalf("claims incorrectos: %+v", c)
	}

	// Otro secreto no valida; un token vencido tampoco.
	if _, err := NewTokens(strings.Repeat("x", 32), time.Minute).Validar(raw); err == nil {
		t.Fatal("validó con otro secreto")
	}
	viejo, _ := tk.Firmar(7, "ADMIN", false, time.Now().Add(-time.Hour))
	if _, err := tk.Validar(viejo); err == nil {
		t.Fatal("validó un token vencido")
	}
	// Algoritmo "none" debe rechazarse.
	if _, err := tk.Validar("eyJhbGciOiJub25lIiwidHlwIjoiSldUIn0.eyJzdWIiOiIxIiwicm9sIjoiU1lTQURNSU4ifQ."); err == nil {
		t.Fatal("aceptó alg none")
	}
}
