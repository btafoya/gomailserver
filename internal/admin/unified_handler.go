package admin

import (
	"io"
	"io/fs"
	"net/http"
	"path"
	"strings"

	"github.com/btafoya/gomailserver/internal/config"
	unified "github.com/btafoya/gomailserver/unified-go"
	"go.uber.org/zap"
)

// UnifiedHandler returns an HTTP handler for the unified UI
// In development mode (with -tags dev), it returns 404 for UI routes
// In production mode, it serves embedded static files with SPA fallback
func UnifiedHandler(logger *zap.Logger, webUIConfig *config.WebUIConfig) http.Handler {
	if unified.DevMode {
		// Development mode: Nuxt dev server handles all UI routes
		// Return 404 for UI routes since they're handled by Nuxt
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Debug("Development mode: rejecting UI route on API server",
				zap.String("path", r.URL.Path))
			http.Error(w, "UI routes not available on API server. Use Nuxt dev server on port configured WebUI port.", http.StatusNotFound)
		})
	}
	return unifiedProdModeHandler(logger)
}

// serveFile attempts to serve a file from embedded filesystem
// Returns true if successful, false if file not found
func serveFile(fsys fs.FS, w http.ResponseWriter, r *http.Request, name string) bool {
	name = path.Clean(name)
	if name == "/" {
		name = "index.html"
	} else {
		name = strings.TrimPrefix(name, "/")
	}

	f, err := fsys.Open(name)
	if err != nil {
		return false
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return false
	}

	if stat.IsDir() {
		indexPath := path.Join(name, "index.html")
		indexFile, err := fsys.Open(indexPath)
		if err != nil {
			return false
		}
		defer indexFile.Close()

		indexStat, err := indexFile.Stat()
		if err != nil {
			return false
		}

		http.ServeContent(w, r, indexPath, indexStat.ModTime(), indexFile.(io.ReadSeeker))
		return true
	}

	http.ServeContent(w, r, name, stat.ModTime(), f.(io.ReadSeeker))
	return true
}

// unifiedProdModeHandler serves embedded static files with SPA fallback
func unifiedProdModeHandler(logger *zap.Logger) http.Handler {
	// Get the embedded filesystem
	distFS, err := fs.Sub(unified.UI, ".output/public")
	if err != nil {
		logger.Error("Failed to create unified dist filesystem", zap.Error(err))
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Unified UI not available", http.StatusServiceUnavailable)
		})
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Strip any of the supported prefixes (/admin, /portal, /webmail)
		path := r.URL.Path
		path = strings.TrimPrefix(path, "/admin")
		path = strings.TrimPrefix(path, "/portal")
		path = strings.TrimPrefix(path, "/webmail")
		if path == "" {
			path = "/"
		}

		// Try to serve the requested file
		if serveFile(distFS, w, r, path) {
			return
		}

		// File not found - serve index.html for SPA routing
		// This allows Nuxt Router to handle the route client-side
		logger.Debug("Serving index.html for SPA route", zap.String("path", r.URL.Path))

		indexFile, err := distFS.Open("index.html")
		if err != nil {
			logger.Error("Failed to open index.html", zap.Error(err))
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		defer indexFile.Close()

		stat, err := indexFile.Stat()
		if err != nil {
			logger.Error("Failed to stat index.html", zap.Error(err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		http.ServeContent(w, r, "index.html", stat.ModTime(), indexFile.(io.ReadSeeker))
	})
}
