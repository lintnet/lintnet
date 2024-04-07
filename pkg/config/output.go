package config

type Outputs []*Output

func (os Outputs) Output(id string) *Output {
	for _, o := range os {
		if o.ID == id {
			return o
		}
	}
	return nil
}

type Output struct {
	ID string `json:"id"`
	// text/template, html/template, jsonnet
	Renderer string `json:"renderer"`
	// path to a template file
	Template string `json:"template"`
	// parameter
	Config map[string]any `json:"config"`
	// Transform is a transformation file path.
	// A transformation file transforms lint results before the results are outputted.
	// A tranformation file must be a Jsonnet.
	// A file path must be an absolute path, a relative path from the configuration file, or a module path.
	// e.g.
	// transform.jsonnnet
	// /home/foo/.lintent/transform.jsonnnet
	// github_archive/github.com/lintnet/modules/transform.jsonnet@32ca3be646ec5b5861aab72fed30cd71f6eba9bf:v0.1.2
	Transform string `json:"transform"`

	TemplateModule  *Module `json:"-"`
	TransformModule *Module `json:"-"`
}
