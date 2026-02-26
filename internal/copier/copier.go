package copier

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// CopyResult holds the results of a copy operation.
type CopyResult struct {
	Copied []string `json:"copied"`
	Errors []string `json:"errors"`
}

// CopyFiles copies the given relative paths from srcRoot to dstRoot, preserving directory structure.
func CopyFiles(files []string, srcRoot string, dstRoot string) *CopyResult {
	result := &CopyResult{}

	for _, relPath := range files {
		srcPath := filepath.Join(srcRoot, filepath.FromSlash(relPath))
		dstPath := filepath.Join(dstRoot, filepath.FromSlash(relPath))

		// Create destination directory
		dstDir := filepath.Dir(dstPath)
		if err := os.MkdirAll(dstDir, 0755); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("mkdir %s: %v", relPath, err))
			continue
		}

		// Copy file
		if err := copyFile(srcPath, dstPath); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("copy %s: %v", relPath, err))
			continue
		}

		result.Copied = append(result.Copied, relPath)
	}

	return result
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}
