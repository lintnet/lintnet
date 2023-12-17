package lint

import (
	"fmt"
	hTemplate "html/template"
	"io"
	"text/template"
)

type TemplateRenderer interface {
	Render(out io.Writer, s string, param interface{}) error
}

type TextTemplateRenderer struct{}

func (t *TextTemplateRenderer) Render(out io.Writer, s string, param interface{}) error {
	tpl, err := template.New("_").Parse(s)
	if err != nil {
		return fmt.Errorf("parse a template: %w", err)
	}
	if err := tpl.Execute(out, param); err != nil {
		return fmt.Errorf("render a template: %w", err)
	}
	return nil
}

type HTMLTemplateRenderer struct{}

func (t *HTMLTemplateRenderer) Render(out io.Writer, s string, param interface{}) error {
	tpl, err := hTemplate.New("_").Parse(s)
	if err != nil {
		return fmt.Errorf("parse a template: %w", err)
	}
	if err := tpl.Execute(out, param); err != nil {
		return fmt.Errorf("render a template: %w", err)
	}
	return nil
}
