package output

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	gojsonnet "github.com/google/go-jsonnet"
	"github.com/lintnet/lintnet/pkg/config"
	"github.com/lintnet/lintnet/pkg/render"
	"github.com/spf13/afero"
)

type Getter struct {
	stdout   io.Writer
	fs       afero.Fs
	importer gojsonnet.Importer
}

func NewGetter(stdout io.Writer, fs afero.Fs, importer gojsonnet.Importer) *Getter {
	return &Getter{
		stdout:   stdout,
		fs:       fs,
		importer: importer,
	}
}

func getOutput(outputs []*config.Output, outputID string) (*config.Output, error) {
	for _, output := range outputs {
		if output.ID == outputID {
			return output, nil
		}
	}
	return nil, errors.New("unknown output id")
}

type Outputter interface {
	Output(result *Output) error
}

type ParamGet struct {
	RootDir     string
	DataRootDir string
	Output      string
}

func setTransform(output *config.Output, param *ParamGet, cfgDir string) error {
	if output.TransformModule != nil {
		output.Transform = filepath.Join(param.RootDir, output.TransformModule.FilePath())
		return nil
	}
	output.Transform = filepath.FromSlash(output.Transform)
	if !filepath.IsAbs(output.Transform) {
		output.Transform = filepath.Join(cfgDir, output.Transform)
	}
	a, err := filepath.Rel(param.DataRootDir, output.Transform)
	if err != nil {
		return fmt.Errorf("get a relative path to transform: %w", err)
	}
	if strings.HasPrefix(a, "..") {
		return errors.New("this transform is unavailable because the transform is out of data root directory")
	}
	return nil
}

func setTemplate(output *config.Output, param *ParamGet, cfgDir string) error {
	if output.TemplateModule != nil {
		output.Template = filepath.Join(param.RootDir, output.TemplateModule.FilePath())
		return nil
	}
	output.Template = filepath.FromSlash(output.Template)
	if !filepath.IsAbs(output.Template) {
		output.Template = filepath.Join(cfgDir, output.Template)
	}
	a, err := filepath.Rel(param.DataRootDir, output.Template)
	if err != nil {
		return fmt.Errorf("get a relative path to template: %w", err)
	}
	if strings.HasPrefix(a, "..") {
		return errors.New("this template is unavailable because the template is out of data root directory")
	}
	return nil
}

func (g *Getter) Get(outputs []*config.Output, param *ParamGet, cfgDir string) (Outputter, error) {
	if param.Output == "" {
		return &jsonOutputter{
			stdout: g.stdout,
		}, nil
	}
	output, err := getOutput(outputs, param.Output)
	if err != nil {
		return nil, err
	}

	if output.Template != "" {
		if err := setTemplate(output, param, cfgDir); err != nil {
			return nil, err
		}
	}

	if output.Transform != "" {
		if err := setTransform(output, param, cfgDir); err != nil {
			return nil, err
		}
	}

	switch output.Renderer {
	case "jsonnet":
		return newJsonnetOutputter(g.fs, g.stdout, output, g.importer)
	case "text/template":
		return newTemplateOutputter(g.stdout, g.fs, &render.TextTemplateRenderer{}, output, g.importer)
	case "html/template":
		return newTemplateOutputter(g.stdout, g.fs, &render.HTMLTemplateRenderer{}, output, g.importer)
	}
	return nil, errors.New("unknown renderer")
}
