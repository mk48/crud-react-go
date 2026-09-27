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
- `statsviz` for live runtime stats

## Getting started

```bash
go run ./cmd/api
```

The server listens on `:8080` by default. On startup it connects to
Postgres and runs any pending migrations.

Copy `.env-sample` to `.env` and fill in real values before running:

| Variable | Purpose |
|---|---|
| `ENV` | `dev` / `prod` |
| `CONNECTION_STRING` | Postgres connection string |
| `CASDOOR_ENDPOINT` | Casdoor server URL |
| `CASDOOR_CLIENT_ID` / `CASDOOR_CLIENT_SECRET` | Casdoor application credentials |
| `CASDOOR_ORGANIZATION_NAME` / `CASDOOR_APPLICATION_NAME` | Casdoor organization/application names |
| `CASDOOR_CERTIFICATE` | PEM certificate (Casdoor's Cert management page) used to verify token signatures |

The frontend (`../web`) needs its own `VITE_CASDOOR_*` values pointing at
the same Casdoor application - see its README.

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

## Runtime stats

http://localhost:8080/kfamily-debug/statsviz

## Package management

```bash
go get <package>          # install
go get <package>@none     # remove
```
