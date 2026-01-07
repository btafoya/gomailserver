package admin

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
	"strings"

	"github.com/btafoya/gomailserver/internal/config"
	webunified "github.com/btafoya/gomailserver/unified-go"
	"go.uber.org/zap"
)

// UnifiedHandler returns an HTTP handler for the unified UI
// In development mode (with -tags dev), it proxies to Vite dev server
// In production mode, it serves embedded static files with SPA fallback
// This handler works for /admin, /portal, and /webmail routes
func UnifiedHandler(logger *zap.Logger, webUIConfig *config.WebUIConfig) http.Handler {
	if webunified.DevMode {
		return unifiedDevModeHandler(logger, webUIConfig)
	}
	return unifiedProdModeHandler(logger)
}

// unifiedDevModeHandler proxies requests to Vite dev server if running
// Falls back to embedded assets if Vite is not available
func unifiedDevModeHandler(logger *zap.Logger, config *config.WebUIConfig) http.Handler {
	viteURL := fmt.Sprintf("http://localhost:%d", config.VitePort)

	// Try to connect to Vite dev server
	if isViteRunning(viteURL) {
		logger.Info("Unified UI: Proxying to Vite dev server", zap.String("url", viteURL))
		target, _ := url.Parse(viteURL)
		proxy := httputil.NewSingleHostReverseProxy(target)

		// Modify the request to handle multiple prefixes (/admin, /portal, /webmail)
		originalDirector := proxy.Director
		proxy.Director = func(req *http.Request) {
			originalDirector(req)
			req.Host = target.Host
			path := req.URL.Path
			path = strings.TrimPrefix(path, "/admin")
			path = strings.TrimPrefix(path, "/portal")
			path = strings.TrimPrefix(path, "/webmail")
			if path == "" {
				path = "/"
			}
			req.URL.Path = path
		}

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var stripPrefix string
			if strings.HasPrefix(r.URL.Path, "/admin") {
				stripPrefix = "/admin"
			} else if strings.HasPrefix(r.URL.Path, "/portal") {
				stripPrefix = "/portal"
			} else if strings.HasPrefix(r.URL.Path, "/webmail") {
				stripPrefix = "/webmail"
			}
			http.StripPrefix(stripPrefix, proxy).ServeHTTP(w, r)
		})
	}

	logger.Warn("Unified UI: Vite dev server not running, serving embedded assets")
	return unifiedProdModeHandler(logger)
}

// serveFile attempts to serve a file from the embedded filesystem
// Returns true if successful, false if file not found
func serveFile(fsys fs.FS, w http.ResponseWriter, r *http.Request, name string) bool {
	// Clean the path
	name = path.Clean(name)
	if name == "/" {
		name = "index.html"
	} else {
		name = strings.TrimPrefix(name, "/")
	}

	// Try to open the file
	f, err := fsys.Open(name)
	if err != nil {
		return false
	}
	defer f.Close()

	// Check if it's a directory
	stat, err := f.Stat()
	if err != nil {
		return false
	}

	if stat.IsDir() {
		// Try index.html in the directory
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

	// Serve the file
	http.ServeContent(w, r, name, stat.ModTime(), f.(io.ReadSeeker))
	return true
}

// isViteRunning checks if Vite dev server is running
func isViteRunning(viteURL string) bool {
	resp, err := http.Get(viteURL)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

// unifiedProdModeHandler serves embedded static files with SPA fallback
func unifiedProdModeHandler(logger *zap.Logger) http.Handler {
	// Get the embedded filesystem
	distFS, err := fs.Sub(webunified.UI, ".output/public")
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
