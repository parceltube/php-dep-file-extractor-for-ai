package filetree

import (
	"sort"
	"strings"
)

// TreeNode represents a node in the file tree.
type TreeNode struct {
	Name     string      `json:"name"`
	Path     string      `json:"path,omitempty"` // relative path for files
	IsDir    bool        `json:"isDir"`
	Children []*TreeNode `json:"children,omitempty"`
}

// Build creates a tree structure from a flat list of relative file paths.
func Build(files []string) *TreeNode {
	root := &TreeNode{Name: "/", IsDir: true}

	for _, filePath := range files {
		parts := strings.Split(filePath, "/")
		current := root

		for i, part := range parts {
			isLast := i == len(parts)-1

			// Find existing child
			var child *TreeNode
			for _, c := range current.Children {
				if c.Name == part {
					child = c
					break
				}
			}

			if child == nil {
				child = &TreeNode{
					Name:  part,
					IsDir: !isLast,
				}
				if isLast {
					child.Path = filePath
				}
				current.Children = append(current.Children, child)
			}

			current = child
		}
	}

	// Sort children: directories first, then alphabetically
	sortTree(root)
	return root
}

func sortTree(node *TreeNode) {
	if node.Children == nil {
		return
	}

	sort.Slice(node.Children, func(i, j int) bool {
		a, b := node.Children[i], node.Children[j]
		if a.IsDir != b.IsDir {
			return a.IsDir // dirs first
		}
		return strings.ToLower(a.Name) < strings.ToLower(b.Name)
	})

	for _, child := range node.Children {
		sortTree(child)
	}
}
