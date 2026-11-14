// Comando de desarrollo: imprime el hash argon2id de una contraseña leída de
// la entrada estándar (no por argumento, para que no quede en el historial).
// Lo usan los scripts de prueba para crear usuarios temporales.
//
//	echo -n "contraseña" | go run ./cmd/hash
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"salus/backend/internal/security"
)

func main() {
	linea, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	plano := strings.TrimRight(linea, "\r\n")
	if plano == "" {
		fmt.Fprintln(os.Stderr, "falta la contraseña en la entrada estándar")
		os.Exit(1)
	}
	h, err := security.HashPassword(plano)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Print(h)
}
