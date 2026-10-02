// Package webui serves the React web app from the API binary - the built
// files (embedded with -tags embedui) plus /config.js, the runtime settings
// the app reads at startup.
package webui

import (
	"encoding/json"
	"io/fs"
	"net/http"
	"path"
	"strings"

	"kfamily/internal/util"

	"github.com/labstack/echo/v5"
)

// runtimeConfig is exposed to the browser as window.__KFAMILY_CONFIG__ - so
// public only. The Casdoor secret/certificate never leave the server.
type runtimeConfig struct {
	CasdoorEndpoint string `json:"casdoorEndpoint"`
	CasdoorClientId string `json:"casdoorClientId"`
	// The app is embedded in this binary, so its version is the API's -
	// sent back as X-Client-Version.
	Version string `json:"version"`
	// Browser tracing - see docs/tracing.md and AppENV.WebOtel*.
	OtelTracesUrl    string  `json:"otelTracesUrl,omitempty"`
	OtelSampleRatio  float64 `json:"otelSampleRatio"`
	TraceUrlTemplate string  `json:"traceUrlTemplate,omitempty"`
}

// Register adds GET /config.js and, when the web app is embedded, serves it
// for every other GET path not claimed by a more specific route.
func Register(e *echo.Echo, env *util.AppENV) {
	e.GET("/config.js", configHandler(env))

	if files := files(); files != nil {
		e.GET("/*", spaHandler(files))
	}
}

// configHandler serves the web app's runtime settings as a script, so one
// build works in every environment - they come from this server's env vars
// rather than being compiled into the bundle. It's a same-origin script
// file (not inline), so the CSP can stay script-src 'self'.
func configHandler(env *util.AppENV) echo.HandlerFunc {
	// json.Marshal escapes <, > and & - safe to embed in a script.
	body, err := json.Marshal(runtimeConfig{
		CasdoorEndpoint:  env.CasdoorEndpoint,
		CasdoorClientId:  env.CasdoorClientId,
		Version:          env.Version,
		OtelTracesUrl:    env.WebOtelTracesUrl,
		OtelSampleRatio:  env.WebOtelSampleRatio,
		TraceUrlTemplate: env.TraceUrlTemplate,
	})
	if err != nil {
		panic(err) // plain strings - can't fail
	}
	script := []byte("window.__KFAMILY_CONFIG__ = " + string(body) + ";\n")

	return func(c *echo.Context) error {
		c.Response().Header().Set("Cache-Control", "no-store")
		return c.Blob(http.StatusOK, "text/javascript; charset=utf-8", script)
	}
}

// spaHandler serves a file from the built app when one matches the path,
// and index.html otherwise - client-side routes like /samples/123 only
// exist in the browser router.
func spaHandler(files fs.FS) echo.HandlerFunc {
	return func(c *echo.Context) error {
		name := strings.TrimPrefix(path.Clean("/"+c.Param("*")), "/")

		// A mistyped API path is a JSON 404, not the HTML app shell.
		if name == "api" || strings.HasPrefix(name, "api/") {
			return c.JSON(http.StatusNotFound, util.HttpErrorMessage("Not found"))
		}

		if name != "" && name != "index.html" {
			if info, err := fs.Stat(files, name); err == nil && !info.IsDir() {
				if strings.HasPrefix(name, "assets/") {
					// Vite content-hashes everything under assets/.
					c.Response().Header().Set("Cache-Control", "public, max-age=31536000, immutable")
				} else {
					c.Response().Header().Set("Cache-Control", "no-cache")
				}
				return c.FileFS(name, files)
			}

			// A missing hashed asset (e.g. a stale tab after a deploy) must
			// 404 - answering with index.html would be cached as JS/CSS.
			if strings.HasPrefix(name, "assets/") {
				return echo.ErrNotFound
			}
		}

		// Always revalidated, so a deploy's new asset hashes are picked up.
		c.Response().Header().Set("Cache-Control", "no-cache")
		return c.FileFS("index.html", files)
	}
}
