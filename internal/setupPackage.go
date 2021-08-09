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
		packageVersion = repo.Workspace.Version
	} else {
		packageVersion = pkg.Version
	}

	metadata := PackageMetadata{
		Name:            pkg.Name,
		Version:         packageVersion,
		Description:     "GENERATED FILE: DO NOT EDIT! This file is managed by monoclean.",
		Dependencies:    make(map[string]string),
		DevDependencies: make(map[string]string),
		Repository:      repo.Url,
		Scripts: map[string]string{
			"clean":        "monoclean clean",
			"buildWithTSC": "tsc -p ./tsconfig.build.json",
			"buildDTS":     "dts-bundle-generator --project ./tsconfig.build.json --verbose",
		},
	}

	for dependencyName, dependency := range pkg.Dependencies {
		metadata.Dependencies[dependencyName] = dependency.Version
	}
	for dependencyName, dependency := range pkg.devDependencies {
		metadata.DevDependencies[dependencyName] = dependency.Version
	}

	if pkg.Executable {
		metadata.Bin = "dist/main.js"
	} else {
		metadata.Types = "dist/index.d.ts"
		metadata.Main = "dist/index.js"
	}

	err := WritePackageJSON(metadata, repo.RootDir+"/"+pkg.Folder)
	if err == nil {
		err = configurePkgTsConfigTest(repo, pkg)
	}
	if err == nil {
		err = configurePkgTsConfigBuild(repo, pkg)
	}
	if err == nil {
		err = initializeEmptySources(repo, pkg)
	}
	return err
}

func configurePkgTsConfigTest(repo *Repository, pkg *Package) error {
	relativePathToRoot, err := filepath.Rel("/"+pkg.Folder, "/")
	if err != nil {
		return err
	}
	meta := TsConfigMetadata{
		Extends: relativePathToRoot + "/tsconfig.settings.json",
		CompilerOptions: TsConfigCompileOptionsMetadata{
			OutDir:          "./dist",
			RootDir:         "./src",
			TsBuildInfoFile: "dist/.tsbuildinfo",
			BaseURL:         ".",
			Composite:       true,
			Paths:           make(map[string][]string),
			Lib:             []string{"ESNext"},
		},
		Exclude: []string{
			"dist",
		},
		Include: []string{
			"src/**/*.ts",
			"src/**/*.tsx",
			"src/**/*.test.ts",
			"src/**/*.test.tsx",
			"src/**/*.spec.ts",
			"src/**/*.spec.tsx",
		},
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

func configurePkgTsConfigBuild(repo *Repository, pkg *Package) error {
	relativePathToRoot, err := filepath.Rel("/"+pkg.Folder, "/")
	if err != nil {
		return err
	}
	meta := TsConfigMetadata{
		Extends: relativePathToRoot + "/tsconfig.settings.json",
		CompilerOptions: TsConfigCompileOptionsMetadata{
			OutDir:          "./dist",
			RootDir:         "./src",
			TsBuildInfoFile: "dist/.tsbuildinfo",
			BaseURL:         ".",
			Composite:       true,
			Paths:           make(map[string][]string),
			Lib:             []string{"ESNext"},
		},
		Exclude: []string{
			"dist",
			"src/**/*.test.ts",
			"src/**/*.test.tsx",
			"src/**/*.spec.ts",
			"src/**/*.spec.tsx",
		},
		Include: []string{
			"src/**/*.ts",
			"src/**/*.tsx",
		},
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
	return WriteTsConfigJSON(meta, path.Join(repo.RootDir, pkg.Folder, "tsconfig.build.json"))
}

func initializeEmptySources(repo *Repository, pkg *Package) error {
	sourceDir := path.Join(repo.RootDir, pkg.Folder, "src")

	err := os.MkdirAll(sourceDir, 0755)
	if err == nil {
		// TODO
		// if pkg.Layer == BusinessRulesLayer {
		// 	err = initializeBusinessRulesLayer(sourceDir, indexTs)
		// }
		if pkg.Executable {
			err = initializeExecutable(sourceDir)
		} else {
			err = initializeGenericLayer(sourceDir)
		}
		if err != nil {
			return err
		}
	}
	return err
}

func initializeGenericLayer(sourceDir string) error {

	indexTs := path.Join(sourceDir, "index.ts")
	_, err := os.Stat(indexTs)
	if err == nil {
		return nil
	}

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

func initializeExecutable(sourceDir string) error {
	mainTs := path.Join(sourceDir, "main.ts")
	_, err := os.Stat(mainTs)
	if err == nil {
		return nil
	}
	const mainTsContent = `
export async function main() {
  return 'something';
}

main().catch(console.log)
`
	return ioutil.WriteFile(mainTs, []byte(mainTsContent), 0644)
}
