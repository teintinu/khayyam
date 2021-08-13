package internal

import (
	"os/exec"
	"path"
)

type TestOptions struct {
	Watch    bool
	Coverage bool
	Colors   bool
}

func CreateJestCommand(repo *Repository, opts TestOptions) *exec.Cmd {
	bin := path.Join(repo.RootDir, "node_modules", ".bin", "jest")
	var args = []string{}
	if opts.Watch {
		args = append(args, "--watch")
	}
	if opts.Coverage {
		args = append(args, "--coverage")
	}
	if !opts.Colors {
		args = append(args, "--no-colors")
	}
	jestCmd := exec.Command(bin, args...)
	return jestCmd
}
