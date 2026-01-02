package webmail

import (
	"io"
	"io/fs"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
	"strings"

	webmail "github.com/btafoya/gomailserver/web/webmail"
	"go.uber.org/zap"
)

// Handler returns an HTTP handler for the webmail UI
// In development mode, it proxies to Nuxt dev server
// In production mode, it serves embedded static files with SPA fallback
func Handler(logger *zap.Logger) http.Handler {
	if webmail.DevMode {
		return devModeHandler(logger)
	}
	return prodModeHandler(logger)
}

// devModeHandler proxies requests to Nuxt dev server if running
// Falls back to embedded assets if Nuxt is not available
func devModeHandler(logger *zap.Logger) http.Handler {
	nuxtURL := "http://localhost:3000"

	// Try to connect to Nuxt dev server
	if isNuxtRunning(nuxtURL) {
		logger.Info("Webmail UI: Proxying to Nuxt dev server", zap.String("url", nuxtURL))
		target, _ := url.Parse(nuxtURL)
		proxy := httputil.NewSingleHostReverseProxy(target)

		// Modify the request to handle /webmail prefix
		originalDirector := proxy.Director
		proxy.Director = func(req *http.Request) {
			originalDirector(req)
			req.Host = target.Host
			// Keep /webmail prefix for Nuxt since it's configured with baseURL
			// No need to remove prefix
		}

		return proxy
	}

	logger.Warn("Webmail UI: Nuxt dev server not running, serving embedded assets")
	return prodModeHandler(logger)
}

// prodModeHandler serves embedded static files with SPA fallback
func prodModeHandler(logger *zap.Logger) http.Handler {
	// Get the embedded filesystem (.output/public from Nuxt build)
	publicFS, err := fs.Sub(webmail.UI, ".output/public")
	if err != nil {
		logger.Error("Failed to create public filesystem", zap.Error(err))
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Webmail UI not available", http.StatusServiceUnavailable)
		})
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Remove /webmail prefix
		urlPath := strings.TrimPrefix(r.URL.Path, "/webmail")
		if urlPath == "" {
			urlPath = "/"
		}

		// Try to serve the requested file
		if serveFile(publicFS, w, r, urlPath) {
			return
		}

		// File not found - serve index.html for SPA routing
		// This allows Nuxt Router to handle the route client-side
		logger.Debug("Serving index.html for SPA route", zap.String("path", r.URL.Path))

		indexFile, err := publicFS.Open("index.html")
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

// isNuxtRunning checks if Nuxt dev server is running
func isNuxtRunning(nuxtURL string) bool {
	resp, err := http.Get(nuxtURL)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}
