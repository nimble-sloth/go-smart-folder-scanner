# Smart Folder Scanner (Go + Fyne)

Desktop app to pick a folder, recursively list directories/files, and show safe previews.

## Run
```bash
go run ./cmd/gui
```

## Build
Option 1
```bash
cd C:\Repo\go-smart-folder-scanner
go build -o bin/folder-scanner.exe ./cmd/gui
```

Option 2 - fancy icon
```bash
cd C:\Repo\go-smart-folder-scanner
rsrc -ico assets\icon.ico -o cmd\gui\rsrc.syso
go build -ldflags="-H=windowsgui" -o bin\folder-scanner.exe .\cmd\gui
```

## Notes
- Previews are limited to 128 KiB and skip binary files.
- Optional AI clients live in `internal/clients/ai` (OpenAI, Grok).
