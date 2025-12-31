//go:build dev

package admin

import "embed"

// UI is not used in dev mode (we proxy to Vite instead)
var UI embed.FS

// DevMode indicates whether the server is running in development mode
const DevMode = true
