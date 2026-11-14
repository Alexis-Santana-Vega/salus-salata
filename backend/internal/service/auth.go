// Package service contiene las reglas de negocio. No conoce HTTP: recibe
// datos ya validados en forma y devuelve modelos o errores de dominio.
package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"
	"unicode/utf8"

	"salus/backend/internal/model"
	"salus/backend/internal/repository"
	"salus/backend/internal/security"
)

// Política anti fuerza bruta (ventana deslizante sobre AUDITORIA_ACCESO).
const (
	VentanaIntentos      = 15 * time.Minute
	MaxFallidosPorCuenta = 5  // por usuario/correo (exista o no)
	MaxFallidosPorIP     = 20 // por IP, sumando todas las cuentas
	MinLargoContrasena   = 12
	MaxLargoContrasena   = 128 // tope para que argon2 no se use como DoS
)

// Errores de dominio; el handler los traduce a respuestas HTTP.
var (
	ErrCredenciales       = errors.New("credenciales inválidas")
	ErrDemasiadosIntentos = errors.New("demasiados intentos")
	ErrSesionInvalida     = errors.New("sesión inválida")
)

// ErrContrasenaDebil describe por qué una contraseña no cumple la política.
type ErrContrasenaDebil struct{ Motivo string }

func (e ErrContrasenaDebil) Error() string { return e.Motivo }

// Sesion es lo que recibe el cliente al iniciar sesión o renovarla.
type Sesion struct {
	AccessToken   string
	AccessExpira  time.Time
	RefreshToken  string // va en cookie httpOnly, nunca en el cuerpo
	RefreshExpira time.Time
	Usuario       model.Usuario
}

type Auth struct {
	db          *sql.DB
	tokens      *security.Tokens
	refreshVida time.Duration
	hashFalso   string // para igualar tiempos cuando el usuario no existe
	log         *slog.Logger
	ahora       func() time.Time
}

func NewAuth(db *sql.DB, tokens *security.Tokens, refreshVida time.Duration, log *slog.Logger) (*Auth, error) {
	falso, err := security.HashPassword("hash-falso-para-igualar-tiempos")
	if err != nil {
		return nil, err
	}
	return &Auth{db: db, tokens: tokens, refreshVida: refreshVida, hashFalso: falso, log: log,
		ahora: func() time.Time { return time.Now().UTC() }}, nil
}

// Normalizar deja usuario/correo en minúsculas y sin espacios.
func Normalizar(s string) string { return strings.ToLower(strings.TrimSpace(s)) }

// Login valida credenciales con límites por cuenta y por IP. Ante usuario
// inexistente, contraseña incorrecta o cuenta inactiva responde lo mismo, y
// tarda lo mismo, para no revelar qué cuentas existen.
func (a *Auth) Login(ctx context.Context, ident, contrasena string, o model.Origen) (*Sesion, error) {
	ident = Normalizar(ident)
	ahora := a.ahora()
	usuarios := repository.NewUsuarios(a.db)
	auditoria := repository.NewAuditoria(a.db)

	// 1. ¿Se superó el límite de intentos fallidos recientes?
	porIdent, porIP, err := auditoria.FallidosDesde(ctx, ident, o.IP, ahora.Add(-VentanaIntentos))
	if err != nil {
		return nil, err
	}
	if porIdent >= MaxFallidosPorCuenta || porIP >= MaxFallidosPorIP {
		a.auditar(ctx, model.EventoLoginBloqueado, 0, ident, o)
		return nil, ErrDemasiadosIntentos
	}

	// 2. Buscar y verificar. Si no existe, se verifica contra un hash falso
	//    para que la respuesta tarde lo mismo.
	u, err := usuarios.PorIdentificador(ctx, ident)
	if err != nil && !errors.Is(err, repository.ErrNoEncontrado) {
		return nil, err
	}
	hash := a.hashFalso
	if u != nil {
		hash = u.HashContrasena
	}
	ok, err := security.VerifyPassword(contrasena, hash)
	if err != nil {
		return nil, fmt.Errorf("verificar contraseña: %w", err)
	}
	if u == nil || !ok || !u.Activo {
		var uid int64
		if u != nil {
			uid = u.ID
		}
		a.auditar(ctx, model.EventoLoginFallido, uid, ident, o)
		return nil, ErrCredenciales
	}

	// 3. Si los parámetros de argon2 cambiaron, se aprovecha para rehashear.
	if security.NeedsRehash(u.HashContrasena) {
		if nuevo, err := security.HashPassword(contrasena); err == nil {
			_ = usuarios.CambiarHash(ctx, u.ID, nuevo, u.DebeCambiarContrasena, ahora)
		}
	}

	// 4. Abrir sesión nueva (familia nueva) y registrar el acceso.
	familia, err := security.RandomHex(16)
	if err != nil {
		return nil, err
	}
	sesiones := repository.NewSesiones(a.db)
	_ = sesiones.BorrarVencidas(ctx, u.ID, ahora)
	s, err := a.emitir(ctx, sesiones, u, familia, o, ahora)
	if err != nil {
		return nil, err
	}
	_ = usuarios.RegistrarAcceso(ctx, u.ID, ahora)
	a.auditar(ctx, model.EventoLoginOK, u.ID, ident, o)
	return s, nil
}

