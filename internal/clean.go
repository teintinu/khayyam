package internal

import (
	"os"
	"path"
)

func CleanRepository(repo *Repository) error {
	for _, pkg := range repo.Packages {
		err := CleanPackage(repo, pkg)
		if err != nil {
			return err
		}
	}
	return nil
}

func CleanPackage(repo *Repository, pkg *Package) error {
	distDir := path.Join(repo.RootDir, pkg.Folder, "dist")
	Logger.Debug("clearing ", distDir)
	err := os.RemoveAll(distDir)
	if err != nil {
		Logger.ErrorObj(err)
	}
	return err
}
