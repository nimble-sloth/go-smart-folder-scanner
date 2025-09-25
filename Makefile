run:
	go run ./cmd/gui

build:
	go build -ldflags "-X github.com/nimble-sloth/go-smart-folder-scanner/internal/version.Version=$$(git rev-parse --short HEAD)" -o bin/folder-scanner ./cmd/gui

test:
	go test ./...

fmt:
	go fmt ./...
