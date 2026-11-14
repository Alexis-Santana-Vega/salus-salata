// Package database abre la conexión a Firebird y aplica las migraciones.
package database

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net"
	"net/url"
	"time"

	_ "github.com/nakagami/firebirdsql" // registra el driver "firebirdsql"

	"salus/backend/internal/config"
)

// DSN arma la cadena de conexión por alias:
// firebirdsql://USUARIO:CLAVE@host:3050/salus?charset=UTF8&wire_crypt=true
func DSN(c config.DBConfig) string {
	u := url.URL{
		User:     url.UserPassword(c.User, c.Password),
		Host:     net.JoinHostPort(c.Host, c.Port),
		Path:     "/" + c.Alias,
		RawQuery: "charset=UTF8&wire_crypt=true&auth_plugin_name=Srp256",
	}
	return u.String()[2:] // quita el "//" inicial: el driver espera usuario:clave@host...
}

// Open abre el pool y espera a que Firebird responda (útil cuando el
// contenedor de la BD todavía está arrancando).
func Open(ctx context.Context, c config.DBConfig, log *slog.Logger) (*sql.DB, error) {
	db, err := sql.Open("firebirdsql", DSN(c))
	if err != nil {
		return nil, fmt.Errorf("abrir conexión: %w", err)
	}
	// Pool pequeño: Firebird crea un proceso/hilo por conexión.
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	// sql.Open no conecta; el Ping sí. Se reintenta ~30 s por si la BD arranca.
	const intentos = 15
	for i := 1; ; i++ {
		pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
		err = db.PingContext(pingCtx)
		cancel()
		if err == nil {
			log.Info("conectado a Firebird", "host", c.Host, "alias", c.Alias, "usuario", c.User)
			return db, nil
		}
		if i == intentos || ctx.Err() != nil {
			_ = db.Close()
			return nil, fmt.Errorf("no se pudo conectar a Firebird (%s:%s/%s): %w", c.Host, c.Port, c.Alias, err)
		}
		log.Warn("Firebird no responde, reintentando", "intento", i, "error", err)
		// Espera 2 s o sale antes si se cancela el arranque (Ctrl+C / SIGTERM).
		select {
		case <-time.After(2 * time.Second):
		case <-ctx.Done():
		}
	}
}
