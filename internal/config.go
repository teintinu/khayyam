package internal

// TODO: Replace all maps with slices, so that we can detect duplicate keys and
// preserve order if necessary.

type Config struct {
	Engines       map[string]string
	Repository    string
	Registry      string
	Workspace     WorkspaceConfig
	Packages      map[string]PackageConfig
	Domains       map[string]PackageConfig `yaml:"domains"` 
	Adapters      map[string]PackageConfig
	Applications  map[string]PackageConfig
}

type WorkspaceConfig struct {
	Name         string
	Version      string
	Dependencies map[string]string
}

type PackageConfig struct {
	Publish      string
	Description  string
	Version      string
	Index        string
	Folder       string
	Executable   bool
	Dependencies map[string]string
}
