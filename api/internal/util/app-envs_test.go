package util

import "testing"

func TestCheckOrigin(t *testing.T) {
	valid := []string{"https://admin.example.com", "http://localhost:5173", "https://example.com:8443"}
	for _, o := range valid {
		if err := CheckOrigin(o); err != nil {
			t.Errorf("CheckOrigin(%q) = %v, want ok", o, err)
		}
	}

	invalid := []string{
		"*", "", "example.com", "https://example.com/", "https://example.com/path",
		"ftp://example.com", "https://example.com?x=1", "https://user@example.com", "https://",
	}
	for _, o := range invalid {
		if err := CheckOrigin(o); err == nil {
			t.Errorf("CheckOrigin(%q) succeeded, want error", o)
		}
	}
}

func TestOriginOf(t *testing.T) {
	got, err := OriginOf("https://otel.example.com:4318/v1/traces")
	if err != nil || got != "https://otel.example.com:4318" {
		t.Errorf("OriginOf = %q, %v", got, err)
	}
	if _, err := OriginOf("/v1/traces"); err == nil {
		t.Error("a relative URL must be rejected")
	}
}

func TestAppENVValidate(t *testing.T) {
	ok := AppENV{
		CorsAllowedOrigins: []string{"https://admin.example.com"},
		WebOtelTracesUrl:   "https://otel.example.com/v1/traces",
		WebOtelSampleRatio: 0.5,
		TraceUrlTemplate:   "https://jaeger.example.com/trace/{traceId}",
	}
	if err := ok.validate(); err != nil {
		t.Fatalf("valid config rejected: %v", err)
	}

	bad := map[string]func(e *AppENV){
		"wildcard origin":     func(e *AppENV) { e.CorsAllowedOrigins = []string{"*"} },
		"relative traces url": func(e *AppENV) { e.WebOtelTracesUrl = "/v1/traces" },
		"ratio above 1":       func(e *AppENV) { e.WebOtelSampleRatio = 1.5 },
		"negative ratio":      func(e *AppENV) { e.WebOtelSampleRatio = -0.1 },
		"template without id": func(e *AppENV) { e.TraceUrlTemplate = "https://jaeger.example.com/trace/" },
		"relative trace url":  func(e *AppENV) { e.TraceUrlTemplate = "/trace/{traceId}" },
	}
	for name, mutate := range bad {
		t.Run(name, func(t *testing.T) {
			e := ok
			mutate(&e)
			if err := e.validate(); err == nil {
				t.Error("validate succeeded, want error")
			}
		})
	}
}
