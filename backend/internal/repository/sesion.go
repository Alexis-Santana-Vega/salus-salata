package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"salus/backend/internal/model"
)

type Sesiones struct{ db DBTX }

func NewSesiones(db DBTX) *Sesiones { return &Sesiones{db: db} }

func (r *Sesiones) Crear(ctx context.Context, s model.Sesion, o model.Origen) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO SESION (ID, USUARIO_ID, FAMILIA, TOKEN_HASH, CREADA_EN, EXPIRA_EN, IP, USER_AGENT)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		model.NuevoID(), s.UsuarioID, s.Familia, s.TokenHash, s.CreadaEn, s.ExpiraEn, recortar(o.IP, 45), recortar(o.UserAgent, 255))
	return err
}

func (r *Sesiones) PorTokenHash(ctx context.Context, hash string) (*model.Sesion, error) {
	var s model.Sesion
	var revocada sql.NullTime
	err := r.db.QueryRowContext(ctx, `
		SELECT ID, USUARIO_ID, FAMILIA, TOKEN_HASH, CREADA_EN, EXPIRA_EN, REVOCADA_EN
		FROM SESION WHERE TOKEN_HASH = ?`, hash).
		Scan(&s.ID, &s.UsuarioID, &s.Familia, &s.TokenHash, &s.CreadaEn, &s.ExpiraEn, &revocada)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNoEncontrado
	}
	if err != nil {
		return nil, err
	}
	s.Familia = strings.TrimSpace(s.Familia)
	if revocada.Valid {
		s.RevocadaEn = &revocada.Time
	}
	return &s, nil
}

// Revocar marca la sesión como usada/cerrada. Devuelve false si ya estaba
// revocada: así dos peticiones simultáneas con el mismo token no pasan ambas.
func (r *Sesiones) Revocar(ctx context.Context, id model.ID, ahora time.Time) (bool, error) {
	res, err := r.db.ExecContext(ctx,
		`UPDATE SESION SET REVOCADA_EN = ? WHERE ID = ? AND REVOCADA_EN IS NULL`, ahora, id)
	if err != nil {
		return false, err
	}
	n, err := res.RowsAffected()
	return n == 1, err
}

// RevocarFamilia cierra todas las rotaciones de una sesión (reuso detectado).
func (r *Sesiones) RevocarFamilia(ctx context.Context, familia string, ahora time.Time) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE SESION SET REVOCADA_EN = ? WHERE FAMILIA = ? AND REVOCADA_EN IS NULL`, ahora, familia)
	return err
}

// RevocarDeUsuario cierra todas las sesiones abiertas del usuario.
func (r *Sesiones) RevocarDeUsuario(ctx context.Context, usuarioID model.ID, ahora time.Time) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE SESION SET REVOCADA_EN = ? WHERE USUARIO_ID = ? AND REVOCADA_EN IS NULL`, ahora, usuarioID)
	return err
}

// BorrarVencidas limpia las sesiones del usuario que ya no sirven (se llama
// al iniciar sesión; no hace falta un proceso en segundo plano).
func (r *Sesiones) BorrarVencidas(ctx context.Context, usuarioID model.ID, ahora time.Time) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM SESION WHERE USUARIO_ID = ? AND EXPIRA_EN < ?`, usuarioID, ahora)
	return err
}

func recortar(s string, max int) string {
	if len([]rune(s)) <= max {
		return s
	}
	return string([]rune(s)[:max])
}
