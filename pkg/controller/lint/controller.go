package lint

import (
	"io"

	"github.com/spf13/afero"
)

type Controller struct {
	fs     afero.Fs
	stdout io.Writer
	gh     GitHub
	http   HTTPClient
}

func NewController(fs afero.Fs, stdout io.Writer, gh GitHub, httpClient HTTPClient) *Controller {
	return &Controller{
		fs:     fs,
		stdout: stdout,
		gh:     gh,
		http:   httpClient,
	}
}
