package config

import (
	"fmt"
	"strings"
)

type Outputs []*Output

func (os Outputs) Output(id string) *Output {
	for _, o := range os {
		if o.ID == id {
			return o
		}
	}
	return nil
}

func (os Outputs) Preprocess(modules map[string]*ModuleArchive) error {
	for _, output := range os {
		if err := output.Preprocess(modules); err != nil {
			return err
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
	// A transformation file must be a Jsonnet.
	// A file path must be an absolute path, a relative path from the configuration file, or a module path.
	// e.g.
	// transform.jsonnnet
	// /home/foo/.lintent/transform.jsonnnet
	// github_archive/github.com/lintnet/modules/transform.jsonnet@32ca3be646ec5b5861aab72fed30cd71f6eba9bf:v0.1.2
	Transform string `json:"transform"`

	TemplateModule  *Module `json:"-"`
	TransformModule *Module `json:"-"`
}

func (o *Output) Preprocess(modules map[string]*ModuleArchive) error {
	if strings.HasPrefix(o.Template, "github_archive/github.com/") {
		m, err := ParseImport(o.Template)
		if err != nil {
			return fmt.Errorf("parse a module path: %w", err)
		}
		o.TemplateModule = m
		modules[m.Archive.String()] = m.Archive
	}
	if strings.HasPrefix(o.Transform, "github_archive/github.com/") {
		m, err := ParseImport(o.Transform)
		if err != nil {
			return fmt.Errorf("parse a module path: %w", err)
		}
		o.TransformModule = m
		modules[m.Archive.String()] = m.Archive
	}
	return nil
}
