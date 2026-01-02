//go:build !dev

package webmail

import "embed"

// UI contains the compiled webmail UI static files
// This is populated at build time with the contents of dist/
//
//go:embed all:dist
var UI embed.FS

// DevMode indicates whether the server is running in development mode
const DevMode = false