// Refresh rota el refresh token: el usado se revoca y se emite otro de la
// misma familia. Si llega un token ya revocado, alguien lo copió: se cierra
// toda la familia (el ladrón y el usuario legítimo deben volver a entrar).
func (a *Auth) Refresh(ctx context.Context, refresh string, o model.Origen) (*Sesion, error) {
	if refresh == "" {
		return nil, ErrSesionInvalida
	}
	ahora := a.ahora()

	tx, err := a.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()
	sesiones := repository.NewSesiones(tx)

	actual, err := sesiones.PorTokenHash(ctx, security.HashToken(refresh))
	if errors.Is(err, repository.ErrNoEncontrado) {
		return nil, ErrSesionInvalida
	}
	if err != nil {
		return nil, err
	}

	if actual.RevocadaEn != nil {
		// Reuso detectado: se revoca la familia fuera de esta transacción.
		_ = tx.Rollback()
		_ = repository.NewSesiones(a.db).RevocarFamilia(ctx, actual.Familia, ahora)
		a.auditar(ctx, model.EventoRefreshReuso, actual.UsuarioID, "", o)
		return nil, ErrSesionInvalida
	}
	if !ahora.Before(actual.ExpiraEn) {
		return nil, ErrSesionInvalida
	}

	u, err := repository.NewUsuarios(tx).PorID(ctx, actual.UsuarioID)
	if err != nil || !u.Activo {
		return nil, ErrSesionInvalida
	}

	// Revocar solo si sigue vigente: si otra petición lo rotó primero, esta pierde.
	revocada, err := sesiones.Revocar(ctx, actual.ID, ahora)
	if err != nil {
		return nil, err
	}
	if !revocada {
		return nil, ErrSesionInvalida
	}
	s, err := a.emitir(ctx, sesiones, u, actual.Familia, o, ahora)
	if err != nil {
		return nil, err
	}
	return s, tx.Commit()
}

// Logout revoca la sesión del refresh token (si existe). Es idempotente.
func (a *Auth) Logout(ctx context.Context, refresh string, o model.Origen) error {
	if refresh == "" {
		return nil
	}
	sesiones := repository.NewSesiones(a.db)
	s, err := sesiones.PorTokenHash(ctx, security.HashToken(refresh))
	if errors.Is(err, repository.ErrNoEncontrado) {
		return nil
	}
	if err != nil {
		return err
	}
	if _, err := sesiones.Revocar(ctx, s.ID, a.ahora()); err != nil {
		return err
	}
	a.auditar(ctx, model.EventoLogout, s.UsuarioID, "", o)
	return nil
}

