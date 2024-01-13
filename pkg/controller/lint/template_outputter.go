package lint

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/lintnet/lintnet/pkg/config"
	"github.com/lintnet/lintnet/pkg/jsonnet"
	"github.com/lintnet/lintnet/pkg/render"
	"github.com/spf13/afero"
)

type templateOutputter struct {
	stdout   io.Writer
	fs       afero.Fs
	output   *config.Output
	template render.Template
	node     jsonnet.Node
	importer *jsonnet.Importer
}

func newTemplateOutputter(stdout io.Writer, fs afero.Fs, renderer render.TemplateRenderer, output *config.Output, importer *jsonnet.Importer) (*templateOutputter, error) {
	if output.Template == "" {
		return nil, errors.New("template is required")
	}
	b, err := afero.ReadFile(fs, output.Template)
	if err != nil {
		return nil, fmt.Errorf("read a template: %w", err)
	}
	tpl, err := renderer.Compile(string(b))
	if err != nil {
		return nil, fmt.Errorf("parse a template: %w", err)
	}
	var node jsonnet.Node
	if output.Transform != "" {
		n, err := jsonnet.ReadToNode(fs, output.Transform)
		if err != nil {
			return nil, fmt.Errorf("read a transform as Jsonnet: %w", err)
		}
		node = n
	}
	return &templateOutputter{
		stdout:   stdout,
		fs:       fs,
		output:   output,
		template: tpl,
		importer: importer,
		node:     node,
	}, nil
}

func (o *templateOutputter) Output(result *Output) error {
	r := *result
	r.Config = o.output.Config
	b, err := json.Marshal(r)
	if err != nil {
		return fmt.Errorf("marshal result as JSON: %w", err)
	}
	var param any
	if o.node != nil {
		vm := jsonnet.NewVM(string(b), o.importer)
		s, err := vm.Evaluate(o.node)
		if err != nil {
			return fmt.Errorf("evaluate a jsonnet: %w", err)
		}
		if err := json.Unmarshal([]byte(s), &param); err != nil {
			return fmt.Errorf("unmarshal transformed result as JSON: %w", err)
		}
	} else {
		if err := json.Unmarshal(b, &param); err != nil {
			return fmt.Errorf("unmarshal result as JSON: %w", err)
		}
	}
	if err := o.template.Execute(o.stdout, param); err != nil {
		return fmt.Errorf("render a template: %w", err)
	}
	return nil
}
