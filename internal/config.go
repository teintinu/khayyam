package internal

// TODO: Replace all maps with slices, so that we can detect duplicate keys and
// preserve order if necessary.

type Config struct {
	Engines       map[string]string
	Repository    string
	Registry      string
	Workspace     WorkspaceConfig
	Dependencies  map[string]string
	Packages      map[string]PackageConfig
	BusinessRules map[string]PackageConfig `yaml:"business-rules"`
	Executables   map[string]PackageConfig
	Adapters      map[string]PackageConfig
}

type WorkspaceConfig struct {
	Name         string
	Version      string
	Dependencies map[string]string
}

type PackageConfig struct {
	Public       bool
	Description  string
	Version      string
	Index        string
	Folder       string
	Executables  map[string]string
	Dependencies map[string]string
}
