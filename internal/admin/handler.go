package admin

import (
	"io"
	"io/fs"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
	"strings"

	webadmin "github.com/btafoya/gomailserver/web/admin"
	"go.uber.org/zap"
)

// Handler returns an HTTP handler for the admin UI
// In development mode (with -tags dev), it proxies to Vite dev server
// In production mode, it serves embedded static files with SPA fallback
func Handler(logger *zap.Logger) http.Handler {
	if webadmin.DevMode {
		return devModeHandler(logger)
	}
	return prodModeHandler(logger)
}

// devModeHandler proxies requests to Vite dev server if running
// Falls back to embedded assets if Vite is not available
func devModeHandler(logger *zap.Logger) http.Handler {
	viteURL := "http://localhost:5173"

	// Try to connect to Vite dev server
	if isViteRunning(viteURL) {
		logger.Info("Admin UI: Proxying to Vite dev server", zap.String("url", viteURL))
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

	logger.Warn("Admin UI: Vite dev server not running, serving embedded assets")
	return prodModeHandler(logger)
}

// prodModeHandler serves embedded static files with SPA fallback
func prodModeHandler(logger *zap.Logger) http.Handler {
	// Get the embedded filesystem
	distFS, err := fs.Sub(webadmin.UI, "dist")
	if err != nil {
		logger.Error("Failed to create dist filesystem", zap.Error(err))
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Admin UI not available", http.StatusServiceUnavailable)
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

