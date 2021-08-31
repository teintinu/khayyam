package internal

import (
	"os"
	"path"
)

func CleanRepository(repo *Repository) error {
	err := cleanFolder(repo, "")
	if err != nil {
		return err
	}
	for _, pkg := range repo.Packages {
		err := cleanFolder(repo, pkg.Folder)
		if err != nil {
			return err
		}
	}
	return nil
}

func cleanFolder(repo *Repository, parentfolder string) error {
	subfolders := []string{
		"dist",
		"node_modules",
		"package.json",
		"package-lock.json",
		"yarn.lock",
		"yarn-error.log",
		"tsconfig.json",
		"tsconfig.settings.json",
		"tsconfig.build.json",
		"tsconfig.test.json",
		".vscode",
		".jest",
		"jest.config.js",
		".eslintrc.json",
	}
	for _, subfolder := range subfolders {
		distDir := path.Join(repo.RootDir, parentfolder, subfolder)
		Logger.Debug("clearing ", distDir)
		err := os.RemoveAll(distDir)
		if err != nil {
			Logger.ErrorObj(err)
			return err
		}
	}
	return nil
}
