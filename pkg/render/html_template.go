package render

import (
	"fmt"
	hTemplate "html/template"
	"io"
)

type HTMLTemplateRenderer struct{}

func (t *HTMLTemplateRenderer) Compile(s string) (Template, error) {
	tpl, err := hTemplate.New("_").Parse(s)
	if err != nil {
		return nil, fmt.Errorf("parse a template: %w", err)
	}
	return tpl, nil
}

func (t *HTMLTemplateRenderer) Render(out io.Writer, s string, param any) error {
	tpl, err := hTemplate.New("_").Parse(s)
	if err != nil {
		return fmt.Errorf("parse a template: %w", err)
	}
	if err := tpl.Execute(out, param); err != nil {
		return fmt.Errorf("render a template: %w", err)
	}
	return nil
}
