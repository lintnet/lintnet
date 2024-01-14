package lint

import (
	"io"

	"github.com/lintnet/lintnet/pkg/jsonnet"
	"github.com/lintnet/lintnet/pkg/module"
	"github.com/spf13/afero"
)

type Controller struct {
	fs                afero.Fs
	stdout            io.Writer
	moduleInstaller   *module.Installer
	importer          *jsonnet.Importer
	param             *ParamController
	lintFileParser    *LintFileParser
	lintFileEvaluator *LintFileEvaluator
	dataFileParser    *DataFileParser
	linter            *Linter
	fileFinder        *FileFinder
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
		linter: &Linter{
			lintFileParser: &LintFileParser{
				fs: fs,
			},
			lintFileEvaluator: &LintFileEvaluator{
				importer: importer,
			},
			dataFileParser: &DataFileParser{
				fs: fs,
			},
		},
		lintFileParser: &LintFileParser{
			fs: fs,
		},
		lintFileEvaluator: &LintFileEvaluator{
			importer: importer,
		},
		dataFileParser: &DataFileParser{
			fs: fs,
		},
		fileFinder: &FileFinder{
			fs: fs,
		},
	}
}
