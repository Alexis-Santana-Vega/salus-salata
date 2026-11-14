package security

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"salus/backend/internal/model"
)

const emisor = "salus-api"

// Claims del token de acceso. Lleva lo mínimo para autorizar sin ir a la BD.
type Claims struct {
	Rol string `json:"rol"`
	// DebeCambiar bloquea todo salvo cambiar la contraseña (primer acceso).
	DebeCambiar bool `json:"chg,omitempty"`
	jwt.RegisteredClaims
}

// UsuarioID devuelve el id del usuario guardado en "sub".
func (c Claims) UsuarioID() (model.ID, error) {
	return model.ParseID(c.Subject)
}

// Tokens firma y valida tokens de acceso (JWT HS256, vida corta).
type Tokens struct {
	secreto []byte
	vida    time.Duration
}

func NewTokens(secreto string, vida time.Duration) *Tokens {
	return &Tokens{secreto: []byte(secreto), vida: vida}
}

func (t *Tokens) Vida() time.Duration { return t.vida }

// Firmar emite un token de acceso para el usuario y rol dados.
func (t *Tokens) Firmar(usuarioID model.ID, rol string, debeCambiar bool, ahora time.Time) (string, error) {
	jti, err := RandomHex(16)
	if err != nil {
		return "", err
	}
	claims := Claims{
		Rol:         rol,
		DebeCambiar: debeCambiar,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    emisor,
			Subject:   usuarioID.String(),
			IssuedAt:  jwt.NewNumericDate(ahora),
			ExpiresAt: jwt.NewNumericDate(ahora.Add(t.vida)),
			ID:        jti,
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(t.secreto)
}

// Validar verifica firma, algoritmo (solo HS256), emisor y vigencia.
func (t *Tokens) Validar(raw string) (Claims, error) {
	var c Claims
	_, err := jwt.ParseWithClaims(raw, &c,
		func(*jwt.Token) (any, error) { return t.secreto, nil },
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer(emisor),
		jwt.WithExpirationRequired(),
		jwt.WithLeeway(5*time.Second),
	)
	if err != nil {
		return Claims{}, err
	}
	if c.Subject == "" || c.Rol == "" {
		return Claims{}, errors.New("token sin usuario o rol")
	}
	return c, nil
}

// RandomHex devuelve n bytes aleatorios criptográficos en hexadecimal.
func RandomHex(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// HashToken es el SHA-256 de un refresh token; es lo único que se guarda en BD.
// (Basta SHA-256: el token es aleatorio de 256 bits, no una contraseña.)
func HashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}
