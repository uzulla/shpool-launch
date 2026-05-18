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
		msg := strings.TrimSpace(strings.Join([]string{stderr.String(), stdout.String()}, "\n"))
		if listUnavailable(msg) {
			if msg == "" {
				return nil, fmt.Errorf("%w: %v", ErrListUnavailable, err)
			}
			return nil, fmt.Errorf("%w: %s", ErrListUnavailable, msg)
		}
		if msg != "" {
			return nil, fmt.Errorf("shpool list failed: %w: %s", err, msg)
		}
		return nil, fmt.Errorf("shpool list failed: %w", err)
	}
	return ParseList(stdout.Bytes()), nil
}

func listUnavailable(msg string) bool {
	msg = strings.ToLower(strings.TrimSpace(msg))
	return msg == "" ||
		strings.Contains(msg, "could not connect to daemon") ||
		strings.Contains(msg, "control socket never came up") ||
		strings.Contains(msg, "connection refused")
}

// ParseList extracts session names from `shpool list` output.
//
// The first column on each non-empty, non-header line is treated as the
// session name. Lines matching shpool's table header are skipped.
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
		if looksLikeHeader(fields) {
			continue
		}
		names = append(names, fields[0])
	}
	return names
}

func looksLikeHeader(fields []string) bool {
	return len(fields) >= 2 &&
		strings.EqualFold(fields[0], "NAME") &&
		strings.EqualFold(fields[1], "STATUS")
}
