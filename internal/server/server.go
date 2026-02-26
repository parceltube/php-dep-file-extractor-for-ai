package server

import (
	"embed"
	"io/fs"
	"net/http"

	"php-dep-extractor/internal/scanner"
)

// AppState holds the current application state.
type AppState struct {
	ProjectRoot string
	ScanResult  *scanner.ScanResult
	ClassIndex  *scanner.ClassIndex
	Framework   scanner.Framework
	Mappings    []scanner.PrefixMapping
}

// New creates a new HTTP handler with all routes registered.
func New(webFS embed.FS) http.Handler {
	state := &AppState{
		Framework: scanner.FrameworkZF1,
		Mappings:  scanner.DefaultZF1Mappings(),
	}

	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("/api/browse", handleBrowse(state))
	mux.HandleFunc("/api/scan", handleScan(state))
	mux.HandleFunc("/api/analyze", handleAnalyze(state))
	mux.HandleFunc("/api/copy", handleCopy(state))
	mux.HandleFunc("/api/settings", handleSettings(state))

	// Serve embedded web files
	webSub, _ := fs.Sub(webFS, "web")
	fileServer := http.FileServer(http.FS(webSub))
	mux.Handle("/", fileServer)

	return mux
}
