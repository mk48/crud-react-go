//go:build !embedui

package webui

import "io/fs"

// Without -tags embedui there's no embedded web app: in dev the Vite dev
// server serves it (proxying /api and /config.js to this server).
func files() fs.FS {
	return nil
}
