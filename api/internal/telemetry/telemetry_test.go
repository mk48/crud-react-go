package telemetry

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
)

// newServer is the production wiring in miniature: HTTPHandler around an
// Echo app using Middleware, recording spans in memory.
func newServer(t *testing.T) (http.Handler, *tracetest.SpanRecorder) {
	t.Helper()
	recorder := tracetest.NewSpanRecorder()
	otel.SetTracerProvider(sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(recorder)))
	otel.SetTextMapPropagator(propagation.TraceContext{})

	e := echo.New()
	e.Use(Middleware())
	handler := func(c *echo.Context) error {
		if !trace.SpanContextFromContext(c.Request().Context()).IsValid() {
			return c.String(http.StatusOK, "untraced")
		}
		return c.String(http.StatusOK, "traced")
	}
	e.GET("/api/v1/samples/:id", handler)
	e.GET("/healthcheck", handler)

	return HTTPHandler(e), recorder
}

func TestContinuesCallersTrace(t *testing.T) {
	server, recorder := newServer(t)

	const traceID = "4bf92f3577b34da6a3ce929d0e0e4736"
	req := httptest.NewRequest(http.MethodGet, "/api/v1/samples/123", nil)
	req.Header.Set("traceparent", "00-"+traceID+"-00f067aa0ba902b7-01")
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)

	if got := rec.Header().Get(TraceIDHeader); got != traceID {
		t.Errorf("%s = %q, want the caller's trace id %q", TraceIDHeader, got, traceID)
	}

	spans := recorder.Ended()
	if len(spans) != 1 {
		t.Fatalf("got %d spans, want 1", len(spans))
	}
	span := spans[0]
	if span.Name() != "GET /api/v1/samples/:id" {
		t.Errorf("span name = %q, want the route, not the raw path", span.Name())
	}
	if span.Parent().SpanID().String() != "00f067aa0ba902b7" {
		t.Errorf("span parent = %s, want the caller's span", span.Parent().SpanID())
	}
	var route string
	for _, a := range span.Attributes() {
		if a.Key == "http.route" {
			route = a.Value.AsString()
		}
	}
	if route != "/api/v1/samples/:id" {
		t.Errorf("http.route = %q", route)
	}
}

func TestStartsTraceWithoutTraceparent(t *testing.T) {
	server, _ := newServer(t)

	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/v1/samples/123", nil))

	id, err := trace.TraceIDFromHex(rec.Header().Get(TraceIDHeader))
	if err != nil || !id.IsValid() {
		t.Errorf("%s = %q, want a new trace id", TraceIDHeader, rec.Header().Get(TraceIDHeader))
	}
}

func TestSkipsNonAPIPaths(t *testing.T) {
	server, recorder := newServer(t)

	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/healthcheck", nil))

	if rec.Body.String() != "untraced" || rec.Header().Get(TraceIDHeader) != "" || len(recorder.Ended()) != 0 {
		t.Errorf("health check was traced: body %q, header %q, %d spans",
			rec.Body.String(), rec.Header().Get(TraceIDHeader), len(recorder.Ended()))
	}
}
