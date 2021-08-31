package internal

import (
	"io/ioutil"
	"os"
	"path"
)

func ConfigureRepository(repo *Repository) error {
	var workspaceName string
	var workspaceVersion string

	workspaceName = repo.Workspace.Name
	workspaceVersion = repo.Workspace.Version
	if err := configureGitIgnore(repo); err != nil {
		return err
	}
	if err := configureVsCodeSettings(repo); err != nil {
		return err
	}
	if err := configureVsCodeRecommendedExtensions(repo); err != nil {
		return err
	}
	if err := configureNvmRc(repo); err != nil {
		return err
	}
	if err := configureRootTsConfigSettings(repo); err != nil {
		return err
	}
	if err := configureRootTsConfigReferences(repo); err != nil {
		return err
	}
	if err := configureRootTsConfigTest(repo); err != nil {
		return err
	}
	if err := configureJest(repo); err != nil {
		return err
	}
	if err := configureJestCustomReport(repo); err != nil {
		return err
	}
	if err := configureEsLint(repo); err != nil {
		return err
	}

	metadata := PackageMetadata{
		Name:            workspaceName,
		Version:         workspaceVersion,
		Private:         true,
		Description:     "GENERATED FILE: DO NOT EDIT! This file is managed by khayyam.",
		DevDependencies: make(map[string]string),
		Repository:      repo.Url,
	}

	metadata.Workspaces = []string{}
	for _, pkg := range repo.Packages {
		metadata.Workspaces = append(metadata.Workspaces, pkg.Folder)
	}

	metadata.Scripts = map[string]string{
		"clean":     "khayyam clean",
		"start":     "khayyam run",
		"build":     "khayyam build",
		"publish":   "khayyam publish",
		"tsc-watch": "tsc -p . --watch",
	}

	for dependencyName, dependency := range repo.DevDependencies {
		metadata.DevDependencies[dependencyName] = dependency.Version
	}

	return WritePackageJSON(metadata, repo.RootDir)
}

func configureGitIgnore(repo *Repository) error {

	gitignore := path.Join(repo.RootDir, ".gitignore")
	var content = `# khayyam
node_modules
.nvm.rc
package.json
package-lock.json
yarn-error.log
yarn.lock
tsconfig.**
tmp
out
.eslintrc.json
coverage
.vscode
.jest
jest.config.js
`

	err := ioutil.WriteFile(gitignore, []byte(content), 0644)
	return err
}

func configureVsCodeSettings(repo *Repository) error {
	var vscodeDir = path.Join(repo.RootDir, ".vscode")
	if err := os.MkdirAll(vscodeDir, 0755); err != nil {
		return err
	}
	var extensionsJson = path.Join(vscodeDir, "settings.json")
	var content = `{
	"search.exclude": {
		"**/node_modules": true,
		"**/bower_components": true,
		"**/env": true,
		"**/venv": true
	},
	"files.watcherExclude": {
		"**/.git/objects/**": true,
		"**/.git/subtree-cache/**": true,
		"**/node_modules/**": true,
		"**/env/**": true,
		"**/venv/**": true,
		"env-*": true
	},
  "yaml.schemas": {
    "https://teintinu.github.io/khayyam/khayyam-schema.json": [
      "khayyam.yml"
    ]
  },
	"typescript.disableAutomaticTypeAcquisition": true
}
`
	err := ioutil.WriteFile(extensionsJson, []byte(content), 0644)
	return err
}

func configureVsCodeRecommendedExtensions(repo *Repository) error {
	var vscodeDir = path.Join(repo.RootDir, ".vscode")
	var extensionsJson = path.Join(vscodeDir, "extensions.json")
	var content = `{
  "recommendations": [
    "dbaeumer.vscode-eslint",
		"redhat.vscode-yaml"
  ]
}
`
	err := ioutil.WriteFile(extensionsJson, []byte(content), 0644)
	return err
}

func configureNvmRc(repo *Repository) error {

	npmRc := path.Join(repo.RootDir, ".nvm.rc")
	var content = repo.Engines["node"]
	if len(content) == 0 {
		return nil
	}
	err := ioutil.WriteFile(npmRc, []byte(content), 0644)
	return err
}

func configureRootTsConfigSettings(repo *Repository) error {
	meta := TsConfigMetadata{
		CompilerOptions: TsConfigCompileOptionsMetadata{
			Incremental:      true,
			Target:           "ESNext",
			Module:           "ESNext",
			Declaration:      true,
			SourceMap:        true,
			ImportHelpers:    true,
			Strict:           true,
			ModuleResolution: "node",
			EsModuleInterop:  true,
			RootDir:          ".",
			BaseURL:          ".",
			Paths:            make(map[string][]string),
		},
		Exclude: []string{"node_modules", "dist", "__tests__", "**/*.test*", "**/*.spec*"},
	}
	return WriteTsConfigJSON(meta, path.Join(repo.RootDir, "tsconfig.settings.json"))
}

