# KFamily API

Go backend for the KFamily family-tree app. Backed by Postgres and Casdoor
for auth.

This is a fresh start copied from an older project (`kuppai`) - only the
generic plumbing (auth, CRUD/audit-history framework, `user` table) and one
`sample_items` table survived, as a base to build the real family-tree
domain model on top of.

## Stack

- [Echo v5](https://github.com/labstack/echo) HTTP framework
- Postgres via `sqlx` + `pgx/v5` (driver name `pgx`)
- [sql-migrate](https://github.com/rubenv/sql-migrate)-style SQL migrations, embedded in the binary and run automatically at startup (`migrations.Run`)
- [Casdoor](https://casdoor.org) for authentication (JWT bearer tokens, verified locally against Casdoor's certificate - see `internal/services/auth` and `internal/middleware/auth.go`)
- `statsviz` for live runtime stats, Swagger (swaggo) for API docs - both dev only
- Serves the built web app (`../web`) too, embedded with `-tags embedui` - see `internal/webui` and the root README

## Getting started

```bash
go run ./cmd/api
```

The server listens on `:8080` by default (`PORT`). On startup it connects to
Postgres and runs any pending migrations (under an advisory lock, so
replicas starting together take turns).

Copy `.env-sample` to `.env` and fill in real values before running:

| Variable | Purpose |
|---|---|
| `ENV` | `dev` enables Swagger UI and statsviz; use `prod` otherwise |
| `PORT` | Listen port (default `8080`) |
| `CONNECTION_STRING` | Postgres connection string |
| `DB_MAX_OPEN_CONNS` / `DB_MAX_IDLE_CONNS` / `DB_CONN_MAX_LIFETIME` / `DB_CONN_MAX_IDLE_TIME` | Connection pool (defaults `20` / `5` / `30m` / `5m`) |
| `CORS_ALLOWED_ORIGINS` | Optional, comma-separated - other *browser* apps' origins (no `*`); native apps and batch jobs need none |
| `API_CLIENTS` | Optional JSON array of apps besides the web app (mobile, admin, batch jobs) - see [docs/tracing.md](../docs/tracing.md) |
| `OTEL_EXPORTER_OTLP_ENDPOINT` (and other standard `OTEL_*`) | Optional - export traces over OTLP/HTTP |
| `WEB_OTEL_TRACES_URL` / `WEB_OTEL_SAMPLE_RATIO` | Optional - where the browser exports its spans, and its sampling ratio (default `1`) |
| `TRACE_URL_TEMPLATE` | Optional - link from a trace id to your tracing UI, e.g. `http://localhost:16686/trace/{traceId}` |
| `CASDOOR_ENDPOINT` | Casdoor server URL |
| `CASDOOR_CLIENT_ID` / `CASDOOR_CLIENT_SECRET` | Casdoor application credentials |
| `CASDOOR_ORGANIZATION_NAME` / `CASDOOR_APPLICATION_NAME` | Casdoor organization/application names |
| `CASDOOR_CERTIFICATE` | PEM certificate (Casdoor's Cert management page) used to verify token signatures |

The web app gets `CASDOOR_ENDPOINT` and `CASDOOR_CLIENT_ID` from this server
at runtime (`GET /config.js`) - it has no configuration of its own.

## Auth flow

1. The frontend redirects the user to Casdoor's sign-in page.
2. Casdoor redirects back to the frontend's `/callback` with a `code`.
3. The frontend POSTs `{ code, state }` to `POST /api/v1/auth/signin`, which
   exchanges the code for a Casdoor access token (needs the client secret,
   so it has to happen server-side).
4. The frontend uses that access token as a Bearer token for every other API
   call. `AuthMiddleware` verifies it locally against `CASDOOR_CERTIFICATE`
   (no round-trip to Casdoor) and provisions a local `user` row from the
   token's claims on first sight of a given `sub`.

## Domain model

- **Users** — provisioned automatically on first sign-in (looked up by Casdoor `sub`, embedded in the verified token). New users default to `is_admin = false`; admin-only endpoints are gated server-side by `AdminMiddleware`.
- **Sample items** (`sample_items`) — a placeholder CRUD resource for exercising the auth/CRUD/audit-history plumbing end to end. Replace with the real family-tree tables once the schema is designed.

Most resources follow the same CRUD + audit-history + soft-delete shape; see
`internal/services/*` for the per-table service packages.

## Operations and tracing

Every write runs inside an *operation* (`util.RunOperation`) recording who
did it, from which app (`client`, verified from the access token) and the
request's OpenTelemetry trace id. The trace id is also returned as
`X-Trace-Id` and is on every log line. See
[docs/tracing.md](../docs/tracing.md) for how it fits together, how to add
an app (mobile, admin, batch job) and how to run a local trace viewer.

## Dev tools (`ENV=dev` only)

- Swagger UI: http://localhost:8080/swagger/index.html - regenerate the spec
  after changing handler annotations with `go generate ./cmd/api`
- Runtime stats: http://localhost:8080/kfamily-debug/statsviz/

## Health checks

- `GET /healthcheck` - liveness (process is serving)
- `GET /healthcheck/ready` - readiness (database answers a ping)
- `kfamily healthcheck` - the binary checks its own readiness (used by the
  container `HEALTHCHECK`, since the image has no curl)

## Package management

```bash
go get <package>          # install
go get <package>@none     # remove
```
