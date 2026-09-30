# KFamily

Family-tree app: a Go API ([api/](api/)) and a React web app ([web/](web/)),
shipped as **one binary** — the built web app is embedded in the API
(`-tags embedui`) and served from the same origin.

## Local development

Run both dev servers (Windows Terminal: `pwsh ./dev.ps1`, or two terminals):

```bash
cd api && go run ./cmd/api     # API on :8080 (needs api/.env - see api/README.md)
cd web && pnpm dev             # web on :5173, proxies /api and /config.js to :8080
```

Open http://localhost:5173. With `ENV=dev` the API also serves Swagger UI at
http://localhost:8080/swagger/index.html and runtime stats at
http://localhost:8080/kfamily-debug/statsviz/.

## Build and deploy

```bash
docker build --build-arg VERSION=$(git describe --tags --always) -t kfamily .
```

The [Dockerfile](Dockerfile) builds the web app, embeds it into a static Go
binary and ships it on a distroless, non-root base image. Without Docker:

```bash
cd web && pnpm build                                        # -> api/internal/webui/dist
cd api && go build -tags embedui -o kfamily ./cmd/api
```

Deploying (e.g. Dokploy):

- Build from the **repository root** with the root `Dockerfile`.
- Configure it only through the API's environment variables (see
  [api/README.md](api/README.md)) — the web app reads its Casdoor settings
  from the API at runtime (`/config.js`), so one image works everywhere.
- Set `ENV` to something other than `dev` (e.g. `prod`) so Swagger and
  statsviz aren't exposed.
- Register `https://<your-domain>/callback` as a Redirect URL in the Casdoor
  application.
- Health: `GET /healthcheck` (liveness) and `GET /healthcheck/ready` (checks
  the database; also used by the image's `HEALTHCHECK`).
- Migrations run on startup under a Postgres advisory lock, so several
  replicas can start at once.
