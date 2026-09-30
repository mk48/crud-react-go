package util

import (
	"fmt"
	"log"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type AppENV struct {
	Env              string `env:"ENV"`
	Port             int    `env:"PORT" envDefault:"8080"`
	ConnectionString string `env:"CONNECTION_STRING,notEmpty"`

	// Postgres connection pool. Keep DBMaxOpenConns x replicas below the
	// server's max_connections.
	DBMaxOpenConns    int           `env:"DB_MAX_OPEN_CONNS" envDefault:"20"`
	DBMaxIdleConns    int           `env:"DB_MAX_IDLE_CONNS" envDefault:"5"`
	DBConnMaxLifetime time.Duration `env:"DB_CONN_MAX_LIFETIME" envDefault:"30m"`
	DBConnMaxIdleTime time.Duration `env:"DB_CONN_MAX_IDLE_TIME" envDefault:"5m"`

	// Casdoor application used to authenticate users. Endpoint/ClientId are
	// also used by the frontend to build the sign-in redirect; ClientSecret
	// and Certificate stay server-side.
	CasdoorEndpoint         string `env:"CASDOOR_ENDPOINT,notEmpty"`
	CasdoorClientId         string `env:"CASDOOR_CLIENT_ID,notEmpty"`
	CasdoorClientSecret     string `env:"CASDOOR_CLIENT_SECRET,notEmpty"`
	CasdoorOrganizationName string `env:"CASDOOR_ORGANIZATION_NAME,notEmpty"`
	CasdoorApplicationName  string `env:"CASDOOR_APPLICATION_NAME,notEmpty"`
	// The PEM certificate (from Casdoor's Cert management page) used to
	// verify the signature of tokens Casdoor issues. Newlines in the .env
	// value must be escaped as \n.
	CasdoorCertificate string `env:"CASDOOR_CERTIFICATE,notEmpty"`

	// Other browser origins (scheme://host[:port]) allowed to call the API,
	// comma-separated. Normally empty: the web app is served by this binary
	// (and proxied by Vite in dev), so it's always same-origin.
	CorsAllowedOrigins []string `env:"CORS_ALLOWED_ORIGINS" envSeparator:","`

	// Set by main from the build (-ldflags "-X main.version=...").
	Version string
}

func (e *AppENV) IsDev() bool {
	return e.Env == "dev"
}

func ReadAppENVs(version string) *AppENV {
	cfg := AppENV{Version: version}

	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		//normally .env file is not available in production (that will be set by CI/CD), so we can avoid this error
		fmt.Printf("Error loading .env file. Using system's environment variables. Err: %v\n", err)
	}

	//Read environment variables
	if err := env.Parse(&cfg); err != nil {
		log.Fatal("Error in parsing environment variables. ", err)
	}

	return &cfg
}
