package lint

import (
	"context"
	"io"

	gojsonnet "github.com/google/go-jsonnet"
	"github.com/lintnet/lintnet/pkg/config"
	"github.com/lintnet/lintnet/pkg/config/reader"
	"github.com/lintnet/lintnet/pkg/domain"
	"github.com/lintnet/lintnet/pkg/encoding"
	"github.com/lintnet/lintnet/pkg/filefind"
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
	moduleInstaller ModuleInstaller
	importer        gojsonnet.Importer
	param           *ParamController
	dataFileParser  lint.DataFileParser
	linter          Linter
	fileFinder      FileFinder
	configReader    ConfigReader
	outputGetter    OutputGetter
}

type OutputGetter interface {
	Get(outputs []*config.Output, param *output.ParamGet, cfgDir string) (output.Outputter, error)
}

type ConfigReader interface {
	Read(p string, cfg *config.RawConfig) error
}

type ModuleInstaller interface {
	Installs(ctx context.Context, logE *logrus.Entry, param *module.ParamInstall, modules map[string]*config.ModuleArchive) error
}

type MockModuleInstaller struct{}

func (m *MockModuleInstaller) Installs(ctx context.Context, logE *logrus.Entry, param *module.ParamInstall, modules map[string]*config.ModuleArchive) error { //nolint:revive
	return nil
}

type Linter interface {
	Lint(targets []*filefind.Target) ([]*domain.Result, error)
}

type FileFinder interface {
	Find(logE *logrus.Entry, cfg *config.Config, rootDir, cfgDir string) ([]*filefind.Target, error)
}

type ParamController struct {
	Version string
	Env     string
}

func NewController(param *ParamController, fs afero.Fs, stdout io.Writer, moduleInstaller ModuleInstaller, importer gojsonnet.Importer) *Controller {
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
		configReader:   reader.New(fs, importer),
		outputGetter:   output.NewGetter(stdout, fs, importer),
	}
}
