//go:build dev

package webmail

import "embed"

// UI is empty in development mode since we proxy to the Vite dev server
var UI embed.FS

// DevMode indicates whether the server is running in development mode
const DevMode = true
