// Package repository es la única capa que escribe SQL. Siempre parametrizado.
package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"salus/backend/internal/model"
)

// ErrNoEncontrado se devuelve cuando la fila buscada no existe.
var ErrNoEncontrado = errors.New("no encontrado")

// DBTX lo cumplen *sql.DB y *sql.Tx: el mismo repositorio sirve dentro o
// fuera de una transacción.
type DBTX interface {
	ExecContext(ctx context.Context, q string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, q string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, q string, args ...any) *sql.Row
}

type Usuarios struct{ db DBTX }

func NewUsuarios(db DBTX) *Usuarios { return &Usuarios{db: db} }

// El rol se trae con JOIN en la misma consulta (evita N+1).
const selectUsuario = `
SELECT U.ID, U.USUARIO, U.CORREO, U.NOMBRE, R.CLAVE, U.ACTIVO,
       U.DEBE_CAMBIAR_CONTRASENA, U.ULTIMO_ACCESO, U.CREADO_EN, U.HASH_CONTRASENA
FROM USUARIO U
JOIN ROL R ON R.ID = U.ROL_ID`

func scanUsuario(s interface{ Scan(...any) error }) (*model.Usuario, error) {
	var u model.Usuario
	var ultimo sql.NullTime
	err := s.Scan(&u.ID, &u.Usuario, &u.Correo, &u.Nombre, &u.Rol, &u.Activo,
		&u.DebeCambiarContrasena, &ultimo, &u.CreadoEn, &u.HashContrasena)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNoEncontrado
	}
	if err != nil {
		return nil, err
	}
	u.Rol = strings.TrimSpace(u.Rol)
	if ultimo.Valid {
		u.UltimoAcceso = &ultimo.Time
	}
	return &u, nil
}

// PorIdentificador busca por usuario o correo (ya normalizados en minúsculas).
func (r *Usuarios) PorIdentificador(ctx context.Context, ident string) (*model.Usuario, error) {
	return scanUsuario(r.db.QueryRowContext(ctx,
		selectUsuario+` WHERE U.USUARIO = ? OR U.CORREO = ?`, ident, ident))
}

func (r *Usuarios) PorID(ctx context.Context, id int64) (*model.Usuario, error) {
	return scanUsuario(r.db.QueryRowContext(ctx, selectUsuario+` WHERE U.ID = ?`, id))
}

// Listar pagina por OFFSET: la tabla de usuarios es pequeña. Para tablas que
// crezcan se usará paginación por cursor (ver docs/09).
func (r *Usuarios) Listar(ctx context.Context, offset, limite int) ([]model.Usuario, int, error) {
	var total int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM USUARIO`).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.db.QueryContext(ctx,
		selectUsuario+` ORDER BY U.NOMBRE, U.ID OFFSET ? ROWS FETCH NEXT ? ROWS ONLY`, offset, limite)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	out := make([]model.Usuario, 0, limite)
	for rows.Next() {
		u, err := scanUsuario(rows)
		if err != nil {
			return nil, 0, err
		}
		out = append(out, *u)
	}
	return out, total, rows.Err()
}

// HaySysadminActivo indica si ya existe al menos un Sysadmin activo.
func (r *Usuarios) HaySysadminActivo(ctx context.Context) (bool, error) {
	var n int
	err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM USUARIO U JOIN ROL R ON R.ID = U.ROL_ID
		WHERE R.CLAVE = ? AND U.ACTIVO`, model.RolSysadmin).Scan(&n)
	return n > 0, err
}

// Crear inserta el usuario y devuelve su id (RETURNING de Firebird).
func (r *Usuarios) Crear(ctx context.Context, u model.Usuario, ahora time.Time) (int64, error) {
	var id int64
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO USUARIO (USUARIO, CORREO, NOMBRE, HASH_CONTRASENA, ROL_ID, ACTIVO,
		                     DEBE_CAMBIAR_CONTRASENA, CREADO_EN, ACTUALIZADO_EN)
		VALUES (?, ?, ?, ?, (SELECT ID FROM ROL WHERE CLAVE = ?), ?, ?, ?, ?)
		RETURNING ID`,
		u.Usuario, u.Correo, u.Nombre, u.HashContrasena, u.Rol, u.Activo,
		u.DebeCambiarContrasena, ahora, ahora).Scan(&id)
	return id, err
}

func (r *Usuarios) RegistrarAcceso(ctx context.Context, id int64, ahora time.Time) error {
	_, err := r.db.ExecContext(ctx, `UPDATE USUARIO SET ULTIMO_ACCESO = ? WHERE ID = ?`, ahora, id)
	return err
}

// CambiarHash guarda un hash nuevo (cambio de contraseña o rehash por
// parámetros nuevos) y fija si debe cambiarla en el siguiente acceso.
func (r *Usuarios) CambiarHash(ctx context.Context, id int64, hash string, debeCambiar bool, ahora time.Time) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE USUARIO SET HASH_CONTRASENA = ?, DEBE_CAMBIAR_CONTRASENA = ?, ACTUALIZADO_EN = ?
		WHERE ID = ?`, hash, debeCambiar, ahora, id)
	return err
}
