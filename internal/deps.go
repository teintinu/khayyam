package internal

import (
	"os"
	"os/exec"
)

type InstallDependenciesOptions struct {
	Frozen bool
}

func InstallDependencies(repo *Repository, opts InstallDependenciesOptions) error {

	CleanRepository(repo)
	if err := ConfigureRepository(repo); err != nil {
		return err
	}
	for _, pkg := range repo.Packages {
		if err := configurePkg(repo, pkg); err != nil {
			return err
		}
	}

	execYarn := exec.Command("yarn")
	execYarn.Stdin = os.Stdin
	execYarn.Stdout = os.Stdout
	execYarn.Stderr = os.Stderr
	return execYarn.Run()

}
