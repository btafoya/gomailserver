//go:build !dev
// +build !dev

package unified

import "embed"

//go:embed all:dist
var UI embed.FS

const DevMode = false