// CambiarContrasena verifica la actual, guarda la nueva, cierra todas las
// sesiones del usuario y abre una nueva para el dispositivo actual.
func (a *Auth) CambiarContrasena(ctx context.Context, usuarioID int64, actual, nueva string, o model.Origen) (*Sesion, error) {
	ahora := a.ahora()
	usuarios := repository.NewUsuarios(a.db)
	u, err := usuarios.PorID(ctx, usuarioID)
	if err != nil || !u.Activo {
		return nil, ErrSesionInvalida
	}
	ok, err := security.VerifyPassword(actual, u.HashContrasena)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrCredenciales
	}
	if err := ValidarContrasena(nueva, u.Usuario, u.Correo); err != nil {
		return nil, err
	}
	if nueva == actual {
		return nil, ErrContrasenaDebil{Motivo: "La contraseña nueva debe ser distinta de la actual"}
	}

	hash, err := security.HashPassword(nueva)
	if err != nil {
		return nil, err
	}
	if err := usuarios.CambiarHash(ctx, u.ID, hash, false, ahora); err != nil {
		return nil, err
	}
	sesiones := repository.NewSesiones(a.db)
	if err := sesiones.RevocarDeUsuario(ctx, u.ID, ahora); err != nil {
		return nil, err
	}
	u.DebeCambiarContrasena = false
	familia, err := security.RandomHex(16)
	if err != nil {
		return nil, err
	}
	s, err := a.emitir(ctx, sesiones, u, familia, o, ahora)
	if err != nil {
		return nil, err
	}
	a.auditar(ctx, model.EventoCambioClave, u.ID, "", o)
	return s, nil
}

// Usuario devuelve el usuario de la sesión actual.
func (a *Auth) Usuario(ctx context.Context, id int64) (*model.Usuario, error) {
	u, err := repository.NewUsuarios(a.db).PorID(ctx, id)
	if errors.Is(err, repository.ErrNoEncontrado) || (err == nil && !u.Activo) {
		return nil, ErrSesionInvalida
	}
	return u, err
}

// emitir crea el refresh token (guardando solo su hash) y firma el de acceso.
func (a *Auth) emitir(ctx context.Context, sesiones *repository.Sesiones, u *model.Usuario, familia string, o model.Origen, ahora time.Time) (*Sesion, error) {
	refresh, err := security.RandomHex(32)
	if err != nil {
		return nil, err
	}
	expira := ahora.Add(a.refreshVida)
	if err := sesiones.Crear(ctx, model.Sesion{
		UsuarioID: u.ID, Familia: familia, TokenHash: security.HashToken(refresh),
		CreadaEn: ahora, ExpiraEn: expira,
	}, o); err != nil {
		return nil, err
	}
	access, err := a.tokens.Firmar(u.ID, u.Rol, u.DebeCambiarContrasena, ahora)
	if err != nil {
		return nil, err
	}
	return &Sesion{
		AccessToken: access, AccessExpira: ahora.Add(a.tokens.Vida()),
		RefreshToken: refresh, RefreshExpira: expira, Usuario: *u,
	}, nil
}

// auditar registra el evento; si falla, solo se anota en el log para no
// romper el acceso del usuario por un problema de bitácora.
func (a *Auth) auditar(ctx context.Context, evento string, usuarioID int64, ident string, o model.Origen) {
	if err := repository.NewAuditoria(a.db).Registrar(ctx, evento, usuarioID, ident, o, a.ahora()); err != nil {
		a.log.Error("no se pudo auditar", "evento", evento, "error", err)
	}
}

// ValidarContrasena aplica la política: largo 12–128 y que no contenga el
// usuario ni el correo. Sigue NIST 800-63B: largo antes que reglas raras.
func ValidarContrasena(c, usuario, correo string) error {
	n := utf8.RuneCountInString(c)
	switch {
	case n < MinLargoContrasena:
		return ErrContrasenaDebil{Motivo: fmt.Sprintf("Usa al menos %d caracteres", MinLargoContrasena)}
	case n > MaxLargoContrasena:
		return ErrContrasenaDebil{Motivo: fmt.Sprintf("Usa como máximo %d caracteres", MaxLargoContrasena)}
	case strings.TrimSpace(c) == "":
		return ErrContrasenaDebil{Motivo: "La contraseña no puede ser solo espacios"}
	}
	baja := strings.ToLower(c)
	if (usuario != "" && strings.Contains(baja, strings.ToLower(usuario))) ||
		(correo != "" && strings.Contains(baja, strings.ToLower(correo))) {
		return ErrContrasenaDebil{Motivo: "La contraseña no debe contener tu usuario ni tu correo"}
	}
	return nil
}
