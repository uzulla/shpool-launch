// Package shpool wraps the shpool CLI: locating it on PATH, replacing the
// current process with `shpool attach`, and parsing `shpool list` output.
package shpool

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// ErrNotFound is returned when the shpool binary is not on PATH.
var ErrNotFound = errors.New("shpool command not found in PATH. Please install shpool first")

// ErrListUnavailable is returned when shpool cannot provide the session list.
var ErrListUnavailable = errors.New("shpool session list unavailable")

// LookPath returns the absolute path to the shpool binary or ErrNotFound.
func LookPath() (string, error) {
	p, err := exec.LookPath("shpool")
	if err != nil {
		return "", ErrNotFound
	}
	return p, nil
}

// Attach replaces the current process with `shpool attach [-f] --dir . <name>`.
// On success it does not return.
func Attach(name string, force bool) error {
	path, err := LookPath()
	if err != nil {
		return err
	}
	args := attachArgs(name, force)
	if err := syscall.Exec(path, args, os.Environ()); err != nil {
		return fmt.Errorf("exec shpool: %w", err)
	}
	return nil
}

func attachArgs(name string, force bool) []string {
	args := []string{"shpool", "attach"}
	if force {
		args = append(args, "-f")
	}
	args = append(args, "--dir", ".")
	// Terminate option parsing with "--" so a session name that begins with
	// "-" is always treated as the positional <NAME>, never as a flag.
	args = append(args, "--", name)
	return args
}
