package lint

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/google/go-jsonnet"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/lintnet/pkg/config"
	"github.com/suzuki-shunsuke/lintnet/pkg/render"
)

func (c *Controller) Output(logE *logrus.Entry, cfg *config.Config, logLevel ErrorLevel, results map[string]*FileResult) error {
	if !isFailed(results) {
		return nil
	}
	fes := c.formatResultToOutput(logLevel, results)
	if len(fes) == 0 {
		return nil
	}
	outputs := cfg.Outputs
	if len(cfg.Outputs) == 0 {
		outputs = []*config.Output{
			{
				Type:     "stdout",
				Renderer: "jsonnet",
			},
		}
	}
	for _, output := range outputs {
		if err := c.output(output, fes); err != nil {
			logE.WithError(err).Error("output errors")
		}
	}
	return errors.New("lint failed")
}

func (c *Controller) outputByJsonnet(output *config.Output, fes []*FlatError) error {
	out := c.stdout
	if output.Type == "file" {
		f, err := c.fs.Create(output.Path)
		if err != nil {
			return fmt.Errorf("create a file: %w", err)
		}
		defer f.Close()
		out = f
	}
	if output.Template != "" {
		node, err := c.readJsonnet(output.Template)
		if err != nil {
			return fmt.Errorf("read a template as Jsonnet: %w", err)
		}
		b, err := json.Marshal(fes)
		if err != nil {
			return fmt.Errorf("marshal results as JSON: %w", err)
		}
		vm := jsonnet.MakeVM()
		vm.ExtCode("input", string(b))
		setNativeFunctions(vm)
		result, err := vm.Evaluate(node)
		if err != nil {
			return fmt.Errorf("evaluate a Jsonnet: %w", err)
		}
		fmt.Fprintln(out, result)
		return nil
	}
	return c.outputJSON(out, fes)
}

func (c *Controller) outputByTemplate(output *config.Output, fes []*FlatError, renderer render.TemplateRenderer) error {
	out := c.stdout
	if output.Type == "file" {
		f, err := c.fs.Create(output.Path)
		if err != nil {
			return fmt.Errorf("create a file: %w", err)
		}
		defer f.Close()
		out = f
	}
	if output.Template != "" {
		b, err := afero.ReadFile(c.fs, output.Template)
		if err != nil {
			return fmt.Errorf("read a template: %w", err)
		}
		if err := renderer.Render(out, string(b), map[string]interface{}{
			"input": fes,
		}); err != nil {
			return fmt.Errorf("render a template: %w", err)
		}
		return nil
	}
	return nil
}

func (c *Controller) output(output *config.Output, fes []*FlatError) error {
	switch output.Renderer {
	case "jsonnet":
		return c.outputByJsonnet(output, fes)
	case "text/template":
		return c.outputByTemplate(output, fes, &render.TextTemplateRenderer{})
	case "html/template":
		return c.outputByTemplate(output, fes, &render.HTMLTemplateRenderer{})
	}
	return errors.New("unknown renderer")
}

type FlatError struct {
	RuleName     string      `json:"rule,omitempty"`
	Level        string      `json:"level,omitempty"`
	Message      string      `json:"message,omitempty"`
	LintFilePath string      `json:"lint_file,omitempty"`
	DataFilePath string      `json:"data_file,omitempty"`
	Location     interface{} `json:"location,omitempty"`
}

func (c *Controller) formatResultToOutput(logLevel ErrorLevel, results map[string]*FileResult) []*FlatError {
	list := make([]*FlatError, 0, len(results))
	for dataFilePath, fileResult := range results {
		list = append(list, fileResult.flattenError(logLevel, dataFilePath)...)
	}
	if len(list) == 0 {
		return nil
	}
	return list
}

func (c *Controller) outputJSON(w io.Writer, result interface{}) error {
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(result); err != nil {
		return fmt.Errorf("encode the result as JSON: %w", err)
	}
	return nil
}
