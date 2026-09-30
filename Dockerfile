# syntax=docker/dockerfile:1

# One image, one binary: the Go API with the built React app embedded
# (-tags embedui). Build from the repo root:
#   docker build --build-arg VERSION=$(git describe --tags --always) -t kfamily .

### Web: build the React app into api/internal/webui/dist (see web/vite.config.ts)
FROM node:24-alpine AS web
RUN npm install --global pnpm@11
WORKDIR /src/web
COPY web/package.json web/pnpm-lock.yaml web/pnpm-workspace.yaml ./
RUN pnpm install --frozen-lockfile
COPY web/ ./
RUN pnpm build

### API: compile a static binary with the web app embedded
FROM golang:1.25-alpine AS api
WORKDIR /src/api
COPY api/go.mod api/go.sum ./
RUN go mod download
COPY api/ ./
COPY --from=web /src/api/internal/webui/dist ./internal/webui/dist
ARG VERSION=dev
RUN CGO_ENABLED=0 GOOS=linux go build -tags embedui -trimpath \
    -ldflags "-s -w -X main.version=${VERSION}" \
    -o /out/kfamily ./cmd/api

### Runtime: distroless - no shell or package manager, CA certs and tzdata
### included, runs as an unprivileged user (uid 65532).
FROM gcr.io/distroless/static-debian12:nonroot
LABEL maintainer="M Kumaran"
COPY --from=api /out/kfamily /kfamily
USER nonroot:nonroot
EXPOSE 8080
# Ready = serving and the database answers (see /healthcheck/ready).
HEALTHCHECK --interval=30s --timeout=6s --start-period=30s --retries=3 \
    CMD ["/kfamily", "healthcheck"]
ENTRYPOINT ["/kfamily"]
