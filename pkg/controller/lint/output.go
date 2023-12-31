package lint

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"runtime"

	"github.com/lintnet/lintnet/pkg/config"
	"github.com/lintnet/lintnet/pkg/errlevel"
	"github.com/lintnet/lintnet/pkg/jsonnet"
	"github.com/lintnet/lintnet/pkg/render"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

func (c *Controller) Output(logE *logrus.Entry, errLevel errlevel.Level, results []*Result, outputters []Outputter, outputSuccess bool) error {
	fes := c.formatResultToOutput(results)
	failed, err := isFailed(fes.Errors, errLevel)
	if err != nil {
		return err
	}
	if !outputSuccess && !failed {
		return nil
	}
	for _, outputter := range outputters {
		if err := outputter.Output(fes); err != nil {
			logE.WithError(err).Error("output errors")
		}
	}
	if failed {
		return errors.New("lint failed")
	}
	return nil
}

type jsonOutputter struct {
	stdout io.Writer
}

func (o *jsonOutputter) Output(result *Output) error {
	r := *result
	for i, e := range result.Errors {
		fe := *e
		fe.Description = ""
		r.Errors[i] = &fe
	}
	return outputJSON(o.stdout, &r)
}

type jsonnetOutputter struct {
	stdout    io.Writer
	output    *config.Output
	transform jsonnet.Node
	node      jsonnet.Node
	importer  *jsonnet.Importer
	config    map[string]any
}

func newJsonnetOutputter(fs afero.Fs, stdout io.Writer, output *config.Output, importer *jsonnet.Importer) (*jsonnetOutputter, error) {
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

func outputJSON(w io.Writer, result any) error {
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(result); err != nil {
		return fmt.Errorf("encode the result as JSON: %w", err)
	}
	return nil
}

type Outputter interface {
	Output(result *Output) error
}

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

func (c *Controller) getOutput(outputs []*config.Output, outputID string) (*config.Output, error) {
	for _, output := range outputs {
		if output.ID == outputID {
			return output, nil
		}
	}
	return nil, errors.New("unknown output id")
}

func (c *Controller) getOutputter(outputs []*config.Output, outputID, rootDir string) (Outputter, error) {
	if outputID == "" {
		return &jsonOutputter{
			stdout: c.stdout,
		}, nil
	}
	output, err := c.getOutput(outputs, outputID)
	if err != nil {
		return nil, err
	}

	if output.TemplateModule != nil {
		output.Template = filepath.Join(rootDir, output.TemplateModule.FilePath())
	} else {
		output.Template = filepath.FromSlash(output.Template)
	}

	if output.TransformModule != nil {
		output.Transform = filepath.Join(rootDir, output.TransformModule.FilePath())
	} else {
		output.Transform = filepath.FromSlash(output.Transform)
	}

	switch output.Renderer {
	case "jsonnet":
		return newJsonnetOutputter(c.fs, c.stdout, output, c.importer)
	case "text/template":
		return newTemplateOutputter(c.stdout, c.fs, &render.TextTemplateRenderer{}, output, c.importer)
	case "html/template":
		return newTemplateOutputter(c.stdout, c.fs, &render.HTMLTemplateRenderer{}, output, c.importer)
	}
	return nil, errors.New("unknown renderer")
}

type FlatError struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Level       string `json:"level,omitempty"`
	Message     string `json:"message,omitempty"`
	LintFile    string `json:"lint_file,omitempty"`
	DataFile    string `json:"data_file,omitempty"`
	// DataFilePaths []string `json:"data_files,omitempty"`
	TargetID string `json:"target_id,omitempty"`
	Location any    `json:"location,omitempty"`
	Custom   any    `json:"custom,omitempty"`
}

func (e *FlatError) Failed(errLevel errlevel.Level) (bool, error) {
	level := errlevel.Error
	if e.Level != "" {
		feErrLevel, err := errlevel.New(e.Level)
		if err != nil {
			return false, fmt.Errorf("verify the error level of a result: %w", err)
		}
		level = feErrLevel
	}
	return level >= errLevel, nil
}

type Output struct {
	LintnetVersion string         `json:"lintnet_version"`
	Env            string         `json:"env"`
	Errors         []*FlatError   `json:"errors,omitempty"`
	Config         map[string]any `json:"config,omitempty"`
}

func (c *Controller) formatResultToOutput(results []*Result) *Output {
	list := make([]*FlatError, 0, len(results))
	for _, result := range results {
		list = append(list, result.FlatErrors()...)
	}
	return &Output{
		LintnetVersion: c.param.Version,
		Env:            fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		Errors:         list,
	}
}
