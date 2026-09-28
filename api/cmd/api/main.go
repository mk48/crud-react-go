package main

import (
	"context"
	"kfamily/internal/services"
	"kfamily/internal/util"
	"kfamily/migrations"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/arl/statsviz"
	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func main() {
	env := util.ReadAppENVs()

	//DB open
	db, err := openDB(env)
	if err != nil {
		log.Fatal("Error while connecting DB.", err)
	}
	defer db.Close()

	//migration
	err = migrations.Run(db.DB)
	if err != nil {
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
	e.Use(middleware.RequestID())
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS(env.CorsAllowedOrigins...))

	services.CreateServices(context.Background(), db, e, env)

	// statsviz exposes runtime internals and is opened directly in a browser
	// (so it can't carry a Bearer token) - only serve it in dev.
	if env.Env == "dev" {
		setupStateViz(e)
	}

	//-------------- Graceful shutdown --------------
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	sc := echo.StartConfig{
		Address:         ":8080",
		GracefulTimeout: 10 * time.Second,
	}
	// A graceful shutdown returns nil - any error here means the server
	// failed (e.g. port in use), so exit non-zero for the orchestrator.
	if err := sc.Start(ctx, e); err != nil {
		e.Logger.Error("failed to start server", "error", err)
		os.Exit(1)
	}

	e.Logger.Debug("shutting down")
}

func openDB(env *util.AppENV) (*sqlx.DB, error) {
	db, err := sqlx.Open("pgx", env.ConnectionString)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func setupStateViz(e *echo.Echo) {
	mux := http.NewServeMux()

	// Register statsviz handlers on the mux.
	statsviz.Register(mux)

	// Use echo WrapHandler to wrap statsviz ServeMux as echo HandleFunc
	e.GET("/kfamily-debug/statsviz", echo.WrapHandler(mux))
	// Serve static content for statsviz UI
	e.GET("/kfamily-debug/statsviz/*", echo.WrapHandler(mux))
}
