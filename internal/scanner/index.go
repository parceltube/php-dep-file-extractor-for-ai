package scanner

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Framework represents the PHP framework type.
type Framework string

const (
	FrameworkZF1     Framework = "zf1"
	FrameworkCakePHP Framework = "cakephp"
	FrameworkLaravel Framework = "laravel"
)

// PrefixMapping maps a class prefix to a directory (for ZF1 style).
type PrefixMapping struct {
	Prefix string `json:"prefix"`
	Dir    string `json:"dir"`
}

// ClassIndex maps class names to their relative file paths.
type ClassIndex struct {
	ClassToFile map[string]string // className -> relative path
	FileToClass map[string]string // relative path -> className
}

// DefaultZF1Mappings returns the default ZF1 prefix→directory mappings.
func DefaultZF1Mappings() []PrefixMapping {
	return []PrefixMapping{
		{Prefix: "Parent_", Dir: "parents/"},
		{Prefix: "DbTable_", Dir: "dbs/"},
		{Prefix: "Service_", Dir: "services/"},
		{Prefix: "Model_", Dir: "models/"},
		{Prefix: "Form_", Dir: "forms/"},
	}
}

// BuildIndex creates a class name index from scanned files.
func BuildIndex(result *ScanResult, fw Framework, mappings []PrefixMapping) *ClassIndex {
	idx := &ClassIndex{
		ClassToFile: make(map[string]string),
		FileToClass: make(map[string]string),
	}

	for _, relPath := range result.Files {
		var className string
		switch fw {
		case FrameworkZF1:
			className = zf1ClassFromPath(relPath, mappings)
		case FrameworkCakePHP:
			className = cakeClassFromPath(relPath)
		case FrameworkLaravel:
			className = laravelClassFromPath(relPath)
		}

		if className == "" {
			// Fallback: read file header to find class declaration
			className = classFromFileContent(filepath.Join(result.Root, filepath.FromSlash(relPath)))
		}

		if className != "" {
			idx.ClassToFile[className] = relPath
			idx.FileToClass[relPath] = className
		}
	}

	return idx
}

// zf1ClassFromPath derives class name from ZF1 path conventions.
// e.g. "application/models/Car/CarrierCust.php" -> "Model_Car_CarrierCust"
func zf1ClassFromPath(relPath string, mappings []PrefixMapping) string {
	// Normalize path
	p := strings.TrimSuffix(relPath, ".php")

	// Try to strip "application/" prefix
	appPath := p
	if strings.HasPrefix(p, "application/") {
		appPath = strings.TrimPrefix(p, "application/")
	} else {
		// Not under application/, skip ZF1 mapping
		return ""
	}

	// Check each mapping (order matters - longer/more specific prefixes first)
	for _, m := range mappings {
		dir := strings.TrimSuffix(m.Dir, "/")
		if strings.HasPrefix(appPath, dir+"/") {
			rest := strings.TrimPrefix(appPath, dir+"/")
			// Convert path separators to underscores
			className := m.Prefix + strings.ReplaceAll(rest, "/", "_")
			return className
		}
	}

	// Controllers: application/controllers/V3/CustomersController.php
	if strings.HasPrefix(appPath, "controllers/") {
		return ""
	}

	return ""
}

// cakeClassFromPath derives class name from CakePHP path conventions.
func cakeClassFromPath(relPath string) string {
	p := strings.TrimSuffix(relPath, ".php")

	// CakePHP 2.x: app/Model/Post.php -> Post
	prefixes := []string{"app/", "src/"}
	for _, prefix := range prefixes {
		if strings.HasPrefix(p, prefix) {
			rest := strings.TrimPrefix(p, prefix)
			// The class name is the filename (last component)
			parts := strings.Split(rest, "/")
			if len(parts) > 0 {
				return parts[len(parts)-1]
			}
		}
	}
	return ""
}

// laravelClassFromPath derives fully-qualified class name from Laravel PSR-4 conventions.
func laravelClassFromPath(relPath string) string {
	p := strings.TrimSuffix(relPath, ".php")

	// app/Models/User.php -> App\Models\User
	if strings.HasPrefix(p, "app/") {
		rest := strings.TrimPrefix(p, "app/")
		return "App\\" + strings.ReplaceAll(rest, "/", "\\")
	}
	return ""
}

var classRegex = regexp.MustCompile(`(?m)^\s*(?:abstract\s+|final\s+)?class\s+(\w+)`)
var interfaceRegex = regexp.MustCompile(`(?m)^\s*interface\s+(\w+)`)
var traitRegex = regexp.MustCompile(`(?m)^\s*trait\s+(\w+)`)

// classFromFileContent reads the first part of a PHP file to find class declaration.
func classFromFileContent(absPath string) string {
	f, err := os.Open(absPath)
	if err != nil {
		return ""
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
		if lineCount > 100 {
			break
		}
		line := scanner.Text()

		if m := classRegex.FindStringSubmatch(line); len(m) > 1 {
			return m[1]
		}
		if m := interfaceRegex.FindStringSubmatch(line); len(m) > 1 {
			return m[1]
		}
		if m := traitRegex.FindStringSubmatch(line); len(m) > 1 {
			return m[1]
		}
	}
	return ""
}
