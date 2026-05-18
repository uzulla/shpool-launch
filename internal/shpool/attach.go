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

// LookPath returns the absolute path to the shpool binary or ErrNotFound.
func LookPath() (string, error) {
	p, err := exec.LookPath("shpool")
	if err != nil {
		return "", ErrNotFound
	}
	return p, nil
}

// Attach replaces the current process with `shpool attach [-f] <name>`.
// On success it does not return.
func Attach(name string, force bool) error {
	path, err := LookPath()
	if err != nil {
		return err
	}
	args := []string{"shpool", "attach"}
	if force {
		args = append(args, "-f")
	}
	args = append(args, name)
	if err := syscall.Exec(path, args, os.Environ()); err != nil {
		return fmt.Errorf("exec shpool: %w", err)
	}
	return nil
}
