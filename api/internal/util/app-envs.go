package util

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"
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
	// comma-separated - e.g. an admin app on its own domain. The web app is
	// served by this binary (and proxied by Vite in dev), so it's always
	// same-origin and needs no entry; native mobile apps and batch jobs
	// aren't browsers and need none either. No "*".
	CorsAllowedOrigins []string `env:"CORS_ALLOWED_ORIGINS" envSeparator:","`

	// Apps besides the web app allowed to call the API - mobile, admin,
	// batch jobs - as a JSON array (see util.NewClientRegistry). Each is a
	// Casdoor application; its name is recorded on every operation it does.
	ApiClients string `env:"API_CLIENTS"`

	// Browser tracing for the web app (see docs/tracing.md). The OTLP/HTTP
	// traces URL the browser exports its spans to, e.g.
	// https://otel.example.com/v1/traces. Empty: the browser still starts
	// each trace and passes it on to the API, but doesn't export its spans.
	WebOtelTracesUrl string `env:"WEB_OTEL_TRACES_URL"`
	// Share of browser-started traces to sample, 0-1. The API follows the
	// browser's decision for those requests.
	WebOtelSampleRatio float64 `env:"WEB_OTEL_SAMPLE_RATIO" envDefault:"1"`
	// Optional link from a trace id to your tracing UI, with {traceId} as
	// placeholder, e.g. https://jaeger.example.com/trace/{traceId}.
	TraceUrlTemplate string `env:"TRACE_URL_TEMPLATE"`

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

	if err := cfg.validate(); err != nil {
		log.Fatal("Invalid environment variables. ", err)
	}

	return &cfg
}

// validate checks the values env.Parse can't, so a bad deploy fails at
// startup instead of on the first request.
func (e *AppENV) validate() error {
	for _, origin := range e.CorsAllowedOrigins {
		if err := CheckOrigin(origin); err != nil {
			return fmt.Errorf("CORS_ALLOWED_ORIGINS: %w", err)
		}
	}

	if e.WebOtelTracesUrl != "" {
		if _, err := OriginOf(e.WebOtelTracesUrl); err != nil {
			return fmt.Errorf("WEB_OTEL_TRACES_URL: %w", err)
		}
	}

	if e.WebOtelSampleRatio < 0 || e.WebOtelSampleRatio > 1 {
		return fmt.Errorf("WEB_OTEL_SAMPLE_RATIO must be between 0 and 1 (got %v)", e.WebOtelSampleRatio)
	}

	if e.TraceUrlTemplate != "" {
		if !strings.Contains(e.TraceUrlTemplate, "{traceId}") {
			return errors.New("TRACE_URL_TEMPLATE must contain {traceId}")
		}
		if _, err := OriginOf(e.TraceUrlTemplate); err != nil {
			return fmt.Errorf("TRACE_URL_TEMPLATE: %w", err)
		}
	}

	return nil
}

// CheckOrigin accepts a browser origin exactly as browsers send it in the
// Origin header - scheme://host[:port], http(s), no path or trailing slash.
// A wildcard isn't allowed: every app gets an explicit entry.
func CheckOrigin(origin string) error {
	if origin == "*" {
		return errors.New(`"*" isn't allowed - list each origin`)
	}
	u, err := url.Parse(origin)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" ||
		u.Path != "" || u.RawQuery != "" || u.Fragment != "" || u.User != nil {
		return fmt.Errorf("%q is not an origin like https://app.example.com", origin)
	}
	return nil
}

// OriginOf returns the origin (scheme://host[:port]) of an absolute
// http(s) URL.
func OriginOf(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return "", fmt.Errorf("%q is not an absolute http(s) URL", rawURL)
	}
	return u.Scheme + "://" + u.Host, nil
}
