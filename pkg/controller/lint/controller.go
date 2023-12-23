package lint

import (
	"io"

	"github.com/lintnet/lintnet/pkg/jsonnet"
	"github.com/lintnet/lintnet/pkg/module"
	"github.com/spf13/afero"
)

type Controller struct {
	fs              afero.Fs
	stdout          io.Writer
	moduleInstaller *module.Installer
	importer        *jsonnet.Importer
}

func NewController(fs afero.Fs, stdout io.Writer, moduleInstaller *module.Installer, importer *jsonnet.Importer) *Controller {
	return &Controller{
		fs:              fs,
		stdout:          stdout,
		moduleInstaller: moduleInstaller,
		importer:        importer,
	}
}
