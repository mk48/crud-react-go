package util

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type AppENV struct {
	Env              string `env:"ENV"`
	ConnectionString string `env:"CONNECTION_STRING,notEmpty"`

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

	// Browser origins (scheme://host[:port]) allowed to call the API, e.g.
	// the web app's URL. Comma-separated in the env var.
	CorsAllowedOrigins []string `env:"CORS_ALLOWED_ORIGINS" envSeparator:"," envDefault:"http://localhost:5173"`

	Version string
}

func ReadAppENVs() *AppENV {
	cfg := AppENV{Version: "1.0.0"} //TODO: set version automatically

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
