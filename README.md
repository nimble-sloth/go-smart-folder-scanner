# Smart Folder Scanner (Go + Fyne)

Desktop app to pick a folder, recursively list directories/files, and show safe previews.

## Run
```bash
go run ./cmd/gui
```

## Build
```bash
go build -o bin/folder-scanner ./cmd/gui
```

## Notes
- Previews are limited to 128 KiB and skip binary files.
- Optional AI clients live in `internal/clients/ai` (OpenAI, Grok).
