package lint

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/lintnet/lintnet/pkg/config"
	"github.com/lintnet/lintnet/pkg/errlevel"
	"github.com/lintnet/lintnet/pkg/render"
	"github.com/sirupsen/logrus"
)

type Outputter interface {
	Output(result *Output) error
}

func (c *Controller) Output(logE *logrus.Entry, errLevel, shownErrLevel errlevel.Level, results []*Result, outputters []Outputter, outputSuccess bool) error {
	fes := c.formatResultToOutput(results, shownErrLevel)
	failed, err := isFailed(fes.Errors, errLevel)
	if err != nil {
		return err
	}
	if !outputSuccess && len(fes.Errors) == 0 {
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

func outputJSON(w io.Writer, result any) error {
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(result); err != nil {
		return fmt.Errorf("encode the result as JSON: %w", err)
	}
	return nil
}

func getOutput(outputs []*config.Output, outputID string) (*config.Output, error) {
	for _, output := range outputs {
		if output.ID == outputID {
			return output, nil
		}
	}
	return nil, errors.New("unknown output id")
}

func (c *Controller) getOutputter(outputs []*config.Output, param *ParamLint, cfgDir string) (Outputter, error) { //nolint:cyclop
	if param.Output == "" {
		return &jsonOutputter{
			stdout: c.stdout,
		}, nil
	}
	output, err := getOutput(outputs, param.Output)
	if err != nil {
		return nil, err
	}

	if output.TemplateModule != nil {
		output.Template = filepath.Join(param.RootDir, output.TemplateModule.FilePath())
	} else {
		output.Template = filepath.FromSlash(output.Template)
		if !filepath.IsAbs(output.Template) {
			output.Template = filepath.Join(cfgDir, output.Template)
		}
		a, err := filepath.Rel(param.DataRootDir, output.Template)
		if err != nil {
			return nil, fmt.Errorf("get a relative path to template: %w", err)
		}
		if strings.HasPrefix(a, "..") {
			return nil, errors.New("this template is unavailable because the template is out of data root directory")
		}
	}

	if output.TransformModule != nil {
		output.Transform = filepath.Join(param.RootDir, output.TransformModule.FilePath())
	} else {
		output.Transform = filepath.FromSlash(output.Transform)
		if !filepath.IsAbs(output.Transform) {
			output.Transform = filepath.Join(cfgDir, output.Transform)
		}
		a, err := filepath.Rel(param.DataRootDir, output.Transform)
		if err != nil {
			return nil, fmt.Errorf("get a relative path to transform: %w", err)
		}
		if strings.HasPrefix(a, "..") {
			return nil, errors.New("this transform is unavailable because the transform is out of data root directory")
		}
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
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	Links       []*Link `json:"links,omitempty"`
	Level       string  `json:"level,omitempty"`
	Message     string  `json:"message,omitempty"`
	LintFile    string  `json:"lint_file,omitempty"`
	DataFile    string  `json:"data_file,omitempty"`
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

func (c *Controller) formatResultToOutput(results []*Result, errLevel errlevel.Level) *Output {
	list := make([]*FlatError, 0, len(results))
	for _, result := range results {
		for _, fe := range result.FlatErrors() {
			el, err := errlevel.New(fe.Level) // TODO output warning
			if err != nil || el >= errLevel {
				list = append(list, fe)
			}
		}
	}
	return &Output{
		LintnetVersion: c.param.Version,
		Env:            fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		Errors:         list,
	}
}