func configureRootTsConfigReferences(repo *Repository) error {
	meta := TsConfigMetadata{
		References: []TsConfigReferenceMetadata{},
	}
	for _, pkg := range repo.Packages {
		ref := TsConfigReferenceMetadata{
			Path: "./" + pkg.Folder,
		}
		meta.References = append(meta.References, ref)
	}
	return WriteTsConfigJSON(meta, path.Join(repo.RootDir, "tsconfig.json"))
}

func configureRootTsConfigTest(repo *Repository) error {
	meta := TsConfigMetadata{
		Extends: "./tsconfig.settings.json",
		CompilerOptions: TsConfigCompileOptionsMetadata{
			Target:  "ES2021",
			RootDir: ".",
			BaseURL: ".",
			Paths:   make(map[string][]string),
		},
		Exclude: []string{"node_modules", "dist"},
	}
	for _, pkg := range repo.Packages {
		meta.CompilerOptions.Paths[pkg.Name] = []string{"./" + pkg.Folder + "/src"}
	}
	return WriteTsConfigJSON(meta, path.Join(repo.RootDir, "tsconfig.test.json"))
}

func configureJest(repo *Repository) error {
	const content = `
const { pathsToModuleNameMapper } = require('ts-jest/utils')
const { compilerOptions } = require('./tsconfig.test')

const moduleNameMapper = pathsToModuleNameMapper(compilerOptions.paths, { prefix: '<rootDir>/' })

module.exports = {
  preset: 'ts-jest',
  modulePathIgnorePatterns: ['dist'],
  testPathIgnorePatterns: ['node_modules', 'dist', '.jest'],
  testRegex: '(\\.(test|spec|steps))\\.(ts|tsx)$',
  globals: {
    'ts-jest': {
      tsConfig: 'tsconfig.test.json'
    }
  },
  moduleNameMapper,
  transform: {
    '^.+\\.tsx?$': [
      'esbuild-jest',
      {
        sourcemap: 'inline',
        target: ['es6', 'node12'],
        loaders: {
          '.spec.ts': 'tsx',
          '.test.ts': 'tsx',
          '.steps.ts': 'tsx'
        }
      }
    ]
  },
  reporters: [
    './.jest/jest.report',
    ['jest-html-reporters', {
      publicPath: './.jest/html-report',
      filename: 'index.html',
      expand: true
    }]
  ],
  coverageDirectory: './.jest/coverage/',
  coverageReporters: [
    'text',
    'html',
    'cobertura',
    'json-summary'
  ],
  coverageThreshold: {
    global: {
      lines: 90,
      statements: 90,
      functions: 90,
      branches: 90
    }
  }
}
`
	jestConfig := path.Join(repo.RootDir, "jest.config.js")
	err := ioutil.WriteFile(jestConfig, []byte(content), 0644)
	return err
}

func configureJestCustomReport(repo *Repository) error {
	var vscodeDir = path.Join(repo.RootDir, ".jest")
	if err := os.MkdirAll(vscodeDir, 0755); err != nil {
		return err
	}
	const content = `
module.exports = class CustomReport {
  onRunStart () {
    this.testcount = 0
    this.failcount = 0
    console.log('')
    console.log('\njestRunStart\n')
  }

  onRunComplete () {
    console.log('\njestRunComplete count=' + this.testcount + ' failed=' + this.failcount +"\n")
  }

  onTestResult (test, testResult) {
		this.testcount++
    this.failcount += testResult.numFailingTests
  }
}
`
	jestConfig := path.Join(repo.RootDir, ".jest/jest.report.js")
	err := ioutil.WriteFile(jestConfig, []byte(content), 0644)
	return err
}

func configureEsLint(repo *Repository) error {
	const content = `{
  "env": {
    "browser": true,
    "es2021": true,
    "node": true,
    "jest": true
  },
  "extends": [
    "plugin:react/recommended",
    "standard"
  ],
  "parser": "@typescript-eslint/parser",
  "parserOptions": {
    "ecmaFeatures": {
      "jsx": true
    },
    "ecmaVersion": 12,
    "sourceType": "module"
  },
  "plugins": [
    "react",
    "@typescript-eslint"
  ],
  "rules": {
  },
  "settings": {
    "react": {
      "version": "detect"
    }
  }
}
`
	eslintrc := path.Join(repo.RootDir, ".eslintrc.json")
	err := ioutil.WriteFile(eslintrc, []byte(content), 0644)
	return err
}
