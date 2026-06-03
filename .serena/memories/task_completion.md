# Task Completion

- Run `gofmt` on changed Go files.
- Run `go test ./...` for coding changes.
- Run `go vet ./...` when behavior crosses package boundaries or before release-level changes; `mise run check` covers vet + tests.
- For TUI behavior changes, add/adjust unit tests in `internal/tui/select_test.go`; avoid relying only on manual interactive checks.
