package output

import (
	"errors"
	"io"
	"path/filepath"

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

type Outputter interface {
	// Outputter outputs the result.
	// Outputter compromises of transformer and renderer.
	Output(result *Output) error
}

type ParamGet struct {
	RootDir string
	Output  string
}

// setTransform set output.Transform.
func setTransform(output *config.Output, param *ParamGet, cfgDir string) {
	if output.TransformModule != nil {
		output.Transform = filepath.Join(param.RootDir, output.TransformModule.FilePath())
		return
	}
	output.Transform = filepath.FromSlash(output.Transform)
	if !filepath.IsAbs(output.Transform) {
		output.Transform = filepath.Join(cfgDir, output.Transform)
	}
}

func setTemplate(output *config.Output, param *ParamGet, cfgDir string) {
	if output.TemplateModule != nil {
		output.Template = filepath.Join(param.RootDir, output.TemplateModule.FilePath())
		return
	}
	output.Template = filepath.FromSlash(output.Template)
	if !filepath.IsAbs(output.Template) {
		output.Template = filepath.Join(cfgDir, output.Template)
	}
}

// Get returns an outputter.
func (g *Getter) Get(outputs config.Outputs, param *ParamGet, cfgDir string) (Outputter, error) {
	if param.Output == "" {
		return &jsonOutputter{
			stdout: g.stdout,
		}, nil
	}
	output := outputs.Output(param.Output)
	if output == nil {
		return nil, errors.New("unknown output id")
	}

	if output.Template != "" {
		setTemplate(output, param, cfgDir)
	}

	if output.Transform != "" {
		setTransform(output, param, cfgDir)
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
