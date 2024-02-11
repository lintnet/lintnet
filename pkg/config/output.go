package config

type Output struct {
	ID string `json:"id"`
	// text/template, html/template, jsonnet
	Renderer string `json:"renderer"`
	// path to a template file
	Template string `json:"template"`
	// parameter
	Config map[string]any `json:"config"`
	// Transform parameter
	Transform       string  `json:"transform"`
	TemplateModule  *Module `json:"-"`
	TransformModule *Module `json:"-"`
}
