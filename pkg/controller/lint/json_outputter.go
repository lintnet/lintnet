package lint

import (
	"io"
)

type jsonOutputter struct {
	stdout io.Writer
}

func (o *jsonOutputter) Output(result *Output) error {
	return outputJSON(o.stdout, result)
}
