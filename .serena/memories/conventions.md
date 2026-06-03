# Conventions

- Keep package boundaries narrow: CLI branching in `cmd/shp`, path-to-name logic in `internal/session`, shpool process/list parsing in `internal/shpool`, picker state/rendering in `internal/tui`.
- Session names are cwd/home-relative path strings with `/` mapped to `.`, only `[A-Za-z0-9._-]` preserved, other chars mapped to `_`, and short hash suffix added when needed for collision mitigation.
- Attach always uses `shpool attach [ -f ] --dir . <name>` so newly-created sessions start in the current directory.
- TUI filtering is case-insensitive substring AND over whitespace-separated tokens; Backspace deletes by rune; default cwd item is prepended/deduped and labeled `(cwd)`.
