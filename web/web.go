// Package web embeds the built Angular SPA. Until the frontend skeleton
// subissue (#9) lands, this serves a single placeholder index.html so the
// static-asset mount has something to return.
package web

import (
	"embed"
	"io/fs"
)

//go:embed dist
var distFS embed.FS

// FS returns the static asset filesystem rooted at the SPA build output.
func FS() fs.FS {
	sub, err := fs.Sub(distFS, "dist")
	if err != nil {
		// Unreachable: dist/ is embedded at compile time.
		panic(err)
	}
	return sub
}
