package config

type Config struct {
	ErrorLevel string `yaml:"error_level"`
	Modules    []*Module
	Targets    []*Target
	Outputs    []*Output
}

type Output struct {
	Type     string
	Renderer string
	Path     string
	Template string
}

type Module struct {
	ID     string
	Source string
}

type Target struct {
	LintFiles *LintFiles `yaml:"lint_files"`
	DataFiles *DataFiles `yaml:"data_files"`
}

type LintFiles struct {
	SearchType string `yaml:"search_type"`
	Paths      []*Path
	Imports    []*Import
}

type DataFiles struct {
	SearchType string `yaml:"search_type"`
	Paths      []*Path
}

type Path struct {
	Path    string
	Exclude bool
}

type Import struct {
	Module string
	Path   string
	Import string
}
