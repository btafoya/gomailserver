//go:build dev
// +build dev

package unified

import "embed"

// UI is empty in development mode since we proxy to the Nuxt dev server
var UI embed.FS

// DevMode indicates whether the server is running in development mode
const DevMode = true
