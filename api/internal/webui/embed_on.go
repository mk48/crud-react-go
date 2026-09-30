//go:build embedui

package webui

import (
	"embed"
	"io/fs"
)

// dist is the built web app - `pnpm build` in ../web writes it here (see
// web/vite.config.ts). Only compiled in with `-tags embedui` (the Docker
// build), so a plain `go run` doesn't need a web build first.
//
//go:embed all:dist
var dist embed.FS

func files() fs.FS {
	sub, err := fs.Sub(dist, "dist")
	if err != nil {
		panic(err) // "dist" is a literal embedded above - can't be missing
	}
	return sub
}
