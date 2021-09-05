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
	"errors"

	"github.com/evanw/esbuild/pkg/api"
)

func Run(repo *Repository, packages []string) error {
	var cmds []*BuildJSResult
	for _, pkgName := range packages {
		pkg := repo.Packages[pkgName]
		if pkg == nil {
			return errors.New("no such package: " + pkgName)
		}
		if !pkg.Executable {
			return errors.New("no is a executable: " + pkgName)
		}
		if cmd, err := BundleWithEsbuild(
			repo,
			pkg,
			&BuildOpts{
				Target: api.ES2015,
				Minify: false,
				tab:    nil,
			},
		); err != nil {
			for _, running := range cmds {
				running.stop()
			}
			return err
		} else {
			cmds = append(cmds, cmd)
		}
	}
	return nil
}
