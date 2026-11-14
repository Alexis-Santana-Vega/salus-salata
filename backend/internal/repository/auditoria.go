package repository

import (
	"context"
	"database/sql"
	"time"

	"salus/backend/internal/model"
)

type Auditoria struct{ db DBTX }

func NewAuditoria(db DBTX) *Auditoria { return &Auditoria{db: db} }

// Registrar guarda un evento de acceso. usuarioID es 0 si no se conoce.
func (r *Auditoria) Registrar(ctx context.Context, evento string, usuarioID int64, ident string, o model.Origen, ahora time.Time) error {
	var uid sql.NullInt64
	if usuarioID > 0 {
		uid = sql.NullInt64{Int64: usuarioID, Valid: true}
	}
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO AUDITORIA_ACCESO (FECHA, EVENTO, USUARIO_ID, IDENTIFICADOR, IP, USER_AGENT)
		VALUES (?, ?, ?, ?, ?, ?)`,
		ahora, evento, uid, recortar(ident, 254), recortar(o.IP, 45), recortar(o.UserAgent, 255))
	return err
}

// FallidosDesde cuenta intentos fallidos para un identificador y para una IP
// desde un instante dado. Usa los índices (IDENTIFICADOR, FECHA) e (IP, FECHA).
func (r *Auditoria) FallidosDesde(ctx context.Context, ident, ip string, desde time.Time) (porIdent, porIP int, err error) {
	err = r.db.QueryRowContext(ctx, `
		SELECT
		  (SELECT COUNT(*) FROM AUDITORIA_ACCESO
		    WHERE IDENTIFICADOR = ? AND FECHA >= ? AND EVENTO = ?),
		  (SELECT COUNT(*) FROM AUDITORIA_ACCESO
		    WHERE IP = ? AND FECHA >= ? AND EVENTO = ?)
		FROM RDB$DATABASE`,
		recortar(ident, 254), desde, model.EventoLoginFallido,
		recortar(ip, 45), desde, model.EventoLoginFallido).Scan(&porIdent, &porIP)
	return porIdent, porIP, err
}
