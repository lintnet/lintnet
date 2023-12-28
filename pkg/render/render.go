package render

import (
	"io"
)

type TemplateRenderer interface {
	Render(out io.Writer, s string, param any) error
	Compile(s string) (Template, error)
}

type Template interface {
	Execute(io.Writer, any) error
}
