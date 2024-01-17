package output

import (
	"encoding/json"
	"fmt"
	"io"

	gojsonnet "github.com/google/go-jsonnet"
	"github.com/lintnet/lintnet/pkg/config"
	"github.com/lintnet/lintnet/pkg/jsonnet"
	"github.com/spf13/afero"
)

type jsonnetOutputter struct {
	stdout    io.Writer
	output    *config.Output
	transform jsonnet.Node
	node      jsonnet.Node
	importer  gojsonnet.Importer
	config    map[string]any
}

func newJsonnetOutputter(fs afero.Fs, stdout io.Writer, output *config.Output, importer gojsonnet.Importer) (*jsonnetOutputter, error) {
	node, err := jsonnet.ReadToNode(fs, output.Template)
	if err != nil {
		return nil, fmt.Errorf("read a template as Jsonnet: %w", err)
	}
	outputter := &jsonnetOutputter{
		stdout:   stdout,
		output:   output,
		node:     node,
		importer: importer,
		config:   output.Config,
	}
	if output.Transform != "" {
		node, err := jsonnet.ReadToNode(fs, output.Transform)
		if err != nil {
			return nil, fmt.Errorf("read a transform as Jsonnet: %w", err)
		}
		outputter.transform = node
	}
	return outputter, nil
}

func (o *jsonnetOutputter) Output(result *Output) error {
	r := *result
	r.Config = o.config
	tla, err := json.Marshal(&r)
	if err != nil {
		return fmt.Errorf("marshal output as JSON: %w", err)
	}
	tlaS := string(tla)
	if o.transform != nil {
		vm := jsonnet.NewVM(string(tla), o.importer)
		s, err := vm.Evaluate(o.transform)
		if err != nil {
			return fmt.Errorf("evaluate a jsonnet: %w", err)
		}
		tlaS = s
	}
	vm := jsonnet.NewVM(tlaS, o.importer)
	s, err := vm.Evaluate(o.node)
	if err != nil {
		return fmt.Errorf("evaluate a jsonnet: %w", err)
	}
	var a any
	if err := json.Unmarshal([]byte(s), &a); err != nil {
		return fmt.Errorf("unmarshal the result as JSON: %w", err)
	}
	return outputJSON(o.stdout, a)
}
