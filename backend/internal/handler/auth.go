package handler

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v3"

	"salus/backend/internal/middleware"
	"salus/backend/internal/model"
	"salus/backend/internal/response"
	"salus/backend/internal/service"
)

// Cookie del refresh token: httpOnly (JS no la lee), SameSite=Strict (no
// viaja desde otros sitios) y limitada a las rutas de /auth.
const (
	cookieRefresh = "salus_refresh"
	cookiePath    = "/api/v1/auth"
)

type AuthHandler struct {
	Svc          *service.Auth
	CookieSegura bool // true en producción (solo HTTPS)
}

type loginReq struct {
	Usuario    string `json:"usuario" validate:"required,max=254"`
	Contrasena string `json:"contrasena" validate:"required,max=128"`
}

type cambioReq struct {
	ContrasenaActual string `json:"contrasenaActual" validate:"required,max=128"`
	ContrasenaNueva  string `json:"contrasenaNueva" validate:"required,max=128"`
}

// sesionResp es la respuesta de login, refresh y cambio de contraseña.
// El refresh token NO va aquí: viaja solo en la cookie.
type sesionResp struct {
	AccessToken string        `json:"accessToken"`
	TokenType   string        `json:"tokenType"`
	ExpiresIn   int           `json:"expiresIn"` // segundos
	Usuario     model.Usuario `json:"usuario"`
}

// Login: POST /api/v1/auth/login
func (h AuthHandler) Login(c fiber.Ctx) error {
	var req loginReq
	if err := bind(c, &req); err != nil {
		return err
	}
	s, err := h.Svc.Login(c.Context(), req.Usuario, req.Contrasena, origen(c))
	if err != nil {
		return traducir(err)
	}
	return h.responderSesion(c, s)
}

// Refresh: POST /api/v1/auth/refresh (usa la cookie, sin cuerpo).
func (h AuthHandler) Refresh(c fiber.Ctx) error {
	s, err := h.Svc.Refresh(c.Context(), c.Cookies(cookieRefresh), origen(c))
	if err != nil {
		h.borrarCookie(c)
		return traducir(err)
	}
	return h.responderSesion(c, s)
}

// Logout: POST /api/v1/auth/logout → 204 siempre (idempotente).
func (h AuthHandler) Logout(c fiber.Ctx) error {
	if err := h.Svc.Logout(c.Context(), c.Cookies(cookieRefresh), origen(c)); err != nil {
		return err
	}
	h.borrarCookie(c)
	return response.NoContent(c)
}

// Me: GET /api/v1/auth/me → usuario de la sesión.
func (h AuthHandler) Me(c fiber.Ctx) error {
	u, err := h.Svc.Usuario(c.Context(), middleware.Sesion(c).UsuarioID)
	if err != nil {
		return traducir(err)
	}
	return response.OK(c, u)
}

// CambiarContrasena: PUT /api/v1/auth/contrasena → nueva sesión.
func (h AuthHandler) CambiarContrasena(c fiber.Ctx) error {
	var req cambioReq
	if err := bind(c, &req); err != nil {
		return err
	}
	s, err := h.Svc.CambiarContrasena(c.Context(), middleware.Sesion(c).UsuarioID,
		req.ContrasenaActual, req.ContrasenaNueva, origen(c))
	if err != nil {
		var debil service.ErrContrasenaDebil
		switch {
		case errors.As(err, &debil):
			return response.Validation(response.FieldError{Field: "contrasenaNueva", Message: debil.Motivo})
		case errors.Is(err, service.ErrCredenciales):
			return response.Validation(response.FieldError{Field: "contrasenaActual", Message: "La contraseña actual no es correcta"})
		}
		return traducir(err)
	}
	return h.responderSesion(c, s)
}

func (h AuthHandler) responderSesion(c fiber.Ctx, s *service.Sesion) error {
	c.Cookie(&fiber.Cookie{
		Name: cookieRefresh, Value: s.RefreshToken, Path: cookiePath,
		Expires: s.RefreshExpira, HTTPOnly: true, Secure: h.CookieSegura, SameSite: fiber.CookieSameSiteStrictMode,
	})
	c.Set(fiber.HeaderCacheControl, "no-store") // los tokens no se cachean
	return response.OK(c, sesionResp{
		AccessToken: s.AccessToken,
		TokenType:   "Bearer",
		ExpiresIn:   int(time.Until(s.AccessExpira).Seconds()),
		Usuario:     s.Usuario,
	})
}

func (h AuthHandler) borrarCookie(c fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name: cookieRefresh, Value: "", Path: cookiePath, MaxAge: -1,
		Expires: time.Unix(0, 0), HTTPOnly: true, Secure: h.CookieSegura, SameSite: fiber.CookieSameSiteStrictMode,
	})
}

// traducir convierte errores de dominio en respuestas con mensaje amigable.
func traducir(err error) error {
	switch {
	case errors.Is(err, service.ErrCredenciales):
		return response.New(fiber.StatusUnauthorized, response.CodeInvalidCredentials, "Usuario o contraseña incorrectos")
	case errors.Is(err, service.ErrDemasiadosIntentos):
		return response.TooManyRequests("Demasiados intentos fallidos. Espera 15 minutos e intenta de nuevo.")
	case errors.Is(err, service.ErrSesionInvalida):
		return response.Unauthorized("Tu sesión terminó. Inicia sesión de nuevo.")
	}
	return err
}

func origen(c fiber.Ctx) model.Origen {
	return model.Origen{IP: c.IP(), UserAgent: c.Get(fiber.HeaderUserAgent)}
}
