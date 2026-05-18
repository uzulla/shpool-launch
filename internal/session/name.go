// Package session generates shpool session names from filesystem paths.
package session

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FromCwd returns a session name derived from the current working directory.
func FromCwd() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("get cwd: %w", err)
	}
	home, _ := os.UserHomeDir()
	return FromPath(cwd, home), nil
}

// FromPath converts an absolute path to a session name.
// If home is non-empty and path is under home, the home prefix is stripped first.
func FromPath(path, home string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		abs = path
	}
	abs = filepath.Clean(abs)

	if home != "" {
		h := filepath.Clean(home)
		if abs == h {
			abs = ""
		} else if strings.HasPrefix(abs, h+string(filepath.Separator)) {
			abs = abs[len(h)+1:]
		}
	}

	abs = strings.Trim(abs, "/")
	abs = strings.ReplaceAll(abs, "/", ".")

	var b strings.Builder
	b.Grow(len(abs))
	for _, r := range abs {
		switch {
		case r >= 'a' && r <= 'z',
			r >= 'A' && r <= 'Z',
			r >= '0' && r <= '9',
			r == '.', r == '-', r == '_':
			b.WriteRune(r)
		default:
			b.WriteRune('_')
		}
	}
	return b.String()
}
