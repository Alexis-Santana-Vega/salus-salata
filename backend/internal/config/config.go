// Package config carga y valida la configuración desde variables de entorno.
// Si falta algo obligatorio, la API no arranca y dice qué falta.
package config

import (
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

const (
	EnvDevelopment = "development"
	EnvProduction  = "production"
)

type Config struct {
	Env         string
	Port        string
	LogLevel    slog.Level
	DocsEnabled bool
	CORSOrigins []string
	DB          DBConfig
}

type DBConfig struct {
	Host     string
	Port     string
	Alias    string
	User     string
	Password string
}

func (c Config) IsProduction() bool { return c.Env == EnvProduction }

// Load lee el entorno. En desarrollo intenta cargar ../.env y .env sin
// sobrescribir variables que ya existan (docker compose tiene prioridad).
func Load() (Config, error) {
	// godotenv se detiene en el primer archivo que no existe: se cargan por separado.
	for _, f := range []string{".env", "../.env"} {
		_ = godotenv.Load(f)
	}

	cfg := Config{
		Env:         get("APP_ENV", EnvDevelopment),
		Port:        get("PORT", "8080"),
		DocsEnabled: getBool("DOCS_ENABLED", false),
		CORSOrigins: getList("CORS_ORIGINS"),
		DB: DBConfig{
			Host:     get("FIREBIRD_HOST", "127.0.0.1"),
			Port:     get("FIREBIRD_PORT", "3050"),
			Alias:    os.Getenv("FIREBIRD_DB_ALIAS"),
			User:     os.Getenv("FIREBIRD_USER"),
			Password: os.Getenv("FIREBIRD_PASSWORD"),
		},
	}

	// Se juntan todos los errores para reportarlos de una vez al arrancar.
	var errs []error
	if cfg.Env != EnvDevelopment && cfg.Env != EnvProduction {
		errs = append(errs, fmt.Errorf("APP_ENV debe ser %q o %q", EnvDevelopment, EnvProduction))
	}
	if err := cfg.LogLevel.UnmarshalText([]byte(get("LOG_LEVEL", "info"))); err != nil {
		errs = append(errs, errors.New("LOG_LEVEL debe ser debug, info, warn o error"))
	}
	if _, err := strconv.Atoi(cfg.Port); err != nil {
		errs = append(errs, errors.New("PORT debe ser numérico"))
	}
	if cfg.DB.User == "" {
		errs = append(errs, errors.New("falta FIREBIRD_USER"))
	}
	if cfg.DB.Password == "" {
		errs = append(errs, errors.New("falta FIREBIRD_PASSWORD"))
	}
	if err := validateAlias(cfg.DB.Alias); err != nil {
		errs = append(errs, err)
	}
	if strings.EqualFold(cfg.DB.User, "SYSDBA") {
		errs = append(errs, errors.New("la API no debe conectarse como SYSDBA; usa el usuario de la aplicación"))
	}
	if cfg.IsProduction() && cfg.DocsEnabled {
		errs = append(errs, errors.New("DOCS_ENABLED no se permite en producción"))
	}
	for _, o := range cfg.CORSOrigins {
		if err := validateOrigin(o, cfg.IsProduction()); err != nil {
			errs = append(errs, err)
		}
	}
	return cfg, errors.Join(errs...)
}

// validateOrigin acepta solo orígenes exactos (esquema + host + puerto opcional).
// Rechaza comodines y rutas; en producción exige https.
func validateOrigin(o string, prod bool) error {
	u, err := url.Parse(o)
	switch {
	case strings.Contains(o, "*"):
		return fmt.Errorf("CORS_ORIGINS no admite comodines: %q", o)
	case err != nil || u.Host == "" || (u.Scheme != "http" && u.Scheme != "https"):
		return fmt.Errorf("CORS_ORIGINS tiene un origen inválido: %q (ej. http://localhost:3000)", o)
	case u.Path != "" || u.RawQuery != "" || u.User != nil:
		return fmt.Errorf("CORS_ORIGINS debe ser solo esquema y host, sin ruta: %q", o)
	case prod && u.Scheme != "https":
		return fmt.Errorf("en producción CORS_ORIGINS debe usar https: %q", o)
	}
	return nil
}

// validateAlias obliga a usar un alias de databases.conf y nunca una ruta.
func validateAlias(alias string) error {
	if alias == "" {
		return errors.New("falta FIREBIRD_DB_ALIAS")
	}
	if strings.ContainsAny(alias, `/\:.`) {
		return errors.New("FIREBIRD_DB_ALIAS debe ser un alias (p. ej. salus), no una ruta de archivo")
	}
	return nil
}

// get devuelve la variable sin espacios o el valor por defecto si está vacía.
func get(key, def string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return def
}

// getBool interpreta true/false/1/0; si el valor no es válido usa el defecto.
func getBool(key string, def bool) bool {
	v, err := strconv.ParseBool(get(key, strconv.FormatBool(def)))
	if err != nil {
		return def
	}
	return v
}

// getList separa una variable por comas e ignora elementos vacíos.
func getList(key string) []string {
	var out []string
	for _, p := range strings.Split(os.Getenv(key), ",") {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}
