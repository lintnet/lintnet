package lint

import (
	"io"

	"github.com/lintnet/lintnet/pkg/config"
	"github.com/lintnet/lintnet/pkg/domain"
	"github.com/lintnet/lintnet/pkg/encoding"
	"github.com/lintnet/lintnet/pkg/filefind"
	"github.com/lintnet/lintnet/pkg/jsonnet"
	"github.com/lintnet/lintnet/pkg/lint"
	"github.com/lintnet/lintnet/pkg/lintfile"
	"github.com/lintnet/lintnet/pkg/module"
	"github.com/lintnet/lintnet/pkg/output"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

type Controller struct {
	fs              afero.Fs
	stdout          io.Writer
	moduleInstaller *module.Installer
	importer        *jsonnet.Importer
	param           *ParamController
	dataFileParser  lint.DataFileParser
	linter          Linter
	fileFinder      FileFinder
	configReader    *ConfigReader
	outputGetter    *output.Getter
}

type Linter interface {
	Lint(targets []*domain.Target) ([]*domain.Result, error)
}

type FileFinder interface {
	Find(logE *logrus.Entry, cfg *config.Config, rootDir, cfgDir string) ([]*domain.Target, error)
}

type ParamController struct {
	Version string
}

func NewController(param *ParamController, fs afero.Fs, stdout io.Writer, moduleInstaller *module.Installer, importer *jsonnet.Importer) *Controller {
	dp := encoding.NewDataFileParser(fs)
	return &Controller{
		param:           param,
		fs:              fs,
		stdout:          stdout,
		moduleInstaller: moduleInstaller,
		importer:        importer,
		linter: lint.NewLinter(
			dp,
			lintfile.NewParser(fs),
			lintfile.NewEvaluator(importer),
		),
		dataFileParser: dp,
		fileFinder:     filefind.NewFileFinder(fs),
		configReader: &ConfigReader{
			fs:       fs,
			importer: importer,
		},
		outputGetter: output.NewGetter(stdout, fs, importer),
	}
}
