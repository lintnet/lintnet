package lint

import (
	"io"

	"github.com/lintnet/lintnet/pkg/config"
	"github.com/lintnet/lintnet/pkg/domain"
	"github.com/lintnet/lintnet/pkg/filefind"
	"github.com/lintnet/lintnet/pkg/jsonnet"
	"github.com/lintnet/lintnet/pkg/lint"
	"github.com/lintnet/lintnet/pkg/module"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

type Controller struct {
	fs              afero.Fs
	stdout          io.Writer
	moduleInstaller *module.Installer
	importer        *jsonnet.Importer
	param           *ParamController
	dataFileParser  *DataFileParser
	linter          Linter
	fileFinder      FileFinder
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
	return &Controller{
		param:           param,
		fs:              fs,
		stdout:          stdout,
		moduleInstaller: moduleInstaller,
		importer:        importer,
		linter: lint.NewLinter(
			&DataFileParser{
				fs: fs,
			},
			&LintFileParser{
				fs: fs,
			},
			&LintFileEvaluator{
				importer: importer,
			},
		),
		dataFileParser: &DataFileParser{
			fs: fs,
		},
		fileFinder: filefind.NewFileFinder(fs),
	}
}
