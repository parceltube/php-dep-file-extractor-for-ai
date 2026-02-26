package parser

import (
	"php-dep-extractor/internal/scanner"
	"strings"
)

// DependencyResult holds analysis results for selected files.
type DependencyResult struct {
	Dependencies []Dependency  `json:"dependencies"`
	Includes     []IncludeItem `json:"includes"`
}

// Dependency represents a resolved class dependency.
type Dependency struct {
	ClassName  string `json:"className"`
	FilePath   string `json:"filePath"`
	RefType    string `json:"refType"`
	ReferencedBy string `json:"referencedBy"` // which selected file references this
}

// IncludeItem represents a found include/require reference.
type IncludeItem struct {
	Type         string `json:"type"`
	RawPath      string `json:"rawPath"`
	Resolved     string `json:"resolved"`
	Line         int    `json:"line"`
	SourceFile   string `json:"sourceFile"`
}

// Resolve takes selected files and finds all their class dependencies.
func Resolve(selectedFiles []string, index *scanner.ClassIndex, projectRoot string, parseIncludes bool) (*DependencyResult, error) {
	result := &DependencyResult{}
	seenDeps := make(map[string]bool)

	for _, relPath := range selectedFiles {
		absPath := projectRoot + "/" + relPath

		// Extract class references
		refs, err := ExtractClassRefs(absPath)
		if err != nil {
			continue
		}

		for _, ref := range refs {
			className := ref.ClassName

			// Try direct lookup
			if depPath, ok := index.ClassToFile[className]; ok {
				key := depPath
				if !seenDeps[key] && !isSelected(depPath, selectedFiles) {
					seenDeps[key] = true
					result.Dependencies = append(result.Dependencies, Dependency{
						ClassName:    className,
						FilePath:     depPath,
						RefType:      ref.RefType,
						ReferencedBy: relPath,
					})
				}
			}

			// For ZF1: try resolving underscore-separated names
			// If className is like "CarrierCust" try common prefixes
			if _, ok := index.ClassToFile[className]; !ok {
				for _, prefix := range []string{"Model_", "DbTable_", "Service_", "Parent_", "Form_"} {
					fullName := prefix + className
					if depPath, ok := index.ClassToFile[fullName]; ok {
						key := depPath
						if !seenDeps[key] && !isSelected(depPath, selectedFiles) {
							seenDeps[key] = true
							result.Dependencies = append(result.Dependencies, Dependency{
								ClassName:    fullName,
								FilePath:     depPath,
								RefType:      ref.RefType,
								ReferencedBy: relPath,
							})
						}
					}
				}
			}

			// For Laravel: resolve short class name via use statements
			if strings.Contains(className, "\\") {
				if depPath, ok := index.ClassToFile[className]; ok {
					key := depPath
					if !seenDeps[key] && !isSelected(depPath, selectedFiles) {
						seenDeps[key] = true
						result.Dependencies = append(result.Dependencies, Dependency{
							ClassName:    className,
							FilePath:     depPath,
							RefType:      ref.RefType,
							ReferencedBy: relPath,
						})
					}
				}
			}
		}

		// Extract includes if enabled
		if parseIncludes {
			includes, err := ExtractIncludes(absPath, projectRoot)
			if err != nil {
				continue
			}
			for _, inc := range includes {
				result.Includes = append(result.Includes, IncludeItem{
					Type:       inc.Type,
					RawPath:    inc.RawPath,
					Resolved:   inc.Resolved,
					Line:       inc.Line,
					SourceFile: relPath,
				})
			}
		}
	}

	return result, nil
}

func isSelected(path string, selected []string) bool {
	for _, s := range selected {
		if s == path {
			return true
		}
	}
	return false
}
