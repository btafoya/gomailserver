package admin

import (
	"io"
	"io/fs"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	webunified "github.com/btafoya/gomailserver/web/unified-go"
	"go.uber.org/zap"
)

// UnifiedHandler returns an HTTP handler for the unified UI
// In development mode (with -tags dev), it proxies to Vite dev server
// In production mode, it serves embedded static files with SPA fallback
func UnifiedHandler(logger *zap.Logger) http.Handler {
	if webunified.DevMode {
		return unifiedDevModeHandler(logger)
	}
	return unifiedProdModeHandler(logger)
}

// unifiedDevModeHandler proxies requests to Vite dev server if running
// Falls back to embedded assets if Vite is not available
func unifiedDevModeHandler(logger *zap.Logger) http.Handler {
	viteURL := "http://localhost:5173"

	// Try to connect to Vite dev server
	if isViteRunning(viteURL) {
		logger.Info("Unified UI: Proxying to Vite dev server", zap.String("url", viteURL))
		target, _ := url.Parse(viteURL)
		proxy := httputil.NewSingleHostReverseProxy(target)

		// Modify the request to handle /admin prefix
		originalDirector := proxy.Director
		proxy.Director = func(req *http.Request) {
			originalDirector(req)
			req.Host = target.Host
			// Remove /admin prefix for Vite
			req.URL.Path = strings.TrimPrefix(req.URL.Path, "/admin")
			if req.URL.Path == "" {
				req.URL.Path = "/"
			}
		}

		return http.StripPrefix("/admin", proxy)
	}

	logger.Warn("Unified UI: Vite dev server not running, serving embedded assets")
	return unifiedProdModeHandler(logger)
}

// unifiedProdModeHandler serves embedded static files with SPA fallback
func unifiedProdModeHandler(logger *zap.Logger) http.Handler {
	// Get the embedded filesystem
	distFS, err := fs.Sub(webunified.UI, "dist")
	if err != nil {
		logger.Error("Failed to create unified dist filesystem", zap.Error(err))
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Unified UI not available", http.StatusServiceUnavailable)
		})
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Remove /admin prefix
		path := strings.TrimPrefix(r.URL.Path, "/admin")
		if path == "" {
			path = "/"
		}

		// Try to serve the requested file
		if serveFile(distFS, w, r, path) {
			return
		}

		// File not found - serve index.html for SPA routing
		// This allows Vue Router to handle the route client-side
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
