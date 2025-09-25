package scanner

import (
    "fmt"
    "io/fs"
    "os"
    "path/filepath"
    "strings"

    "github.com/nimble-sloth/go-smart-folder-scanner/internal/preview"
)

// ASCII markers (portable on all consoles).
const folderIcon = "[DIR]"
const fileIcon = "[FILE]"

func ScanTree(root string, maxPreviewBytes int) (string, error) {
    var b strings.Builder
    abs, err := filepath.Abs(root)
    if err != nil { return "", err }

    fmt.Fprintf(&b, "%s Root: %s\n\n", folderIcon, abs)
    b.WriteString(folderIcon + " " + filepath.Base(abs) + "\n")

    err = filepath.WalkDir(abs, func(path string, d fs.DirEntry, walkErr error) error {
        if walkErr != nil {
            fmt.Fprintf(&b, "[error] %s: %v\n", path, walkErr)
            return nil
        }
        if path == abs { return nil }

        rel, _ := filepath.Rel(abs, path)
        depth := 0
        if rel != "." { depth = strings.Count(rel, string(os.PathSeparator)) }
        indent := strings.Repeat("  ", depth)
        name := filepath.Base(path)

        if d.IsDir() {
            fmt.Fprintf(&b, "%s%s %s\n", indent, folderIcon, name)
            return nil
        }

        fmt.Fprintf(&b, "%s%s %s  (at %s)\n", indent, fileIcon, name, path)

        pv, err := preview.PreviewFile(path, maxPreviewBytes)
        if err != nil {
            fmt.Fprintf(&b, "%s  [could not read: %v]\n", indent, err)
            return nil
        }
        if pv == "" {
            fmt.Fprintf(&b, "%s  [no preview]\n", indent)
            return nil
        }
        b.WriteString(indent + "  --- content preview ---\n")
        for _, line := range strings.Split(pv, "\n") {
            b.WriteString(indent + "  " + line + "\n")
        }
        b.WriteString(indent + "  --- end ---\n")
        return nil
    })
    return b.String(), err
}
