//go:build !dev
// +build !dev

package unified

import "embed"

//go:embed all:.output/public
var UI embed.FS

// DevMode indicates whether the server is running in development mode
const DevMode = false
