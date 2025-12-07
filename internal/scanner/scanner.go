package scanner

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// ASCII markers (portable on all consoles).
const folderIcon = "[DIR]"
const fileIcon = "[FILE]"

func ScanTree(root string, skipFolders, exts []string) (string, error) {
	var b strings.Builder
	abs, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}

	err = filepath.WalkDir(abs, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			fmt.Fprintf(&b, "[error] %s: %v\n", path, walkErr)
			return nil
		}
		if d.IsDir() {
			for _, skip := range skipFolders {
				if strings.EqualFold(d.Name(), skip) {
					return filepath.SkipDir
				}
			}
			return nil
		}

		if len(exts) > 0 {
			name := strings.ToLower(d.Name())
			match := false
			for _, e := range exts {
				e = strings.ToLower(strings.TrimSpace(e))
				if e == "" {
					continue
				}
				if name == e || strings.HasSuffix(name, e) {
					match = true
					break
				}
			}
			if !match {
				return nil
			}
		}

		data, err := os.ReadFile(path)
		if err != nil {
			fmt.Fprintf(&b, "[could not read: %v]\n", err)
			return nil
		}
		fmt.Fprintf(&b, "=== File: %s ===\n%s\n\n", path, string(data))
		return nil
	})

	return b.String(), err
}

func WriteTree(root string) (string, error) {
	var b strings.Builder
	abs, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}

	fmt.Fprintf(&b, "FULL PROJECT STRUCTURE\n%s\\\n", filepath.Base(abs))

	err = filepath.WalkDir(abs, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			fmt.Fprintf(&b, "[error] %s: %v\n", path, walkErr)
			return nil
		}
		rel, _ := filepath.Rel(abs, path)
		depth := strings.Count(rel, string(os.PathSeparator))
		indent := strings.Repeat("|   ", depth)
		prefix := "|-- "
		fmt.Fprintf(&b, "%s%s%s\n", indent, prefix, filepath.Base(path))
		return nil
	})

	b.WriteString("---- END OF TREE ----\n")
	return b.String(), err
}
