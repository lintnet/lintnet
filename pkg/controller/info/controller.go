package info

import (
	_ "embed"
	"io"

	"github.com/lintnet/lintnet/pkg/config"
	"github.com/lintnet/lintnet/pkg/domain"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

type Controller struct {
	fs     afero.Fs
	stdout io.Writer
	param  *ParamController
}

type FileFinder interface {
	Find(logE *logrus.Entry, cfg *config.Config, rootDir, cfgDir string) ([]*domain.Target, error)
}

type ParamController struct {
	Version string
}

func NewController(param *ParamController, fs afero.Fs, stdout io.Writer) *Controller {
	return &Controller{
		param:  param,
		fs:     fs,
		stdout: stdout,
	}
}
