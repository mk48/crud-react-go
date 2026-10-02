// Package telemetry wires up OpenTelemetry tracing: the SDK, the HTTP
// server/client and database instrumentation, and the glue that puts each
// request's trace id on its response, log lines and operations. See
// docs/tracing.md.
package telemetry

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"kfamily/internal/util"

	"github.com/XSAM/otelsql"
	"github.com/labstack/echo/v5"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.43.0"
	"go.opentelemetry.io/otel/trace"
)

// TraceIDHeader carries the request's trace id back to the caller, so an
// error message can quote it (and support can find the trace, logs and
// operations from it).
const TraceIDHeader = "X-Trace-Id"

const tracerName = "kfamily"

// Tracer is the tracer for hand-made spans.
func Tracer() trace.Tracer {
	return otel.Tracer(tracerName)
}

// Setup installs the global TracerProvider and the W3C trace-context and
// baggage propagators. Spans are always created - so every request has a
// trace id to return, log and store on its operations - but only exported
// when an OTLP endpoint is configured through the standard
// OTEL_EXPORTER_OTLP_ENDPOINT / OTEL_EXPORTER_OTLP_TRACES_ENDPOINT env vars
// (plus OTEL_EXPORTER_OTLP_HEADERS etc.). Sampling follows
// OTEL_TRACES_SAMPLER/_ARG, defaulting to parentbased_always_on.
//
// Call the returned shutdown before exiting, so buffered spans are flushed.
func Setup(ctx context.Context, env *util.AppENV, log *slog.Logger) (shutdown func(context.Context) error, err error) {
	res, err := resource.New(ctx,
		resource.WithTelemetrySDK(),
		resource.WithAttributes(
			semconv.ServiceName("kfamily-api"),
			semconv.ServiceVersion(env.Version),
			semconv.DeploymentEnvironmentNameKey.String(env.Env),
		),
		// Last, so OTEL_SERVICE_NAME / OTEL_RESOURCE_ATTRIBUTES win.
		resource.WithFromEnv(),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to build telemetry resource. err: %w", err)
	}

	opts := []sdktrace.TracerProviderOption{sdktrace.WithResource(res)}
	if exportEnabled() {
		exporter, err := otlptracehttp.New(ctx)
		if err != nil {
			return nil, fmt.Errorf("unable to create OTLP trace exporter. err: %w", err)
		}
		opts = append(opts, sdktrace.WithBatcher(exporter))
		log.Info("exporting traces over OTLP/HTTP")
	} else {
		log.Info("trace export off - set OTEL_EXPORTER_OTLP_ENDPOINT to enable it")
	}

	provider := sdktrace.NewTracerProvider(opts...)
	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{}, propagation.Baggage{},
	))

	return provider.Shutdown, nil
}

func exportEnabled() bool {
	if strings.EqualFold(os.Getenv("OTEL_SDK_DISABLED"), "true") || os.Getenv("OTEL_TRACES_EXPORTER") == "none" {
		return false
	}
	return os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT") != "" || os.Getenv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT") != ""
}

// HTTPHandler wraps the whole server in a span per /api request, continuing
// the caller's trace when it sends a traceparent header. Health checks,
// /config.js and the web app's static files aren't traced.
func HTTPHandler(h http.Handler) http.Handler {
	return otelhttp.NewHandler(h, "http.server",
		otelhttp.WithFilter(func(r *http.Request) bool {
			return strings.HasPrefix(r.URL.Path, "/api/")
		}),
		// "<method> <route>", e.g. "PUT /api/v1/samples/:id" - low
		// cardinality, unlike the raw path. Echo sets the route as the
		// request's Pattern once it has routed, and otelhttp renames the span
		// with this after the handler; until then it's just the method.
		otelhttp.WithSpanNameFormatter(func(_ string, r *http.Request) string {
			if r.Pattern == "" {
				return r.Method
			}
			if strings.Contains(r.Pattern, " ") { // net/http style "GET /path"
				return r.Pattern
			}
			return r.Method + " " + r.Pattern
		}),
	)
}

// Middleware records the request's Echo route on its span (http.route),
// returns the trace id in X-Trace-Id and tags the request's logger with
// trace_id/span_id. Register it with e.Use (which Echo runs after routing),
// with HTTPHandler wrapping the server.
func Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			span := trace.SpanFromContext(c.Request().Context())
			sc := span.SpanContext()
			if sc.IsValid() {
				if route := c.Path(); route != "" {
					span.SetAttributes(semconv.HTTPRoute(route))
				}
				c.Response().Header().Set(TraceIDHeader, sc.TraceID().String())
				c.SetLogger(c.Logger().With(
					slog.String("trace_id", sc.TraceID().String()),
					slog.String("span_id", sc.SpanID().String()),
				))
			}
			return next(c)
		}
	}
}

// OpenDB opens driverName/dsn with a span per query, inside the request's
// trace. Only the parameterized SQL is recorded - never the argument values,
// which hold personal data. Database work outside any trace (migrations, the
// pool's own housekeeping) isn't traced.
func OpenDB(driverName, dsn string) (*sql.DB, error) {
	return otelsql.Open(driverName, dsn,
		otelsql.WithAttributes(semconv.DBSystemNamePostgreSQL),
		otelsql.WithSpanOptions(otelsql.SpanOptions{
			OmitConnResetSession: true,
			OmitConnPrepare:      true,
			OmitRows:             true,
			OmitConnectorConnect: true,
			SpanFilter: func(ctx context.Context, _ otelsql.Method, _ string, _ []driver.NamedValue) bool {
				return trace.SpanContextFromContext(ctx).IsValid()
			},
			// A lookup finding nothing is an answer, not a failure.
			RecordError: func(err error) bool {
				return !errors.Is(err, sql.ErrNoRows)
			},
		}),
	)
}

// HTTPClient returns a client whose requests are spans in ctx's trace (and
// carry its traceparent) - for SDKs that build their own requests without
// our context, like casdoorsdk.GetOAuthToken.
func HTTPClient(ctx context.Context) *http.Client {
	return &http.Client{
		Timeout: 15 * time.Second,
		Transport: contextTransport{
			span: trace.SpanFromContext(ctx),
			next: otelhttp.NewTransport(http.DefaultTransport),
		},
	}
}

// contextTransport makes each request a child of span, whatever context the
// request was built with.
type contextTransport struct {
	span trace.Span
	next http.RoundTripper
}

func (t contextTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	return t.next.RoundTrip(r.WithContext(trace.ContextWithSpan(r.Context(), t.span)))
}
