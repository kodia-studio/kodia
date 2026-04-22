package static

import "embed"

// DistFS is the embedded file system containing the frontend build.
//go:embed all:dist/*
var DistFS embed.FS
