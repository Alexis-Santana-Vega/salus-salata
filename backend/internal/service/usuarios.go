package service

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"salus/backend/internal/model"
	"salus/backend/internal/repository"
	"salus/backend/internal/security"
)

type Usuarios struct {
	db  *sql.DB
	log *slog.Logger
}

func NewUsuarios(db *sql.DB, log *slog.Logger) *Usuarios { return &Usuarios{db: db, log: log} }

// Listar devuelve una página de usuarios y el total.
func (s *Usuarios) Listar(ctx context.Context, offset, limite int) ([]model.Usuario, int, error) {
	return repository.NewUsuarios(s.db).Listar(ctx, offset, limite)
}

// SysadminInicial son los datos del primer Sysadmin (variables INITIAL_ADMIN_*).
type SysadminInicial struct {
	Usuario, Correo, Nombre, Contrasena string
}

// AsegurarSysadmin crea el primer Sysadmin solo si no existe ninguno activo.
// Queda marcado para cambiar la contraseña en su primer acceso. La
// contraseña nunca se escribe en el log.
func (s *Usuarios) AsegurarSysadmin(ctx context.Context, ini SysadminInicial) error {
	hay, err := repository.NewUsuarios(s.db).HaySysadminActivo(ctx)
	if err != nil || hay {
		return err
	}
	if ini.Usuario == "" || ini.Contrasena == "" {
		s.log.Warn("no hay ningún Sysadmin; define INITIAL_ADMIN_* en .env para crearlo")
		return nil
	}
	if err := ValidarContrasena(ini.Contrasena, ini.Usuario, ini.Correo); err != nil {
		return err
	}
	hash, err := security.HashPassword(ini.Contrasena)
	if err != nil {
		return err
	}

	// Se vuelve a comprobar dentro de la transacción: dos arranques a la vez
	// no crean dos Sysadmin (el segundo choca con el UNIQUE y se revierte).
	creado := false
	err = enTx(ctx, s.db, func(tx *sql.Tx) error {
		repo := repository.NewUsuarios(tx)
		hay, err := repo.HaySysadminActivo(ctx)
		if err != nil || hay {
			return err
		}
		_, err = repo.Crear(ctx, model.Usuario{
			Usuario: Normalizar(ini.Usuario), Correo: Normalizar(ini.Correo), Nombre: ini.Nombre,
			HashContrasena: hash, Rol: model.RolSysadmin, Activo: true, DebeCambiarContrasena: true,
		}, time.Now().UTC())
		creado = err == nil
		return err
	})
	if creado {
		s.log.Info("Sysadmin inicial creado; debe cambiar su contraseña al entrar", "usuario", Normalizar(ini.Usuario))
	}
	return err
}
