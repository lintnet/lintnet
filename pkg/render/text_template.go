package render

import (
	"fmt"
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
