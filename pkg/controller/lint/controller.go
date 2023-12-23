package lint

import (
	"io"

	"github.com/lintnet/lintnet/pkg/module"
	"github.com/spf13/afero"
)

type Controller struct {
	fs              afero.Fs
	stdout          io.Writer
	moduleInstaller *module.Installer
	importer        *Importer
}

func NewController(fs afero.Fs, stdout io.Writer, moduleInstaller *module.Installer, importer *Importer) *Controller {
	return &Controller{
		fs:              fs,
		stdout:          stdout,
		moduleInstaller: moduleInstaller,
		importer:        importer,
	}
}
