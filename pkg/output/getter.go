package output

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/lintnet/lintnet/pkg/config"
	"github.com/lintnet/lintnet/pkg/jsonnet"
	"github.com/lintnet/lintnet/pkg/render"
	"github.com/spf13/afero"
)

type Getter struct {
	stdout   io.Writer
	fs       afero.Fs
	importer *jsonnet.Importer
}

func NewGetter(stdout io.Writer, fs afero.Fs, importer *jsonnet.Importer) *Getter {
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

func (g *Getter) Get(outputs []*config.Output, param *ParamGet, cfgDir string) (Outputter, error) { //nolint:cyclop
	if param.Output == "" {
		return &jsonOutputter{
			stdout: g.stdout,
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
		return newJsonnetOutputter(g.fs, g.stdout, output, g.importer)
	case "text/template":
		return newTemplateOutputter(g.stdout, g.fs, &render.TextTemplateRenderer{}, output, g.importer)
	case "html/template":
		return newTemplateOutputter(g.stdout, g.fs, &render.HTMLTemplateRenderer{}, output, g.importer)
	}
	return nil, errors.New("unknown renderer")
}
