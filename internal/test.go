package internal

import (
	"fmt"
	"os/exec"
	"path"
	"regexp"
)

type TestOptions struct {
	Watch    bool
	Coverage bool
	Colors   bool
}

func CreateJestCommand(repo *Repository, opts TestOptions) *exec.Cmd {
	// bin := "/bin/bash"
	// var args = []string{"-c"}
	// args = append(args, path.Join(repo.RootDir, "node_modules", ".bin", "jest"))
	bin := path.Join(repo.RootDir, "node_modules", ".bin", "jest")
	var args = []string{
		"--config",
		path.Join(repo.RootDir, "jest.config.js"),
	}
	if opts.Watch {
		args = append(args, "--watch")
	}
	if opts.Coverage {
		args = append(args, "--coverage")
	}
	if !opts.Colors {
		args = append(args, "--no-colors")
	}
	fmt.Println(bin, args)
	jestCmd := exec.Command(bin, args...)
	return jestCmd
}

// Tests:       5 passed, 5 total
// Tests:       1 failed, 4 passed, 5 total
var RegExpJestRunComplete = regexp.MustCompile(`^\s*Tests:\s+(?:(\d+)\s+failed,)?\s*(?:(\d+).+passed,)?\s*(\d+)\s+total\s*$`)

var RegExpJestReportCreated = regexp.MustCompile(`reporter is created on`)
