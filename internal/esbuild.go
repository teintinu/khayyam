package internal

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/evanw/esbuild/pkg/api"
)

type BuildJSResult struct {
	Errors   func() []api.Message
	Warnings func() []api.Message
	stop     func()
}

type BuildOpts struct {
	Target api.Target
	Minify bool
	tab    *WebTermTab
}

func newBuildOpts(repo *Repository, pkg *Package, opts *BuildOpts) (
	pkgRoot string,
	entryPoint string,
	buildOpts api.BuildOptions,
	err error,
) {
	pkgRoot = path.Join(repo.RootDir, pkg.Folder)
	distDir := path.Join(pkgRoot, "dist")
	err = os.MkdirAll(distDir, 0755)
	if err != nil {
		return
	}

	entryPoint, err = GetPackageEntryPoint(repo, pkg)
	if err != nil {
		return
	}

	buildOpts = api.BuildOptions{
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
		LogLevel:          api.LogLevelVerbose,
		Sourcemap:         api.SourceMapLinked,
		Plugins:           []api.Plugin{},
		External:          getExternals(repo),
		Loader:            loaders,
		// TODO: Splitting: true,
	}
	if Logger.logging == int(LoggingDebug) {
		buildOpts.LogLevel = api.LogLevelVerbose
	}
	if Logger.logging == int(LoggingError) {
		buildOpts.LogLevel = api.LogLevelError
	}
	if Logger.logging == int(LoggingWarn) {
		buildOpts.LogLevel = api.LogLevelWarning
	}
	if Logger.logging == int(LoggingInfo) {
		buildOpts.LogLevel = api.LogLevelInfo
	}

	buildOpts.EntryPoints = append(buildOpts.EntryPoints, entryPoint)
	return
}

func BundleWithEsbuild(repo *Repository, pkg *Package, opts *BuildOpts) (*BuildJSResult, error) {
	_, entryPoint, buildOpts, err := newBuildOpts(repo, pkg, opts)
	if err != nil {
		return nil, err
	}
	esbuildResult := api.Build(buildOpts)
	Logger.Debug("esbuild-BuildOnly", entryPoint)
	return &BuildJSResult{
		Errors: func() []api.Message {
			return esbuildResult.Errors
		},
		Warnings: func() []api.Message {
			return esbuildResult.Warnings
		},
		stop: nil,
	}, nil
}

func BundleAndRunWithEsbuild(repo *Repository, pkg *Package, opts *BuildOpts) (*BuildJSResult, error) {

	_, entryPoint, buildOpts, err := newBuildOpts(repo, pkg, opts)
	if err != nil {
		return nil, err
	}

	var cmd *exec.Cmd
	esbuildResult := api.Build(buildOpts)
	Logger.Debug("esbuild-RunOnce", entryPoint)
	cmd = run(repo, pkg, opts.tab, cmd)

	return &BuildJSResult{
		Errors: func() []api.Message {
			return esbuildResult.Errors
		},
		Warnings: func() []api.Message {
			return esbuildResult.Warnings
		},
		stop: func() {
			if cmd != nil {
				cmd.Process.Kill()
			}
			esbuildResult.Stop()
		},
	}, nil
}

func WatchWithEsbuild(repo *Repository, pkg *Package, opts *BuildOpts) (*BuildJSResult, error) {

	pkgRoot, entryPoint, buildOpts, err := newBuildOpts(repo, pkg, opts)
	if err != nil {
		return nil, err
	}
	esbuildResult := api.Build(buildOpts)
	var cmd *exec.Cmd
	rebuildAndRun := func() {
		Logger.Debug("building: ", entryPoint)
		esbuildResult.Rebuild()
		cmd = run(repo, pkg, opts.tab, cmd)
	}
	var stopWatching func()
	cancelRebuild := func(canceled func()) {
		if cmd != nil {
			cmd.Process.Kill()
			cmd = nil
		}
		esbuildResult.Stop()
	}
	stopWatching = WatchFolder(path.Join(pkgRoot, "src"), rebuildAndRun, cancelRebuild, true)

	return &BuildJSResult{
		Errors: func() []api.Message {
			return esbuildResult.Errors
		},
		Warnings: func() []api.Message {
			return esbuildResult.Warnings
		},
		stop: func() {
			stopWatching()
			cancelRebuild(nil)
		},
	}, nil
}

func run(repo *Repository, pkg *Package, tab *WebTermTab, old *exec.Cmd) *exec.Cmd {
	Logger.Debug("esbuild-run", pkg.Bin)
	if old != nil {
		old.Process.Kill()
	}
	node := exec.Command("node", pkg.Bin)
	if tab == nil {
		node.Stdin = os.Stdin
		node.Stdout = os.Stdout
		node.Stderr = os.Stderr
		if err := node.Run(); err != nil {
			runLogger(tab, err)
			return nil
		}
	} else {
		go tab.runInPty(node)
	}
	return node
}

func runLogger(tab *WebTermTab, args ...interface{}) {
	if tab == nil {
		fmt.Println(args...)
	} else {
		tab.routines.setError(fmt.Sprintf("%v", args...))
		tab.consoleOutput(fmt.Sprintf("%v", args...))
	}
}
