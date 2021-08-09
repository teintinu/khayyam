package internal

import (
	"github.com/teintinu/gjobs"
)

func BuildWorkspace(repo *Repository) error {
	jobs := gjobs.NewJobs()
	gjobs.VerbosityEvent = func(args ...interface{}) {
		Logger.Debug(args...)
	}
	MakeJobs(jobs, "build", repo, repo.PackageNames, func(pkg *Package) error {
		if pkg.Executable {
			return bundleExecutable(repo, pkg)
		} else {
			return buildLibrary(repo, pkg)
		}
	})
	jobs.Run()
	return nil
}

func bundleExecutable(repo *Repository, pkg *Package) error {
	Logger.Info("Bundling executable:", pkg.Name)
	msgs, err := BundleWithEsbuild(repo, pkg)
	if err != nil {
		Logger.ErrorObj(err)
	}
	for _, msg := range msgs.Errors {
		Logger.Error(msg.Text)
	}
	Logger.Debug("Bundled  executable:", pkg.Name)
	return err
}

func buildLibrary(repo *Repository, pkg *Package) error {
	Logger.Info("Building library   :", pkg.Name)
	err := BuildWithTSC(repo, pkg)
	if err != nil {
		Logger.ErrorObj(err)
	}
	Logger.Debug("Built   library   :", pkg.Name)
	return err
}
