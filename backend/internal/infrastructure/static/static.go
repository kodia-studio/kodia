package static

import "embed"

// DistFS is the embedded file system containing the frontend build.
// For production, this contains the built SvelteKit frontend.
// For development, this may be empty or contain placeholder files.
//
//go:embed all:dist/*
var DistFS embed.FS

// HasFrontend checks if the frontend build exists.
// Returns true if dist/ contains actual frontend files (not just .gitkeep).
func HasFrontend() bool {
	entries, err := DistFS.ReadDir("dist")
	if err != nil {
		return false
	}

	for _, entry := range entries {
		// If we find anything other than .gitkeep, frontend exists
		if entry.Name() != ".gitkeep" {
			return true
		}
	}
	return false
}
