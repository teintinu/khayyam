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
	folders := []string{
		"dist",
		"package*.json",
		"yarn*.json",
		"yarn.lock",
		"tsconfig*.json",
		".nvm.rc",
		".vscode",
		".jest",
		"jest.config.js",
		".eslintrc.json",
	}
	for _, folder := range folders {
		distDir := path.Join(repo.RootDir, pkg.Folder, folder)
		Logger.Debug("clearing ", distDir)
		err := os.RemoveAll(distDir)
		if err != nil {
			Logger.ErrorObj(err)
			return err
		}
	}
	return nil
}
