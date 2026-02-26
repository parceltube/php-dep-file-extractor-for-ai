package scanner

import (
	"os"
	"path/filepath"
	"strings"
)

// ScanResult holds all PHP files found in a project directory.
type ScanResult struct {
	Files []string // relative paths using forward slashes
	Root  string   // absolute project root
}

// ExcludeDirs are directories to skip during scanning.
var ExcludeDirs = map[string]bool{
	"vendor":       true,
	"node_modules": true,
	".git":         true,
	".svn":         true,
	".idea":        true,
}

// Scan walks the project directory and collects all .php file paths.
func Scan(root string) (*ScanResult, error) {
	root, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}

	var files []string
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			name := info.Name()
			if ExcludeDirs[name] {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasSuffix(strings.ToLower(info.Name()), ".php") {
			rel, _ := filepath.Rel(root, path)
			files = append(files, filepath.ToSlash(rel))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &ScanResult{Files: files, Root: root}, nil
}
