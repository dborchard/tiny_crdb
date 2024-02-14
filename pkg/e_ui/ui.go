package ui

import "net/http"

// Handler returns an http.Handler that serves the UI,
// including index.html, which has some login-related variables
// templated into it, as well as static assets.
func Handler(cfg Config) http.Handler {
	// etags is used to provide a unique per-file checksum for each served file,
	// which enables client-side caching using Cache-Control and ETag headers.

}
