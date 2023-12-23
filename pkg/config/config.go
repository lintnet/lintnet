package config

type Config struct {
	ErrorLevel string `yaml:"error_level"`
	// Modules    []*Module
	Targets []*Target
	Outputs []*Output
}

type Output struct {
	ID       string
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
	LintFiles string   `yaml:"lint_files"`
	Modules   []string `yaml:"modules"`
	DataFiles string   `yaml:"data_files"`
}

type Import struct {
	Module string
	Path   string
	Import string
}
