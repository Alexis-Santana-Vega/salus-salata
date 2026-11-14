package repository

import (
	"context"
	"database/sql"
	"errors"
)

// Siguiente devuelve el siguiente número del consecutivo dado. Debe llamarse
// dentro de la transacción del alta: el UPDATE bloquea la fila hasta que la
// transacción termina, así dos altas simultáneas no toman el mismo número y,
// si el alta se revierte, el número tampoco se consume.
func Siguiente(ctx context.Context, db DBTX, nombre string) (int, error) {
	var valor int
	err := db.QueryRowContext(ctx,
		`UPDATE CONSECUTIVO SET VALOR = VALOR + 1 WHERE NOMBRE = ? RETURNING VALOR`, nombre).Scan(&valor)
	if errors.Is(err, sql.ErrNoRows) {
		// Primera vez: se crea en 1. Si otra transacción lo crea al mismo
		// tiempo, esta falla por llave duplicada y no se guarda nada.
		_, err = db.ExecContext(ctx, `INSERT INTO CONSECUTIVO (NOMBRE, VALOR) VALUES (?, 1)`, nombre)
		return 1, err
	}
	return valor, err
}
