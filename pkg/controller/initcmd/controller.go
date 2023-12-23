package initcmd

import (
	"github.com/spf13/afero"
)

type Controller struct {
	fs afero.Fs
}

func NewController(fs afero.Fs) *Controller {
	return &Controller{
		fs: fs,
	}
}
