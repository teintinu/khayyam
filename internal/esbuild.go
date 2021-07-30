package internal

import (
	"os"
	"path"

	"github.com/evanw/esbuild/pkg/api"
)

type BuildJSResult struct {
	Errors   []api.Message
	Warnings []api.Message
}

func BundleWithEsbuild(repo *Repository, pkg *Package) (BuildJSResult, error) {

	pkgRoot := path.Join(repo.RootDir, pkg.Folder)
	distDir := path.Join(pkgRoot, "dist")
	err := os.MkdirAll(distDir, 0755)
	if err != nil {
		return BuildJSResult{}, err
	}

	buildOpts := api.BuildOptions{
		AbsWorkingDir: repo.RootDir,
		Outdir:        distDir,
		Bundle:        true,
		// Target: ES5,
		// MinifyIdentifiers: true,
		// MinifyWhitespace: true,
		// MinifySyntax: true,
		Platform:  api.PlatformNode,
		Format:    api.FormatCommonJS,
		Write:     true,
		LogLevel:  api.LogLevelWarning,
		Sourcemap: api.SourceMapLinked,
		Plugins:   []api.Plugin{},
		External:  getExternals(repo),
		Loader:    loaders,
		// TODO: Splitting: true,
	}

	entryPoint, err := GetPackageEntryPoint(repo, pkg)
	if err != nil {
		return BuildJSResult{}, err
	}
	buildOpts.EntryPoints = append(buildOpts.EntryPoints, entryPoint)

	result := api.Build(buildOpts)
	Logger.Debug("esbuild", buildOpts, result)
	return BuildJSResult{
		Errors:   result.Errors,
		Warnings: result.Warnings,
	}, nil
}
