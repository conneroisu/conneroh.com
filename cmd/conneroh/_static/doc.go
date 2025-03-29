// Package static contains static files for the web application.
package static

import (
	"embed"
	_ "embed"
)

//go:embed dist/**
var Dist embed.FS

//go:embed dist/favicon.ico
var Favicon []byte
