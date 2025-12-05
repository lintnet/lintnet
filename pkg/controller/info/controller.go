package info

import (
	"io"
	"log/slog"

	"github.com/lintnet/lintnet/pkg/config"
	"github.com/lintnet/lintnet/pkg/filefind"
	"github.com/spf13/afero"
)

type Controller struct {
	fs     afero.Fs
	stdout io.Writer
	param  *ParamController
}

type FileFinder interface {
	Find(logger *slog.Logger, cfg *config.Config, rootDir, cfgDir string) ([]*filefind.Target, error)
}

type ParamController struct {
	Version string
	Commit  string
	Env     string
}

func NewController(param *ParamController, fs afero.Fs, stdout io.Writer) *Controller {
	return &Controller{
		param:  param,
		fs:     fs,
		stdout: stdout,
	}
}
