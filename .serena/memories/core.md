# Core

- Go CLI wrapping `shpool attach`; executable entrypoint `cmd/shp/main.go`.
- Packages: `internal/session` derives session names from cwd/path; `internal/shpool` wraps `shpool` CLI list/attach; `internal/tui` is the Bubble Tea picker.
- `shp` with no args shows TUI using cwd-derived default plus existing `shpool list` sessions; explicit names bypass TUI and attach directly.
- Read `mem:tech_stack` for tools/deps, `mem:conventions` for implementation style, `mem:suggested_commands` for local commands, `mem:task_completion` before finishing coding tasks.
