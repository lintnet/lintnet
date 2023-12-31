package lint

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
	return outputJSON(o.stdout, result)
}

type jsonnetOutputter struct {
	stdout   io.Writer
	output   *config.Output
	node     jsonnet.Node
	importer *jsonnet.Importer
}

func newJsonnetOutputter(fs afero.Fs, stdout io.Writer, output *config.Output, importer *jsonnet.Importer) (*jsonnetOutputter, error) {
	node, err := jsonnet.ReadToNode(fs, output.Template)
	if err != nil {
		return nil, fmt.Errorf("read a template as Jsonnet: %w", err)
	}
	return &jsonnetOutputter{
		stdout:   stdout,
		output:   output,
		node:     node,
		importer: importer,
	}, nil
}

func (o *jsonnetOutputter) Output(result *Output) error {
	tla, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("marshal output as JSON: %w", err)
	}
	vm := jsonnet.NewVM(string(tla), o.importer)
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
}

func newTemplateOutputter(stdout io.Writer, fs afero.Fs, renderer render.TemplateRenderer, output *config.Output) (*templateOutputter, error) {
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
	return &templateOutputter{
		stdout:   stdout,
		fs:       fs,
		output:   output,
		template: tpl,
	}, nil
}

func (o *templateOutputter) Output(result *Output) error {
	if err := o.template.Execute(o.stdout, map[string]any{
		"result": result,
	}); err != nil {
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

func (c *Controller) getOutputter(outputs []*config.Output, outputID string) (Outputter, error) {
	if outputID == "" {
		return &jsonOutputter{
			stdout: c.stdout,
		}, nil
	}
	output, err := c.getOutput(outputs, outputID)
	if err != nil {
		return nil, err
	}
	switch output.Renderer {
	case "jsonnet":
		return newJsonnetOutputter(c.fs, c.stdout, output, c.importer)
	case "text/template":
		return newTemplateOutputter(c.stdout, c.fs, &render.TextTemplateRenderer{}, output)
	case "html/template":
		return newTemplateOutputter(c.stdout, c.fs, &render.HTMLTemplateRenderer{}, output)
	}
	return nil, errors.New("unknown renderer")
}

type FlatError struct {
	RuleName     string `json:"rule,omitempty"`
	Level        string `json:"level,omitempty"`
	Message      string `json:"message,omitempty"`
	LintFilePath string `json:"lint_file,omitempty"`
	DataFilePath string `json:"data_file,omitempty"`
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
	LintnetVersion string       `json:"lintnet_version"`
	Env            string       `json:"env"`
	Errors         []*FlatError `json:"errors,omitempty"`
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
