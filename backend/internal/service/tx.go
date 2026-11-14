package service

import (
	"context"
	"database/sql"
	"fmt"
)

// enTx ejecuta fn en una transacción: si fn falla (o hay panic) se revierte
// todo; si no, se confirma. Toda operación que escribe en varios pasos, o que
// lee para luego escribir, va aquí: nada de datos a medias ni huérfanos.
func enTx(ctx context.Context, db *sql.DB, fn func(tx *sql.Tx) error) (err error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("iniciar transacción: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
		if err != nil {
			_ = tx.Rollback()
		}
	}()
	if err = fn(tx); err != nil {
		return err
	}
	return tx.Commit()
}
