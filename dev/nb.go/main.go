package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	// "strings"
	"syscall"
)

type configuration struct {
	NBDir  string
	NBPath string
}

func (config *configuration) load() error {
	dir := os.Getenv("NB_DIR")

	if dir == "" {
		if runtime.GOOS == "windows" {
			dir = os.Getenv("APPDATA")

			if dir == "" {
				dir = filepath.Join(os.Getenv("USERPROFILE"), "Application Data", "nb")
			}

			dir = filepath.Join(dir, "nb")
		} else {
			dir = filepath.Join(os.Getenv("HOME"), ".nb")
		}
	}

	config.NBDir = dir

	var err error

	if config.NBPath, err = exec.LookPath("nb"); err != nil {
		return err
	}

	return nil
}

type subcommand struct {
	name  string
	usage string
}

func cmdRun(config configuration, args []string, env []string) error {
	var arguments []string

	if len(args) < 2 {
		return errors.New("command required")
	} else if 2 <= len(args) {
		arguments = args[2:]
	} else {
		arguments = []string{}
	}

	cmd := exec.Command(args[1], arguments...)
	cmd.Dir = config.NBDir

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

// run loads the configuration and environment, then runs the subcommand.
func run() error {
	var config configuration

	if err := config.load(); err != nil {
		return err
	}

	args := os.Args
	env := os.Environ()

	if len(args) > 1 && args[1] == "run" {
		if err := cmdRun(config, args[1:], env); err != nil {
			return err
		}
	} else {
		if execErr := syscall.Exec(config.NBPath, args, env); execErr != nil {
			return execErr
		}
	}

	return nil
}

// main is the primary entry point for the program.
func main() {
	os.Exit(presentError(run()))
}

// presentError translates an error into a message presnted to the user and
// returns the appropriate exit code.
func presentError(err error) int {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return 1
	}

	return 0
}