package render

import (
	"io"
)

type TemplateRenderer interface {
	Compile(s string) (Template, error)
}

type Template interface {
	Execute(wr io.Writer, data any) error
}
