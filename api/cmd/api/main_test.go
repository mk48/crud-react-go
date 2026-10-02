package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func corsServer() *echo.Echo {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(corsConfig([]string{"https://admin.example.com"})))
	e.PUT("/api/v1/samples/:id", func(c *echo.Context) error { return c.NoContent(http.StatusOK) })
	return e
}

func preflight(e *echo.Echo, origin string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodOptions, "/api/v1/samples/1", nil)
	req.Header.Set("Origin", origin)
	req.Header.Set("Access-Control-Request-Method", http.MethodPut)
	req.Header.Set("Access-Control-Request-Headers", "authorization,content-type,traceparent,x-client-version")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec
}

func TestCORSPreflightAllowedOrigin(t *testing.T) {
	rec := preflight(corsServer(), "https://admin.example.com")

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "https://admin.example.com" {
		t.Errorf("Allow-Origin = %q", got)
	}
	allowed := strings.ToLower(rec.Header().Get("Access-Control-Allow-Headers"))
	for _, h := range []string{"authorization", "content-type", "traceparent", "x-client-version"} {
		if !strings.Contains(allowed, h) {
			t.Errorf("Allow-Headers %q is missing %s", allowed, h)
		}
	}
	if rec.Header().Get("Access-Control-Allow-Credentials") != "" {
		t.Error("credentials must stay off - auth is a Bearer token")
	}
}

func TestCORSPreflightOtherOrigin(t *testing.T) {
	rec := preflight(corsServer(), "https://evil.example.com")

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Errorf("Allow-Origin = %q for an unlisted origin", got)
	}
}

func TestCORSExposesTraceID(t *testing.T) {
	req := httptest.NewRequest(http.MethodPut, "/api/v1/samples/1", nil)
	req.Header.Set("Origin", "https://admin.example.com")
	rec := httptest.NewRecorder()
	corsServer().ServeHTTP(rec, req)

	if got := rec.Header().Get("Access-Control-Expose-Headers"); !strings.Contains(got, "X-Trace-Id") {
		t.Errorf("Expose-Headers = %q, want X-Trace-Id", got)
	}
}
