package internal

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/goccy/go-yaml"
)

type Repository struct {
	ConfigPath   string
	RootDir      string
	OutDir       string
	DistDir      string
	TmpDir       string
	Engines      map[string]string
	IsWorkspace  bool
	Workspace    Workspace
	Dependencies map[string]*Dependency
	Url          string
	Registry     string

	Packages map[string]*Package
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

type Package struct {
	Name          string
	Version       string
	Folder        string
	Public        bool
	usesNode      bool
	usesDOM       bool
	usesWebWorker bool
	Description   string
	Index         string
	Layer         PackageLayer
	Dependencies  map[string]*Dependency
	IsExecutable  bool
}

type Executable struct {
	Name       string
	Entrypoint string
}

const DefaultRegistry = "https://registry.npmjs.org/"

func LoadRepository(searchDir string) (*Repository, error) {
	f, err := openConfigFile(searchDir)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var repo Repository
	repo.ConfigPath = f.Name()
	repo.RootDir = path.Dir(repo.ConfigPath)
	repo.OutDir = path.Join(repo.RootDir, "out")
	repo.DistDir = path.Join(repo.OutDir, "dist")
	repo.TmpDir = path.Join(repo.OutDir, "tmp")

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

	repo.IsWorkspace = len(cfg.Workspace.Version) > 0

	if repo.IsWorkspace {
		repo.Workspace.Name = cfg.Workspace.Name
		repo.Workspace.Version = cfg.Workspace.Version
		if cfg.Dependencies != nil {
			return nil, errors.New("use dependencies inside workspace")
		}
	}

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
	err = loadPackages(repo, cfg.Packages, NormalLayer)
	if err != nil {
		return nil, err
	}
	err = loadPackages(repo, cfg.BusinessRules, BusinessRulesLayer)
	if err != nil {
		return nil, err
	}
	err = loadPackages(repo, cfg.Executables, ExecutablesLayer)
	if err != nil {
		return nil, err
	}
	err = loadPackages(repo, cfg.Adapters, AdaptersLayer)
	if err != nil {
		return nil, err
	}

	repo.Dependencies = make(map[string]*Dependency)
	addDependency := func(name, version string) {
		repo.Dependencies[name] = &Dependency{
			Name:    name,
			Version: version,
		}
	}
	for dependencyName, dependencyVersion := range cfg.Dependencies {
		addDependency(dependencyName, dependencyVersion)
	}
	var requiredDeps map[string]string
	if repo.IsWorkspace {
		requiredDeps = requiredDependenciesWorkspace
	} else {
		requiredDeps = requiredDependenciesMix
	}
	for dependencyName, dependencyVersion := range requiredDeps {
		if _, ok := repo.Dependencies[dependencyName]; !ok {
			addDependency(dependencyName, dependencyVersion)
		}
	}
	for pkgName := range repo.Packages {
		err = validateLayer(repo, repo.Packages[pkgName])
		if err != nil {
			return nil, err
		}
	}
	return &repo, nil
}

func loadPackages(repo Repository, packages map[string]PackageConfig, layer PackageLayer) error {
	for packageName, packageConfig := range packages {
		if repo.IsWorkspace {
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
		}
		if packageConfig.Dependencies != nil && (!repo.IsWorkspace) {
			return errors.New("package dependencies is supported only inside workspace")
		}
		pkg := &Package{
			Name:        packageName,
			Public:      packageConfig.Public,
			Description: packageConfig.Description,
			Index:       packageConfig.Index,
			Folder:      packageConfig.Folder,
			Layer:       layer,
		}
		if repo.IsWorkspace && packageConfig.Dependencies != nil {
			pkg.Dependencies = make(map[string]*Dependency)
			for dependencyName, dependencyVersion := range packageConfig.Dependencies {
				pkg.Dependencies[dependencyName] = &Dependency{
					Name:    dependencyName,
					Version: dependencyVersion,
				}
				pkg.usesDOM = dependencyName == "react-dom"
				pkg.usesNode = dependencyName == "@types/node"
			}
		}
		repo.Packages[packageName] = pkg
	}
	return nil
}

func validateLayer(repo Repository, pkg *Package) error {
	if pkg.Layer == NormalLayer {
		return nil
	}
	if pkg.Layer == BusinessRulesLayer {
		if pkg.IsExecutable {
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
	if pkg.Layer == AdaptersLayer && pkg.IsExecutable {
		return errors.New("Adapter layer " + pkg.Name + " can't be an executble")
	}
	if pkg.Layer == ExecutablesLayer && pkg.IsExecutable {
		return errors.New("Adapter layer " + pkg.Name + " can't be an executble")
	}
	return nil
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
		return f, err
	}
}
