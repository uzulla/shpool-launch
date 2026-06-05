// Package session generates shpool session names from filesystem paths.
package session

import (
	"crypto/sha1"
	"encoding/hex"
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

// FromCwdOrFallback returns the cwd-derived session name, falling back to
// FallbackName when the cwd yields no name (the home directory itself or the
// filesystem root). The result is always non-empty, so callers can attach or
// print a name without dead-ending.
func FromCwdOrFallback() (string, error) {
	name, err := FromCwd()
	if err != nil {
		return "", err
	}
	if name != "" {
		return name, nil
	}
	return FallbackName()
}

// FallbackName returns a session name to use when FromCwd yields no name,
// which happens in the home directory itself or at the filesystem root. It
// derives the name from the full absolute path (without stripping the home
// prefix), so the home directory maps to e.g. "home.developer" rather than
// dead-ending. This keeps the name unique and consistent with FromPath's
// naming rules. It falls back to "shell" only when even the full path yields
// nothing (the filesystem root).
func FallbackName() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("get cwd: %w", err)
	}
	name := FromPath(cwd, "")
	if name == "" {
		name = "shell"
	}
	return name, nil
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

	rel := strings.Trim(filepath.ToSlash(abs), "/")
	if rel == "" {
		return ""
	}
	needsSuffix := needsHashSuffix(rel)
	nameInput := strings.ReplaceAll(rel, "/", ".")

	name := Sanitize(nameInput)
	if needsSuffix {
		name += "-" + shortHash(rel)
	}
	return name
}

// Sanitize maps a string to the shpool session-name character set used
// throughout this tool: only [A-Za-z0-9._-] are preserved and every other
// rune becomes '_'. It is the single normalization rule that both
// path-derived names and free-text (TUI) names go through, so that all
// session names share one safe alphabet.
func Sanitize(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
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

func needsHashSuffix(rel string) bool {
	for _, part := range strings.Split(rel, "/") {
		for _, r := range part {
			switch {
			case r >= 'a' && r <= 'z',
				r >= 'A' && r <= 'Z',
				r >= '0' && r <= '9',
				r == '-', r == '_':
				continue
			default:
				return true
			}
		}
	}
	return false
}

func shortHash(s string) string {
	sum := sha1.Sum([]byte(s))
	return hex.EncodeToString(sum[:])[:8]
}
