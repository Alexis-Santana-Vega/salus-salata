package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	_ "time/tzdata" // zonas horarias dentro del binario (la imagen distroless no las trae)

	"github.com/gofiber/fiber/v3"

	"salus/backend/internal/config"
	"salus/backend/internal/database"
	"salus/backend/internal/middleware"
	"salus/backend/internal/response"
	"salus/backend/internal/router"
	"salus/backend/internal/security"
	"salus/backend/internal/service"
)

// banner se incrusta en el binario al compilar; no depende de archivos en disco.
//
//go:embed banner.txt
var banner string

// version se puede fijar al compilar: -ldflags "-X main.version=1.2.3".
var version = "0.6.0-dev"

// main solo decide el modo: chequeo de salud para Docker o servidor normal.
func main() {
	healthcheck := flag.Bool("healthcheck", false, "consulta /api/v1/health y sale con 0 si está sano (para Docker)")
	flag.Parse()
	if *healthcheck {
		os.Exit(runHealthcheck())
	}

	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

// run arma la aplicación en orden: config → log → BD → migraciones → HTTP,
// y la apaga limpio al recibir Ctrl+C o SIGTERM (docker stop).
func run() error {
	fmt.Println(banner)

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("configuración inválida:\n%w", err)
	}
	log := newLogger(cfg)
	log.Info("iniciando SALUS API", "version", version, "env", cfg.Env)

	// ctx se cancela con la señal de apagado; todo el arranque lo respeta.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db, err := database.Open(ctx, cfg.DB, log)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := database.Migrate(ctx, db, log); err != nil {
		return fmt.Errorf("migraciones: %w", err)
	}

	// Servicios: se crean una vez y se comparten (son seguros para uso concurrente).
	tokens := security.NewTokens(cfg.Auth.JWTSecret, cfg.Auth.AccessTTL)
	authSvc, err := service.NewAuth(db, tokens, cfg.Auth.RefreshTTL, log)
	if err != nil {
		return err
	}
	usuariosSvc := service.NewUsuarios(db, log)
	ini := cfg.Auth.AdminInicial
	if err := usuariosSvc.AsegurarSysadmin(ctx, service.SysadminInicial{
		Usuario: ini.Usuario, Correo: ini.Correo, Nombre: ini.Nombre, Contrasena: ini.Contrasena,
	}); err != nil {
		return fmt.Errorf("Sysadmin inicial: %w", err)
	}

	app := fiber.New(fiber.Config{
		AppName:      "SALUS API",
		ErrorHandler: response.ErrorHandler(log),
		BodyLimit:    1 << 20, // 1 MB
		Concurrency:  2048,    // tope de conexiones simultáneas
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  30 * time.Second,
		// Detrás de nginx la IP real llega en X-Real-IP. Solo se lee si la
		// petición viene de un proxy de confianza; si no, se usa la IP directa
		// (así nadie la falsifica mandando la cabecera).
		TrustProxy:       len(cfg.TrustedProxies) > 0,
		TrustProxyConfig: fiber.TrustProxyConfig{Proxies: cfg.TrustedProxies},
		ProxyHeader:      "X-Real-IP",
	})
	middleware.Register(app, cfg, log)
	router.Setup(app, router.Deps{
		Config: cfg, DB: db, Version: version,
		Tokens: tokens, Auth: authSvc, Usuarios: usuariosSvc,
		Procedimientos: service.NewProcedimientos(db),
		Tecnicos:       service.NewTecnicos(db),
	})

	// Única gorutina: Listen bloquea, así que corre aparte y avisa por errCh
	// (con buffer, para que no se quede colgada si ya estamos apagando).
	errCh := make(chan error, 1)
	go func() {
		errCh <- app.Listen(":"+cfg.Port, fiber.ListenConfig{DisableStartupMessage: true})
	}()
	log.Info("SALUS API escuchando", "port", cfg.Port, "docs", cfg.DocsEnabled)

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
	}

	// Termina las peticiones en curso (máximo 10 s) antes de cerrar.
	log.Info("apagando…")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return app.ShutdownWithContext(shutdownCtx)
}

// newLogger: texto legible en desarrollo, JSON en producción (para que Docker
// o cualquier colector de logs lo pueda procesar).
func newLogger(cfg config.Config) *slog.Logger {
	opts := &slog.HandlerOptions{Level: cfg.LogLevel}
	var h slog.Handler = slog.NewTextHandler(os.Stdout, opts)
	if cfg.IsProduction() {
		h = slog.NewJSONHandler(os.Stdout, opts)
	}
	l := slog.New(h).With("service", "salus-api")
	slog.SetDefault(l)
	return l
}

// runHealthcheck lo usa Docker (la imagen no tiene curl): 0 = sano, 1 = no.
func runHealthcheck() int {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	client := http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get("http://127.0.0.1:" + port + "/api/v1/health")
	if err != nil {
		return 1
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 1
	}
	return 0
}
