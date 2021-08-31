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

var RegExpJestRunStart = regexp.MustCompile(`^\s+jestRunStart\s+$`)

// jestRunComplete count=0 failed=0
var RegExpJestRunComplete = regexp.MustCompile(`^\s+jestRunComplete\s+count=(\d+)\s+failed=(\d+)\s+$`)

var RegExpJestReportCreated = regexp.MustCompile(`reporter is created on`)
