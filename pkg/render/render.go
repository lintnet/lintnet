package render

import (
	"io"
)

type TemplateRenderer interface {
	Render(out io.Writer, s string, param any) error
}
