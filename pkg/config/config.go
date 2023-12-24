package config

type Config struct {
	ErrorLevel string    `json:"error_level" yaml:"error_level"`
	Targets    []*Target `json:"targets"`
	Outputs    []*Output `json:"outputs"`
}

type Output struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Renderer string `json:"renderer"`
	Path     string `json:"path"`
	Template string `json:"template"`
}

type Target struct {
	LintFiles []string `json:"lint_files" yaml:"lint_files"`
	Modules   []string `json:"modules"`
	DataFiles []string `json:"data_files" yaml:"data_files"`
}
