// Package model contiene las entidades del dominio tal como viajan entre capas.
package model

import "time"

// Roles del sistema (tabla ROL). Solo estos tres existen; User SALUS y
// User Doc están fuera de alcance.
const (
	RolSysadmin = "SYSADMIN"
	RolAdmin    = "ADMIN"
	RolDoctor   = "DOCTOR"
)

type Usuario struct {
	ID                    int64      `json:"id"`
	Usuario               string     `json:"usuario"`
	Correo                string     `json:"correo"`
	Nombre                string     `json:"nombre"`
	Rol                   string     `json:"rol"`
	Activo                bool       `json:"activo"`
	DebeCambiarContrasena bool       `json:"debeCambiarContrasena"`
	UltimoAcceso          *time.Time `json:"ultimoAcceso"`
	CreadoEn              time.Time  `json:"creadoEn"`
	HashContrasena        string     `json:"-"` // nunca sale en JSON
}

// Sesion es un refresh token emitido a un dispositivo.
type Sesion struct {
	ID         int64
	UsuarioID  int64
	Familia    string
	TokenHash  string
	CreadaEn   time.Time
	ExpiraEn   time.Time
	RevocadaEn *time.Time
}

// Eventos de AUDITORIA_ACCESO.
const (
	EventoLoginOK        = "LOGIN_OK"
	EventoLoginFallido   = "LOGIN_FALLIDO"
	EventoLoginBloqueado = "LOGIN_BLOQUEADO"
	EventoLogout         = "LOGOUT"
	EventoRefreshReuso   = "REFRESH_REUSO"
	EventoCambioClave    = "CAMBIO_CONTRASENA"
)

// Origen describe desde dónde llega una petición (para auditoría).
type Origen struct {
	IP        string
	UserAgent string
}
