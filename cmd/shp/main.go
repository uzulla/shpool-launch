// Command shp is a thin wrapper around `shpool attach` that generates
// session names from the current directory and offers a peco-style picker
// over `shpool list`.
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/uzulla/shpool-launch/internal/session"
	"github.com/uzulla/shpool-launch/internal/shpool"
	"github.com/uzulla/shpool-launch/internal/tui"
)

const usage = `shp - shpool attach helper

Usage:
  shp                      Show a TUI picker. The cwd-derived session name is
                           the default selection; existing sessions follow
                           when available. Enter attaches or creates when no
                           match is found.
  shp <session-name>       Attach to the given session name (no TUI).
  shp -f                   Show the TUI picker (like plain shp) and
                           force-attach to the chosen/created session.
  shp -f <session-name>    Force-attach to the given session name (no TUI).
  shp --print-name         Print the session name generated from the current
                           directory and exit.
  shp -h | --help          Show this help.

New sessions are created in the current directory.
In the home directory (or filesystem root), where no name can be derived from
a relative path, the session name falls back to one derived from the full path
(e.g. "home.<user>"), so shp still launches instead of dead-ending.
shpool must be installed and on PATH.
`

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	fs := flag.NewFlagSet("shp", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	fs.Usage = func() { fmt.Fprint(os.Stderr, usage) }

	var (
		force     bool
		printName bool
	)
	fs.BoolVar(&force, "f", false, "force attach")
	fs.BoolVar(&printName, "print-name", false, "print generated session name and exit")

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}
	rest := fs.Args()

	if printName {
		if len(rest) > 0 {
			return fmt.Errorf("--print-name takes no arguments")
		}
		name, err := session.FromCwdOrFallback()
		if err != nil {
			return err
		}
		fmt.Println(name)
		return nil
	}

	// With an explicit session name argument, skip the TUI and attach directly.
	if len(rest) == 1 {
		return shpool.Attach(rest[0], force)
	}
	if len(rest) > 1 {
		fs.Usage()
		return fmt.Errorf("too many arguments")
	}

	// `shp` / `shp -f` (no name): show the picker with the cwd name as the
	// default entry. `force` is carried through to the final attach, so
	// `shp -f` lets the user pick/create a session and force-attaches to it.
	return runPicker(force)
}

func runPicker(force bool) error {
	// Use the full-path fallback when the cwd yields no name (the home
	// directory itself or the filesystem root) so the picker always has a
	// default create-entry, just like any other directory, instead of
	// dead-ending with "no sessions".
	cwdName, err := session.FromCwdOrFallback()
	if err != nil {
		return err
	}

	sessions, err := shpool.List()
	if err != nil {
		if !errors.Is(err, shpool.ErrListUnavailable) {
			return err
		}
		sessions = nil
	}

	picked, err := tui.SelectWithDefault(sessions, cwdName)
	if err != nil {
		if errors.Is(err, tui.ErrCancelled) {
			return nil
		}
		return err
	}
	return shpool.Attach(picked, force)
}
