package lint

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/lintnet/lintnet/pkg/config"
	"github.com/lintnet/lintnet/pkg/errlevel"
	"github.com/lintnet/lintnet/pkg/jsonnet"
	"github.com/lintnet/lintnet/pkg/render"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

func (c *Controller) Output(logE *logrus.Entry, cfg *config.Config, errLevel errlevel.Level, results map[string]*FileResult, outputIDs []string, outputSuccess bool) error {
	fes := c.formatResultToOutput(results)
	if !outputSuccess && len(fes.Errors) == 0 {
		return nil
	}
	failed, err := isFailed(fes.Errors, errLevel)
	if err != nil {
		return err
	}
	if !outputSuccess && !failed {
		return nil
	}
	outputs, err := c.getOutputs(cfg, outputIDs)
	if err != nil {
		return err
	}
	for _, output := range outputs {
		if err := c.output(output, fes); err != nil {
			logE.WithError(err).Error("output errors")
		}
	}
	if failed {
		return errors.New("lint failed")
	}
	return nil
}

func (c *Controller) getOutputs(cfg *config.Config, outputIDs []string) ([]*config.Output, error) {
	outputList := cfg.Outputs
	if len(outputList) == 0 {
		outputList = []*config.Output{
			{
				Type:     "stdout",
				Renderer: "jsonnet",
			},
		}
	}
	if len(outputIDs) == 0 {
		outputIDs = []string{
			"stdout",
		}
	}
	outputs := make([]*config.Output, len(outputIDs))
	outputMap := make(map[string]*config.Output, len(outputList))
	for _, output := range outputList {
		outputMap[output.ID] = output
	}
	for i, outputID := range outputIDs {
		output, ok := outputMap[outputID]
		if !ok {
			return nil, errors.New("unknown output id")
		}
		outputs[i] = output
	}
	return outputs, nil
}

func (c *Controller) outputByJsonnet(output *config.Output, result *Output) error {
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
		node, err := jsonnet.Read(c.fs, output.Template)
		if err != nil {
			return fmt.Errorf("read a template as Jsonnet: %w", err)
		}
		b, err := json.Marshal(result)
		if err != nil {
			return fmt.Errorf("marshal results as JSON: %w", err)
		}
		vm := jsonnet.MakeVM()
		vm.ExtCode("input", string(b))
		jsonnet.SetNativeFunctions(vm)
		result, err := vm.Evaluate(node)
		if err != nil {
			return fmt.Errorf("evaluate a Jsonnet: %w", err)
		}
		fmt.Fprintln(out, result)
		return nil
	}
	return c.outputJSON(out, result)
}

func (c *Controller) outputByTemplate(output *config.Output, result *Output, renderer render.TemplateRenderer) error {
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
			"result": result,
		}); err != nil {
			return fmt.Errorf("render a template: %w", err)
		}
		return nil
	}
	return nil
}

func (c *Controller) output(output *config.Output, out *Output) error {
	switch output.Renderer {
	case "jsonnet":
		return c.outputByJsonnet(output, out)
	case "text/template":
		return c.outputByTemplate(output, out, &render.TextTemplateRenderer{})
	case "html/template":
		return c.outputByTemplate(output, out, &render.HTMLTemplateRenderer{})
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

type Output struct {
	Errors []*FlatError `json:"errors,omitempty"`
}

func (c *Controller) formatResultToOutput(results map[string]*FileResult) *Output {
	list := make([]*FlatError, 0, len(results))
	for dataFilePath, fileResult := range results {
		list = append(list, fileResult.flattenError(dataFilePath)...)
	}
	return &Output{
		Errors: list,
	}
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
