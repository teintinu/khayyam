package internal

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

func configurePkg(repo *Repository, pkg *Package) error {
	var packageVersion string
	if pkg.Version == "" {
		if repo.IsWorkspace {
			packageVersion = repo.Workspace.Version
		} else {
			packageVersion = "0.0.1"
		}
	} else {
		packageVersion = pkg.Version
	}

	metadata := PackageMetadata{
		Name:         pkg.Name,
		Version:      packageVersion,
		Description:  "GENERATED FILE: DO NOT EDIT! This file is managed by monoclean.",
		Dependencies: make(map[string]string),
		Repository:   repo.Url,
		Scripts: map[string]string{
			"postinstall": "patch-package",
		},
	}

	for dependencyName, dependency := range pkg.Dependencies {
		metadata.Dependencies[dependencyName] = dependency.Version
	}

	err := WritePackageJSON(metadata, repo.RootDir+"/"+pkg.Folder)
	if err == nil {
		err = configurePkgTsConfigReferences(repo, pkg)
	}
	if err == nil {
		err = initializeEmptySources(repo, pkg)
	}
	return err
}

func configurePkgTsConfigReferences(repo *Repository, pkg *Package) error {
	relativePathToRoot, err := filepath.Rel("/"+pkg.Folder, "/")
	if err != nil {
		return err
	}
	meta := TsConfigMetadata{
		Extends: relativePathToRoot + "/tsconfig.json",
		CompilerOptions: TsConfigCompileOptionsMetadata{
			OutDir:    "./dist/cjs",
			RootDir:   "./src",
			BaseURL:   ".",
			Composite: true,
			Paths:     make(map[string][]string),
			Lib:       []string{"ESNext"},
		},
		Exclude: []string{"dist"},
		Include: []string{"src/**/*.ts", "src/**/*.tsx"},
	}
	for _, dependency := range pkg.Dependencies {
		if dependency.Version != "*" {
			return errors.New("package dependencies is supported only inside workspace")
		}
		var depPkg = repo.Packages[dependency.Name]
		if depPkg == nil {
			return errors.New("package not found " + dependency.Name)
		}
		relativePathToDep, err := filepath.Rel(pkg.Folder, depPkg.Folder)
		if err != nil {
			return err
		}

		ref := TsConfigReferenceMetadata{
			Path: relativePathToDep,
		}
		meta.References = append(meta.References, ref)
		meta.CompilerOptions.Paths[dependency.Name] = []string{relativePathToDep + "/src"}
		if depPkg.usesDOM {
			meta.CompilerOptions.Lib = append(meta.CompilerOptions.Lib, "DOM")
		}
		if depPkg.usesWebWorker {
			meta.CompilerOptions.Lib = append(meta.CompilerOptions.Lib, "WebWorker")
		}
	}
	return WriteTsConfigJSON(meta, path.Join(repo.RootDir, pkg.Folder, "tsconfig.json"))
}

func initializeEmptySources(repo *Repository, pkg *Package) error {
	sourceDir := path.Join(repo.RootDir, pkg.Folder, "src")
	indexTs := path.Join(sourceDir, "index.ts")

	_, err := os.Stat(indexTs)
	if err == nil {
		return nil
	}

	err = os.MkdirAll(sourceDir, 0755)
	if err == nil {
		// TODO
		// if pkg.Layer == BusinessRulesLayer {
		// 	err = initializeBusinessRulesLayer(sourceDir, indexTs)
		// }
		err = initializeGenericLayer(sourceDir, indexTs)
	}
	return err
}

func initializeGenericLayer(sourceDir string, indexTs string) error {
	const indexTsContent = `
export function doSomething() {
  return 'something';
}
`
	if err := ioutil.WriteFile(indexTs, []byte(indexTsContent), 0644); err != nil {
		return err
	}
	var indexTestTs = path.Join(sourceDir, "index.test.ts")
	const indexTestContent = `
import {doSomething} from "./index"

it('test something' , () => {
  expect('anything').toEqual(doSomething())
})
`
	return ioutil.WriteFile(indexTestTs, []byte(indexTestContent), 0644)
}
