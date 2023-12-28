package render

import (
	"fmt"
	"io"
	"text/template"
)

type TextTemplateRenderer struct{}

func (t *TextTemplateRenderer) Compile(s string) (Template, error) {
	tpl, err := template.New("_").Parse(s)
	if err != nil {
		return nil, fmt.Errorf("parse a template: %w", err)
	}
	return tpl, nil
}

func (t *TextTemplateRenderer) Render(out io.Writer, s string, param any) error {
	tpl, err := template.New("_").Parse(s)
	if err != nil {
		return fmt.Errorf("parse a template: %w", err)
	}
	if err := tpl.Execute(out, param); err != nil {
		return fmt.Errorf("render a template: %w", err)
	}
	return nil
}
