# Tech Stack

- Go module `github.com/uzulla/shpool-launch`; Go version pinned in `go.mod` as `go 1.25.10`, mise pins `go = "1.25"`.
- Dependencies are intentionally small: Bubble Tea v1 and Lipgloss for TUI; standard library `flag` for CLI parsing; `syscall.Exec` for attach handoff.
- Build/install/test tasks are defined in `mise.toml`; plain `go` commands also work.
