// Cases to handle
//
// without watch
//   build error
//   build ok
//     program fails to start
//     program terminates success
//     program terminates failure
//
// with watch
//   build error
//   build ok                              waitForChange
//     program fails to start                  true      wait for change
//     program terminates prematurely          true      wait for change
//     code changes                            false     restart
//     interrupt                               false     exi

package internal

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
)

type TestOptions struct {
	Watch bool
}

// TODO: Need to handle interrupts in order to have a higher chance
// of cleaning up temporary files.

// Status code may be returend within an exec.ExitError return value.
func Test(repo *Repository, opts TestOptions) error {
	if err := EnsureTmp(repo); err != nil {
		return err
	}

	dir, err := TempDir(repo, "test")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	testfiles, err := listTests(repo.RootDir)
	if err != nil {
		return err
	}
	var outfile = path.Join(dir, "bundle.test.js")
	jestCommand, err := makeJestCommand(repo, path.Join(dir, "jest.config.js"), outfile)
	if err != nil {
		return err
	}
	return buildAndWatch{
		Repository: repo,
		Watch:      opts.Watch,
		Esbuild: api.BuildOptions{
			AbsWorkingDir: repo.RootDir,
			EntryPoints:   testfiles,
			Outfile:       outfile,
			Bundle:        true,
			Platform:      api.PlatformNode,
			Format:        api.FormatCommonJS,
			Write:         true,
			LogLevel:      api.LogLevelWarning,
			Sourcemap:     api.SourceMapLinked,
			External:      getExternals(repo),
			Loader:        loaders,
		},
		CreateProcess: func() process {

			nodeArgs := jestCommand
			node := exec.Command("npx", nodeArgs...)
			node.Stdin = os.Stdin
			node.Stdout = os.Stdout
			node.Stderr = os.Stderr

			return &cmdProcess{cmd: node}
		},
	}.Run()
}

func listTests(root string) ([]string, error) {
	var tests = []string{}
	var err = filepath.WalkDir(root, func(path string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if !d.IsDir() && !(strings.Contains(path, "/node_modules/")) {
			if strings.HasSuffix(path, ".test.ts") ||
				strings.HasSuffix(path, ".test.tsx") ||
				strings.HasSuffix(path, ".spec.ts") ||
				strings.HasSuffix(path, ".spec.tsx") {
				println(path)
				tests = append(tests, path)
			}
		}
		return nil
	})
	return tests, err
}

func makeJestCommand(repo *Repository, jestconfig string, testBundle string) ([]string, error) {
	// TODO moduleNameMapper / rootDir
	// https://stackoverflow.com/questions/51799300/run-tests-against-compiled-bundles
	for _, pkg := range repo.Packages {
		fmt.Printf("%+v\n", pkg)
	}
	script := `require('source-map-support').install();
module.exports = {
  verbose: true,
  moduleNameMapper: {
    ".scss$": "scss-stub.js"
  }
}
`
	if err := ioutil.WriteFile(jestconfig, []byte(script), 0644); err != nil {
		return nil, err
	}

	return []string{"jest", "--config", jestconfig, testBundle}, nil
}
