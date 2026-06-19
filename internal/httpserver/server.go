// Package httpserver wires the HTTP routing layer. It deliberately holds no
// business logic — only the entry-point routes (/healthz) and the static
// asset mount. Feature endpoints will be registered by later subissues.
package httpserver

import (
	"encoding/json"
	"io/fs"
	"log/slog"
	"net/http"

	"github.com/zredinger-ccc/git-forge/internal/version"
)

// New returns the root HTTP handler. staticFS is the embedded SPA bundle
// (subtree rooted at index.html). Pass an empty fs.FS in tests if you do not
// need static serving.
func New(logger *slog.Logger, staticFS fs.FS) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", healthz)

	if staticFS != nil {
		mux.Handle("GET /", http.FileServerFS(staticFS))
	}

	return loggingMiddleware(logger)(mux)
}

func healthz(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"status": "ok",
		"build":  version.Info(),
	})
}
