package internal

import (
	"encoding/json"
	"io/ioutil"
	"path"
)

type TsConfigMetadata struct {
	Extends         string                         `json:"extends,omitempty"`
	CompilerOptions TsConfigCompileOptionsMetadata `json:"compilerOptions,omitempty"`
	Exclude         []string                       `json:"exclude,omitempty"`
	Include         []string                       `json:"include,omitempty"`
	References      []TsConfigReferenceMetadata    `json:"references,omitempty"`
}

type TsConfigCompileOptionsMetadata struct {
	Incremental      bool                `json:"incremental,omitempty"`
	Declaration      bool                `json:"declaration,omitempty"`
	SourceMap        bool                `json:"sourceMap,omitempty"`
	Composite        bool                `json:"composite,omitempty"`
	ImportHelpers    bool                `json:"importHelpers,omitempty"`
	Strict           bool                `json:"strict,omitempty"`
	EsModuleInterop  bool                `json:"esModuleInterop,omitempty"`
	Target           string              `json:"target,omitempty"`
	ModuleResolution string              `json:"moduleResolution,omitempty"`
	Module           string              `json:"module,omitempty"`
	Lib              []string            `json:"lib,omitempty"`
	RootDir          string              `json:"rootDir,omitempty"`
	BaseURL          string              `json:"baseUrl,omitempty"`
	Paths            map[string][]string `json:"paths,omitempty"`
	OutDir           string              `json:"outDir,omitempty"`
}

type TsConfigReferenceMetadata struct {
	Path string `json:"path,omitempty"`
}

func WriteTsConfigJSON(metadata TsConfigMetadata, filename string) error {
	return WriteJSON(filename, metadata)
}

func ReadTConfigJSON(filename string) (*TsConfigMetadata, error) {
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var metadata TsConfigMetadata
	if err := json.Unmarshal(bs, &metadata); err != nil {
		return nil, err
	}
	return &metadata, nil
}

func GetPackageEntryPoint(repo *Repository, pkg *Package) (string, error) {

	pkgRoot := path.Join(repo.RootDir, pkg.Folder)
	srcDir := path.Join(pkgRoot, "src")

	if pkg.Layer == ExecutablesLayer {
		return NeedSomeOfTheseFiles(
			srcDir,
			[]string{"main.ts", "main.tsx"},
			"main.ts or main.tsx not found in package "+pkg.Name,
		)
	}

	return NeedSomeOfTheseFiles(
		srcDir,
		[]string{"index.ts", "index.tsx"},
		"main.ts or main.tsx not found in package "+pkg.Name,
	)
}
