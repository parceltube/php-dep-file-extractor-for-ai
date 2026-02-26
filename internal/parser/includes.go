package parser

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// IncludeRef represents a require/include statement found in a PHP file.
type IncludeRef struct {
	Type     string `json:"type"`     // "require", "require_once", "include", "include_once"
	RawPath  string `json:"rawPath"`  // original path expression
	Resolved string `json:"resolved"` // resolved relative path (if possible)
	Line     int    `json:"line"`
}

var (
	reInclude = regexp.MustCompile(`(?:require_once|include_once|require|include)\s*[\(]?\s*(.+?)\s*[\)]?\s*;`)
	reIncType = regexp.MustCompile(`(require_once|include_once|require|include)`)
)

// ExtractIncludes extracts require/include statements from a PHP file.
func ExtractIncludes(filePath string, projectRoot string) ([]IncludeRef, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	var refs []IncludeRef

	fileDir := filepath.Dir(filePath)

	for i, line := range lines {
		lineNum := i + 1
		trimmed := strings.TrimSpace(line)

		// Skip comments
		if strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "*") || strings.HasPrefix(trimmed, "/*") {
			continue
		}

		m := reInclude.FindStringSubmatch(line)
		if len(m) < 2 {
			continue
		}

		typeMatch := reIncType.FindString(line)
		rawPath := strings.TrimSpace(m[1])

		resolved := resolveIncludePath(rawPath, fileDir, projectRoot)

		refs = append(refs, IncludeRef{
			Type:     typeMatch,
			RawPath:  rawPath,
			Resolved: resolved,
			Line:     lineNum,
		})
	}

	return refs, nil
}

// resolveIncludePath attempts to resolve a PHP include path expression to a relative path.
func resolveIncludePath(rawPath string, fileDir string, projectRoot string) string {
	// Remove quotes for simple string paths
	path := strings.Trim(rawPath, "'\"")

	// Handle APPLICATION_PATH . '/...'
	if strings.Contains(rawPath, "APPLICATION_PATH") {
		re := regexp.MustCompile(`APPLICATION_PATH\s*\.\s*['"]([^'"]+)['"]`)
		if m := re.FindStringSubmatch(rawPath); len(m) > 1 {
			// APPLICATION_PATH typically points to application/ dir
			p := strings.TrimPrefix(m[1], "/")
			resolved := "application/" + p
			return filepath.ToSlash(resolved)
		}
	}

	// Handle dirname(__FILE__) . '/...'
	if strings.Contains(rawPath, "dirname") || strings.Contains(rawPath, "__DIR__") {
		re := regexp.MustCompile(`(?:dirname\s*\(\s*__FILE__\s*\)|__DIR__)\s*\.\s*['"]([^'"]+)['"]`)
		if m := re.FindStringSubmatch(rawPath); len(m) > 1 {
			absPath := filepath.Join(fileDir, m[1])
			if rel, err := filepath.Rel(projectRoot, absPath); err == nil {
				return filepath.ToSlash(rel)
			}
		}
	}

	// Simple string path (no variables/concatenation)
	if !strings.ContainsAny(path, "$.(") && (strings.HasSuffix(path, ".php") || strings.Contains(path, "/")) {
		// Try relative to file directory
		absPath := filepath.Join(fileDir, path)
		if _, err := os.Stat(absPath); err == nil {
			if rel, err := filepath.Rel(projectRoot, absPath); err == nil {
				return filepath.ToSlash(rel)
			}
		}
		// Try relative to project root
		absPath = filepath.Join(projectRoot, path)
		if _, err := os.Stat(absPath); err == nil {
			if rel, err := filepath.Rel(projectRoot, absPath); err == nil {
				return filepath.ToSlash(rel)
			}
		}
	}

	return ""
}
