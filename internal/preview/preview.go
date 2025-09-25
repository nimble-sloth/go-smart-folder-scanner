package preview

import (
    "bytes"
    "errors"
    "fmt"
    "io"
    "os"
    "strings"
    "unicode/utf8"
)

const MaxPreviewBytes = 128 * 1024 // 128 KiB

func PreviewFile(path string, limit int) (string, error) {
    fi, err := os.Stat(path)
    if err != nil { return "", err }
    if fi.IsDir() { return "", errors.New("path is a directory") }

    if fi.Size() > int64(limit) {
        return fmt.Sprintf("[content omitted: file is %d bytes (> %d)]", fi.Size(), limit), nil
    }

    f, err := os.Open(path)
    if err != nil { return "", err }
    defer f.Close()

    size := fi.Size()
    if size > int64(limit) {
        size = int64(limit)
    }
    buf := make([]byte, size)
    _, err = io.ReadFull(f, buf)
    if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) { return "", err }

    if looksBinary(buf) { return "[binary file — preview omitted]", nil }

    text := strings.ReplaceAll(string(buf), "\r\n", "\n")
    return text, nil
}

func looksBinary(b []byte) bool {
    if bytes.IndexByte(b, 0x00) >= 0 { return true }
    return !utf8.Valid(b)
}
