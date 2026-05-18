package shpool

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// List runs `shpool list` and returns parsed session names.
func List() ([]string, error) {
	path, err := LookPath()
	if err != nil {
		return nil, err
	}
	cmd := exec.Command(path, "list")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg != "" {
			return nil, fmt.Errorf("shpool list failed: %w: %s", err, msg)
		}
		return nil, fmt.Errorf("shpool list failed: %w", err)
	}
	return ParseList(stdout.Bytes()), nil
}

// ParseList extracts session names from `shpool list` output.
//
// The first column on each non-empty, non-header line is treated as the
// session name. Lines that look like a header row are skipped.
func ParseList(out []byte) []string {
	var names []string
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		first := fields[0]
		if looksLikeHeader(first, fields) {
			continue
		}
		names = append(names, first)
	}
	return names
}

// looksLikeHeader returns true for rows that appear to be column headers
// rather than real session entries.
func looksLikeHeader(first string, fields []string) bool {
	upper := strings.ToUpper(first)
	switch upper {
	case "NAME", "SESSION", "SESSIONS", "SESSION_NAME", "ID":
		return true
	}
	// All-caps multi-column rows are very likely headers (e.g. "NAME STATUS").
	if len(fields) >= 2 && isAllUpperWord(first) && isAllUpperWord(fields[1]) {
		return true
	}
	return false
}

func isAllUpperWord(s string) bool {
	if s == "" {
		return false
	}
	hasLetter := false
	for _, r := range s {
		switch {
		case r >= 'A' && r <= 'Z':
			hasLetter = true
		case r == '_' || r == '-':
			// allowed
		default:
			return false
		}
	}
	return hasLetter
}
