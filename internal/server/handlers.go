package server

import (
	"encoding/json"
	"net/http"
	"os/exec"
	"path/filepath"
	"strings"

	"php-dep-extractor/internal/copier"
	"php-dep-extractor/internal/filetree"
	"php-dep-extractor/internal/parser"
	"php-dep-extractor/internal/scanner"
)

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

// handleBrowse opens a Windows folder picker dialog.
func handleBrowse(state *AppState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeError(w, 405, "Method not allowed")
			return
		}

		cmd := exec.Command("powershell", "-Command",
			`Add-Type -AssemblyName System.Windows.Forms; `+
				`$f = New-Object System.Windows.Forms.FolderBrowserDialog; `+
				`$f.Description = 'Select folder'; `+
				`$f.ShowNewFolderButton = $true; `+
				`if ($f.ShowDialog() -eq 'OK') { $f.SelectedPath }`)

		output, err := cmd.Output()
		if err != nil {
			writeJSON(w, map[string]string{"path": ""})
			return
		}

		path := strings.TrimSpace(string(output))
		// Convert backslashes to forward slashes
		path = filepath.ToSlash(path)
		writeJSON(w, map[string]string{"path": path})
	}
}

// handleScan scans a project directory and returns the file tree.
func handleScan(state *AppState) http.HandlerFunc {
	type scanRequest struct {
		Path      string                  `json:"path"`
		Framework string                  `json:"framework"`
		Mappings  []scanner.PrefixMapping `json:"mappings,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeError(w, 405, "Method not allowed")
			return
		}

		var req scanRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, 400, "Invalid JSON")
			return
		}

		if req.Path == "" {
			writeError(w, 400, "Path is required")
			return
		}

		// Convert forward slashes back for OS operations
		osPath := filepath.FromSlash(req.Path)

		result, err := scanner.Scan(osPath)
		if err != nil {
			writeError(w, 500, "Scan failed: "+err.Error())
			return
		}

		// Set framework
		fw := scanner.Framework(req.Framework)
		if fw == "" {
			fw = scanner.FrameworkZF1
		}

		mappings := req.Mappings
		if len(mappings) == 0 {
			mappings = scanner.DefaultZF1Mappings()
		}

		// Build class index
		index := scanner.BuildIndex(result, fw, mappings)

		// Update state
		state.ProjectRoot = result.Root
		state.ScanResult = result
		state.ClassIndex = index
		state.Framework = fw
		state.Mappings = mappings

		// Build file tree
		tree := filetree.Build(result.Files)

		writeJSON(w, map[string]any{
			"tree":      tree,
			"fileCount": len(result.Files),
			"indexed":   len(index.ClassToFile),
		})
	}
}

// handleAnalyze analyzes selected files for dependencies.
func handleAnalyze(state *AppState) http.HandlerFunc {
	type analyzeRequest struct {
		Files         []string `json:"files"`
		ParseIncludes bool     `json:"parseIncludes"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeError(w, 405, "Method not allowed")
			return
		}

		var req analyzeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, 400, "Invalid JSON")
			return
		}

		if state.ClassIndex == nil {
			writeError(w, 400, "Project not scanned yet")
			return
		}

		if len(req.Files) == 0 {
			writeError(w, 400, "No files selected")
			return
		}

		result, err := parser.Resolve(req.Files, state.ClassIndex, filepath.ToSlash(state.ProjectRoot), req.ParseIncludes)
		if err != nil {
			writeError(w, 500, "Analysis failed: "+err.Error())
			return
		}

		writeJSON(w, result)
	}
}

// handleCopy copies files to the output directory.
func handleCopy(state *AppState) http.HandlerFunc {
	type copyRequest struct {
		Files     []string `json:"files"`
		OutputDir string   `json:"outputDir"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeError(w, 405, "Method not allowed")
			return
		}

		var req copyRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, 400, "Invalid JSON")
			return
		}

		if state.ProjectRoot == "" {
			writeError(w, 400, "Project not scanned yet")
			return
		}

		if req.OutputDir == "" {
			writeError(w, 400, "Output directory is required")
			return
		}

		if len(req.Files) == 0 {
			writeError(w, 400, "No files to copy")
			return
		}

		osOutput := filepath.FromSlash(req.OutputDir)
		result := copier.CopyFiles(req.Files, state.ProjectRoot, osOutput)

		writeJSON(w, result)
	}
}

// handleSettings returns/updates the current prefix mappings.
func handleSettings(state *AppState) http.HandlerFunc {
	type settingsRequest struct {
		Mappings []scanner.PrefixMapping `json:"mappings"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			writeJSON(w, map[string]any{
				"framework": state.Framework,
				"mappings":  state.Mappings,
			})
		case http.MethodPost:
			var req settingsRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				writeError(w, 400, "Invalid JSON")
				return
			}
			state.Mappings = req.Mappings
			writeJSON(w, map[string]string{"status": "ok"})
		default:
			writeError(w, 405, "Method not allowed")
		}
	}
}
