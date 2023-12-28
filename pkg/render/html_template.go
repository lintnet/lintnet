package render

import (
	"fmt"
	hTemplate "html/template"
)

type HTMLTemplateRenderer struct{}

func (t *HTMLTemplateRenderer) Compile(s string) (Template, error) {
	tpl, err := hTemplate.New("_").Parse(s)
	if err != nil {
		return nil, fmt.Errorf("parse a template: %w", err)
	}
	return tpl, nil
}
