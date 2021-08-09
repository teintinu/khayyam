package internal

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/goccy/go-yaml"
	"github.com/teintinu/gjobs"
)

type Repository struct {
	ConfigPath      string
	RootDir         string
	Engines         map[string]string
	Workspace       Workspace
	DevDependencies map[string]*Dependency
	Url             string
	Registry        string

	Packages     map[string]*Package
	PackageNames []string
}

type Workspace struct {
	Name    string
	Version string
}

type Dependency struct {
	Name    string
	Version string
}

type PackageLayer = int

const (
	NormalLayer PackageLayer = iota
	BusinessRulesLayer
	ExecutablesLayer
	AdaptersLayer
)

type PackagePublish = int

const (
	DontPublish PackagePublish = iota
	PublishRestrictly
	PublishPublicly
)

type Package struct {
	Name            string
	Version         string
	Folder          string
	Publish         PackagePublish
	usesNode        bool
	usesDOM         bool
	usesWebWorker   bool
	Description     string
	Dependencies    map[string]*Dependency
	devDependencies map[string]*Dependency
	Layer           PackageLayer
	Executable      bool
	Main            string
	Bin             string
	Types           string
}

type Executable struct {
	Name       string
	Entrypoint string
}

const DefaultRegistry = "https://registry.npmjs.org/"

func LoadRepository(searchDir string, checkEntryPoints bool) (*Repository, error) {
	f, err := openConfigFile(searchDir)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var repo Repository
	repo.ConfigPath = f.Name()
	repo.RootDir = path.Dir(repo.ConfigPath)

	dec := yaml.NewDecoder(f, yaml.Strict())
	var cfg Config
	err = dec.Decode(&cfg)
	if err == io.EOF {
		// TODO: Should entirely empty config files be allowed? Probably not.
		err = nil
	}
	if err != nil {
		return nil, err
	}

	repo.Workspace.Name = cfg.Workspace.Name
	repo.Workspace.Version = cfg.Workspace.Version

	repo.Engines = make(map[string]string)
	for engineName, engineVersion := range cfg.Engines {
		repo.Engines[engineName] = engineVersion
	}

	repo.Url = cfg.Repository
	repo.Registry = cfg.Registry
	if repo.Registry == "" {
		repo.Registry = DefaultRegistry
	}

	repo.Packages = make(map[string]*Package)
	err = loadPackages(&repo, cfg.Packages, NormalLayer, checkEntryPoints)
	if err != nil {
		return nil, err
	}
	err = loadPackages(&repo, cfg.BusinessRules, BusinessRulesLayer, checkEntryPoints)
	if err != nil {
		return nil, err
	}
	err = loadPackages(&repo, cfg.Executables, ExecutablesLayer, checkEntryPoints)
	if err != nil {
		return nil, err
	}
	err = loadPackages(&repo, cfg.Adapters, AdaptersLayer, checkEntryPoints)
	if err != nil {
		return nil, err
	}

	repo.DevDependencies = make(map[string]*Dependency)
	addWorkspaceDependency := func(name, version string) {
		repo.DevDependencies[name] = &Dependency{
			Name:    name,
			Version: version,
		}
	}

	for dependencyName, dependencyVersion := range requiredWorkspaceDevDependencies {
		if _, ok := repo.DevDependencies[dependencyName]; !ok {
			addWorkspaceDependency(dependencyName, dependencyVersion)
		}
	}
	for pkgName := range repo.Packages {
		pkg := repo.Packages[pkgName]
		if checkEntryPoints {
			_, err = GetPackageEntryPoint(&repo, pkg)
			if err != nil {
				return nil, err
			}
		}
		err = validateLayer(&repo, pkg)
		if err != nil {
			return nil, err
		}
	}
	return &repo, nil
}

