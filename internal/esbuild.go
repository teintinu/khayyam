package internal

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/evanw/esbuild/pkg/api"
)

type BuildJSResult struct {
	Errors   []api.Message
	Warnings []api.Message
	stop     func()
}

type EsbuildMode uint8

const (
	BuildOnly EsbuildMode = iota
	RunOnce
	WatchAndRun
)

type BuildOpts struct {
	Target api.Target
	Minify bool
	Mode   EsbuildMode
	tab    *WebTermTab
}

func BundleWithEsbuild(repo *Repository, pkg *Package, opts *BuildOpts) (*BuildJSResult, error) {

	pkgRoot := path.Join(repo.RootDir, pkg.Folder)
	distDir := path.Join(pkgRoot, "dist")
	err := os.MkdirAll(distDir, 0755)
	if err != nil {
		return nil, err
	}

	var watch *api.WatchMode
	var cmd *exec.Cmd

	if opts.Mode == WatchAndRun {
		watch = &api.WatchMode{
			OnRebuild: func(result api.BuildResult) {
				if len(result.Errors) > 0 {
					for _, msg := range result.Errors {
						if opts.tab == nil {
							fmt.Println(msg.Text)
						} else {
							opts.tab.consoleOutput(msg.Text)
						}
					}
					if opts.tab == nil {
						fmt.Printf("%v has %v errors", pkg.Name, len(result.Errors))
					} else {
						opts.tab.routines.setError(fmt.Sprintf("%v errors", len(result.Errors)))
						opts.tab.consoleOutput(fmt.Sprintf("%v has %v errors", pkg.Name, len(result.Errors)))
					}
				} else {
					cmd, err = run(repo, pkg, opts.tab, cmd)
					if err != nil {
						if opts.tab == nil {
							fmt.Println(err)
							cmd = nil
						} else {
							opts.tab.routines.setError(fmt.Sprintf("%v", err))
							opts.tab.consoleOutput(fmt.Sprintf("%v", err))
						}
					} else {
						if opts.tab != nil {
							opts.tab.routines.setSuccess("")
						}
					}
				}
			},
		}
	}

	buildOpts := api.BuildOptions{
		AbsWorkingDir:     repo.RootDir,
		Outdir:            distDir,
		Bundle:            true,
		Target:            opts.Target,
		MinifyIdentifiers: opts.Minify,
		MinifyWhitespace:  opts.Minify,
		MinifySyntax:      opts.Minify,
		Platform:          api.PlatformNode,
		Format:            api.FormatCommonJS,
		Write:             true,
		LogLevel:          api.LogLevelWarning,
		Sourcemap:         api.SourceMapLinked,
		Plugins:           []api.Plugin{},
		External:          getExternals(repo),
		Loader:            loaders,
		Watch:             watch,
		// TODO: Splitting: true,
	}

	entryPoint, err := GetPackageEntryPoint(repo, pkg)
	if err != nil {
		return nil, err
	}
	buildOpts.EntryPoints = append(buildOpts.EntryPoints, entryPoint)

	if watch == nil {
		esbuildResult := api.Build(buildOpts)
		Logger.Debug("esbuild", buildOpts, esbuildResult)
		if opts.Mode != BuildOnly {
			cmd, err = run(repo, pkg, opts.tab, cmd)
			if err != nil {
				if opts.tab == nil {
					fmt.Println(err)
					cmd = nil
				} else {
					opts.tab.routines.setError(fmt.Sprintf("%v", err))
					opts.tab.consoleOutput(fmt.Sprintf("%v", err))
				}
			}
		}
		return &BuildJSResult{
			Errors:   esbuildResult.Errors,
			Warnings: esbuildResult.Warnings,
			stop: func() {
				if cmd != nil {
					cmd.Process.Kill()
				}
				esbuildResult.Stop()
			},
		}, nil
	} else {
		result := api.Build(buildOpts)
		Logger.Debug("esbuild-watch", buildOpts, result)
		return nil, nil
	}
}

func run(repo *Repository, pkg *Package, tab *WebTermTab, old *exec.Cmd) (*exec.Cmd, error) {
	if old != nil {
		old.Process.Kill()
	}
	node := exec.Command("node", pkg.Bin)
	if tab == nil {
		node.Stdin = os.Stdin
		node.Stdout = os.Stdout
		node.Stderr = os.Stderr
		if err := node.Run(); err != nil {
			return nil, err
		}
	} else {
		go tab.runInPty(node)
	}
	return node, nil
}
