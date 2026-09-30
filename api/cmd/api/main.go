// KFamily API
//
//	@title						KFamily API
//	@version					1.0
//	@description				Backend for the KFamily family-tree app. Every /api/v1 endpoint except /auth/signin needs a Casdoor access token.
//	@BasePath					/
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description				"Bearer <Casdoor access token>"
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	_ "kfamily/internal/docs" // registers the generated Swagger spec
	"kfamily/internal/services"
	"kfamily/internal/util"
	"kfamily/internal/webui"
	"kfamily/migrations"

	"github.com/arl/statsviz"
	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

//go:generate go run github.com/swaggo/swag/cmd/swag@v1.16.4 init --generalInfo main.go --dir ./,../../internal --output ../../internal/docs --outputTypes go,json --parseInternal --parseDependency

// version is set at build time: -ldflags "-X main.version=<version>".
var version = "dev"

// contentSecurityPolicy only allows this origin's own scripts, so an
// injected <script> can't run (and read the token in localStorage).
// 'unsafe-inline' styles are needed for the UI components' style attributes.
const contentSecurityPolicy = "default-src 'self'; script-src 'self'; " +
	"style-src 'self' 'unsafe-inline'; img-src 'self' data:; font-src 'self' data:; " +
	"connect-src 'self'; object-src 'none'; base-uri 'self'; " +
	"frame-ancestors 'none'; form-action 'self'"

func main() {
	// `kfamily healthcheck` - for the container HEALTHCHECK; the distroless
	// image has no curl/wget.
	if len(os.Args) > 1 && os.Args[1] == "healthcheck" {
		os.Exit(healthcheck())
	}

	env := util.ReadAppENVs(version)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	//DB open
	db, err := openDB(ctx, env)
	if err != nil {
		log.Fatal("Error while connecting DB.", err)
	}
	defer db.Close()

	//migration
	if err := migrations.Run(ctx, db.DB); err != nil {
		log.Fatal("Error while executing migration:", err)
	}

	//casdoor setup: tokens are verified locally against this certificate, no
	//network call needed per request.
	casdoorsdk.InitConfig(
		env.CasdoorEndpoint,
		env.CasdoorClientId,
		env.CasdoorClientSecret,
		env.CasdoorCertificate,
		env.CasdoorOrganizationName,
		env.CasdoorApplicationName,
	)

	// create Echo server
	e := echo.New()
	// Behind Dokploy's Traefik: take the client IP from X-Forwarded-For, but
	// only as set by proxies on private/loopback addresses - so a client
	// can't spoof its IP (e.g. to dodge the sign-in rate limit).
	e.IPExtractor = echo.ExtractIPFromXFFHeader()

	e.Use(middleware.RequestID())
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	// Every request body here is small JSON.
	e.Use(middleware.BodyLimit(1 << 20)) // 1 MiB
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         "0", // legacy filter; CSP replaces it
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "DENY",
		HSTSMaxAge:            31536000, // only sent over HTTPS
		ContentSecurityPolicy: contentSecurityPolicy,
		ReferrerPolicy:        "strict-origin-when-cross-origin",
		// The dev-only tool pages below use inline scripts.
		Skipper: func(c *echo.Context) bool {
			p := c.Request().URL.Path
			return strings.HasPrefix(p, "/swagger/") || strings.HasPrefix(p, "/kfamily-debug/")
		},
	}))
	if len(env.CorsAllowedOrigins) > 0 {
		e.Use(middleware.CORS(env.CorsAllowedOrigins...))
	}

	services.CreateServices(ctx, db, e, env)
	webui.Register(e, env)

	// statsviz and Swagger expose internals and are opened directly in a
	// browser (so they can't carry a Bearer token) - only serve them in dev.
	if env.IsDev() {
		setupStateViz(e)
		e.GET("/swagger/*", echo.WrapHandler(httpSwagger.Handler(httpSwagger.URL("/swagger/doc.json"))))
	}

	//-------------- Graceful shutdown --------------
	sc := echo.StartConfig{
		Address:         fmt.Sprintf(":%d", env.Port),
		GracefulTimeout: 10 * time.Second,
		BeforeServeFunc: func(s *http.Server) error {
			// Without these, slow or idle clients can hold connections open
			// indefinitely (e.g. slowloris).
			s.ReadHeaderTimeout = 10 * time.Second
			s.ReadTimeout = 30 * time.Second
			s.IdleTimeout = 120 * time.Second
			// statsviz streams over a long-lived websocket in dev.
			if !env.IsDev() {
				s.WriteTimeout = 60 * time.Second
			}
			return nil
		},
	}
	// A graceful shutdown returns nil - any error here means the server
	// failed (e.g. port in use), so exit non-zero for the orchestrator.
	if err := sc.Start(ctx, e); err != nil {
		e.Logger.Error("failed to start server", "error", err)
		os.Exit(1)
	}

	e.Logger.Debug("shutting down")
}

func openDB(ctx context.Context, env *util.AppENV) (*sqlx.DB, error) {
	db, err := sqlx.Open("pgx", env.ConnectionString)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(env.DBMaxOpenConns)
	db.SetMaxIdleConns(env.DBMaxIdleConns)
	db.SetConnMaxLifetime(env.DBConnMaxLifetime)
	db.SetConnMaxIdleTime(env.DBConnMaxIdleTime)

	pingCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := db.PingContext(pingCtx); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

// healthcheck asks this container's own server whether it's ready (the
// database answers). Returns the process exit code: 0 healthy, 1 not.
func healthcheck() int {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get("http://127.0.0.1:" + port + "/healthcheck/ready")
	if err != nil {
		fmt.Fprintln(os.Stderr, "healthcheck:", err)
		return 1
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintln(os.Stderr, "healthcheck: status", resp.StatusCode)
		return 1
	}
	return 0
}

func setupStateViz(e *echo.Echo) {
	mux := http.NewServeMux()

	// Register statsviz handlers on the mux, under the same path Echo
	// forwards to it (its default root is /debug/statsviz).
	if err := statsviz.Register(mux, statsviz.Root("/kfamily-debug/statsviz")); err != nil {
		e.Logger.Error("unable to set up statsviz", "err", err)
	}

	// Use echo WrapHandler to wrap statsviz ServeMux as echo HandleFunc
	e.GET("/kfamily-debug/statsviz", echo.WrapHandler(mux))
	// Serve static content for statsviz UI
	e.GET("/kfamily-debug/statsviz/*", echo.WrapHandler(mux))
}