func loadPackages(repo *Repository, packages map[string]PackageConfig, layer PackageLayer, checkEntryPoints bool) error {
	for packageName, packageConfig := range packages {

		if len(packageConfig.Index) > 0 {
			return errors.New("use folder instead index inside workspace")
		}
		if packageConfig.Folder == "" {
			return errors.New("use folder for each package in workspace")
		}
		if strings.Contains(packageConfig.Folder, "\\") {
			return errors.New(packageConfig.Folder + " user normal slashes")
		}
		if strings.HasPrefix(packageConfig.Folder, ".") || strings.HasPrefix(packageConfig.Folder, "/") {
			return errors.New(packageConfig.Folder + " must be relative to workspace folder")
		}
		if strings.HasSuffix(packageConfig.Folder, ".") || strings.HasSuffix(packageConfig.Folder, "/") {
			return errors.New(packageConfig.Folder + " folder ends with a invalid char")
		}

		pkg := &Package{
			Name:            packageName,
			Description:     packageConfig.Description,
			Folder:          packageConfig.Folder,
			Layer:           layer,
			Executable:      layer == ExecutablesLayer || (layer == NormalLayer && packageConfig.Executable),
			Dependencies:    make(map[string]*Dependency),
			devDependencies: make(map[string]*Dependency),
		}
		pkgFolder := path.Join(repo.RootDir, pkg.Folder)
		if packageConfig.Publish == "public" {
			pkg.Publish = PublishPublicly
		} else if packageConfig.Publish == "restrict" {
			pkg.Publish = PublishRestrictly
		} else {
			pkg.Publish = DontPublish
		}

		for dependencyName, dependencyVersion := range packageConfig.Dependencies {
			pkg.Dependencies[dependencyName] = &Dependency{
				Name:    dependencyName,
				Version: dependencyVersion,
			}
			pkg.usesDOM = dependencyName == "react-dom"
			pkg.usesNode = dependencyName == "@types/node"
		}

		for depName, depVersion := range requiredPackageDevDependencies {
			pkg.devDependencies[depName] = &Dependency{
				Name:    depName,
				Version: depVersion,
			}
		}

		if pkg.Executable {
			if _, err := GetPackageEntryPoint(repo, pkg); err != nil && checkEntryPoints {
				return err
			} else {
				pkg.Bin = path.Join(pkgFolder, "dist/main.js")
			}
		} else {
			if _, err := GetPackageEntryPoint(repo, pkg); err != nil && checkEntryPoints {
				return err
			} else {
				pkg.Main = path.Join(pkgFolder, "dist/index.js")
				pkg.Types = path.Join(pkgFolder, "dist/index.d.ts")
			}
		}
		repo.Packages[packageName] = pkg
		repo.PackageNames = append(repo.PackageNames, packageName)
	}
	return nil
}

func MakeJobs(jobs *gjobs.GJobs, jobPrefix string, repo *Repository, witchPackages []string, fn func(pkg *Package) error) error {

	var errWalk error
	errMutex := sync.Mutex{}
	for _, pkg := range repo.Packages {
		if Contains(witchPackages, pkg.Name) {
			dep := []string{}
			for depName := range pkg.Dependencies {
				dep = append(dep, jobPrefix+depName)
			}
			pkgInternal := pkg
			jobs.NewJob(jobPrefix+pkg.Name, dep, func() (interface{}, error) {
				if err := fn(pkgInternal); err != nil {
					errMutex.Lock()
					errWalk = err
					errMutex.Unlock()
				}
				return nil, nil
			})
		}
	}

	return errWalk
}

func validateLayer(repo *Repository, pkg *Package) error {

	if pkg.Layer == BusinessRulesLayer {
		if pkg.Executable {
			return errors.New("Business layer " + pkg.Name + " can't be an executble")
		}
		for depName := range pkg.Dependencies {
			var depPkg = repo.Packages[depName]
			if depPkg == nil {
				return errors.New("don't use external dependency " + depName + " on business layer " + pkg.Name)
			}
			if depPkg.Layer == AdaptersLayer {
				return errors.New("business layer " + pkg.Name + "can't depends of adapter layer " + depName)
			}
			if depPkg.Layer == ExecutablesLayer {
				return errors.New("business layer " + pkg.Name + "can't depends of executable layer " + depName)
			}
		}
	}
	if pkg.Layer == AdaptersLayer && pkg.Executable {
		return errors.New("Adapter layer " + pkg.Name + " can't be an executble")
	}

	return avoidCircularReferences(repo, repo.Packages, map[string]bool{})
}

const configName = "monoclean.yml"

var ErrNoConfig = fmt.Errorf("cannot find %s config file", configName)

func openConfigFile(searchDir string) (*os.File, error) {
	for {
		configPath := path.Join(searchDir, configName)
		f, err := os.Open(configPath)
		if os.IsNotExist(err) {
			searchDir = path.Dir(searchDir)
			if len(searchDir) <= 1 {
				return nil, ErrNoConfig
			}
			continue
		}
		println("workspace: " + configPath)
		return f, err
	}
}

func avoidCircularReferences(repo *Repository, packages map[string]*Package, used map[string]bool) error {
	for _, pkg := range packages {
		if err := avoidCircularReferencesForPackage(repo, pkg, used); err != nil {
			return err
		}
	}
	return nil
}

func avoidCircularReferencesForPackage(repo *Repository, pkg *Package, used map[string]bool) error {
	for dependencyName := range pkg.Dependencies {
		depPkg := repo.Packages[dependencyName]
		if depPkg.Executable {
			return errors.New(pkg.Name + " references to executable " + dependencyName)
		}
		var _, circularRef = used[dependencyName]
		if circularRef {
			return errors.New("Package " + dependencyName + " has circular reference")
		}
	}
	return nil
}
