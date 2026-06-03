# Suggested Commands

- `mise install` installs the pinned Go toolchain.
- `mise run build` builds `./bin/shp`; equivalent: `go build -o ./bin/shp ./cmd/shp`.
- `mise run test` runs `go test ./...`; `mise run vet` runs `go vet ./...`; `mise run check` runs both.
- `go run ./cmd/shp --print-name` exercises session-name generation without requiring `shpool` runtime behavior.
- `go run ./cmd/shp --help` checks CLI usage text.