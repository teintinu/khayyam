package internal

import (
	"os"
	"os/exec"
	"path"
)

type TestOptions struct {
	Watch    bool
	Coverage bool
}

func Test(repo *Repository, opts TestOptions) error {
	bin := path.Join(repo.RootDir, "node_modules", ".bin", "jest")
	var args = []string{}
	if opts.Watch {
		args = append(args, "--watch")
	}
	if opts.Coverage {
		args = append(args, "--coverage")
	}
	jest := exec.Command(bin, args...)
	jest.Stdin = os.Stdin
	jest.Stdout = os.Stdout
	jest.Stderr = os.Stderr
	return jest.Run()
}
