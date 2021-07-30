package internal

import (
	"fmt"
	"os"
	"os/exec"
	"path"
)

func BuildWithTSC(repo *Repository, pkg *Package) error {
	args := []string{
		"--project",
		path.Join(repo.RootDir, pkg.Folder),
	}

	bin := path.Join(repo.RootDir, "node_modules", ".bin", "tsc")
	cmd := exec.Command(
		bin,
		args...,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("bundling type declarations: %w", err)
	}
	Logger.Debug(bin, args)
	return nil
}
