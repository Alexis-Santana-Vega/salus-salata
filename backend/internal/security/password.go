// Package security agrupa el hash de contraseñas y los tokens de sesión.
package security

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// Parámetros argon2id recomendados por OWASP (19 MiB, 2 pasadas, 1 hilo).
// Si se cambian, los hashes viejos siguen validando y NeedsRehash avisa.
const (
	argonMemoria = 19 * 1024 // KiB
	argonTiempo  = 2
	argonHilos   = 1
	argonSal     = 16 // bytes
	argonLargo   = 32 // bytes
)

var errHashInvalido = errors.New("hash de contraseña con formato inválido")

// HashPassword genera un hash argon2id con sal aleatoria, en formato PHC:
// $argon2id$v=19$m=19456,t=2,p=1$<sal>$<hash>
func HashPassword(plain string) (string, error) {
	sal := make([]byte, argonSal)
	if _, err := rand.Read(sal); err != nil {
		return "", err
	}
	h := argon2.IDKey([]byte(plain), sal, argonTiempo, argonMemoria, argonHilos, argonLargo)
	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, argonMemoria, argonTiempo, argonHilos,
		b64(sal), b64(h)), nil
}

// VerifyPassword compara en tiempo constante usando los parámetros guardados
// en el propio hash.
func VerifyPassword(plain, encoded string) (bool, error) {
	p, sal, esperado, err := decodificar(encoded)
	if err != nil {
		return false, err
	}
	h := argon2.IDKey([]byte(plain), sal, p.tiempo, p.memoria, p.hilos, uint32(len(esperado)))
	return subtle.ConstantTimeCompare(h, esperado) == 1, nil
}

// NeedsRehash indica si el hash se hizo con parámetros distintos a los actuales.
func NeedsRehash(encoded string) bool {
	p, _, _, err := decodificar(encoded)
	return err != nil || p.memoria != argonMemoria || p.tiempo != argonTiempo || p.hilos != argonHilos
}

type parametros struct {
	memoria, tiempo uint32
	hilos           uint8
}

func decodificar(encoded string) (parametros, []byte, []byte, error) {
	var p parametros
	partes := strings.Split(encoded, "$") // "", argon2id, v=19, m=..,t=..,p=.., sal, hash
	if len(partes) != 6 || partes[1] != "argon2id" {
		return p, nil, nil, errHashInvalido
	}
	var version int
	if _, err := fmt.Sscanf(partes[2], "v=%d", &version); err != nil || version != argon2.Version {
		return p, nil, nil, errHashInvalido
	}
	if _, err := fmt.Sscanf(partes[3], "m=%d,t=%d,p=%d", &p.memoria, &p.tiempo, &p.hilos); err != nil {
		return p, nil, nil, errHashInvalido
	}
	sal, err1 := base64.RawStdEncoding.DecodeString(partes[4])
	hash, err2 := base64.RawStdEncoding.DecodeString(partes[5])
	if err1 != nil || err2 != nil || len(hash) == 0 {
		return p, nil, nil, errHashInvalido
	}
	return p, sal, hash, nil
}

func b64(b []byte) string { return base64.RawStdEncoding.EncodeToString(b) }
