package render_test

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/lintnet/lintnet/pkg/render"
)

func TestRenderer_Render(t *testing.T) {
	t.Parallel()
	data := []struct {
		name     string
		s        string
		param    interface{}
		renderer render.TemplateRenderer
		exp      string
		isErr    bool
	}{
		{
			name:     "text/template",
			s:        "Name: {{.name}}",
			renderer: &render.TextTemplateRenderer{},
			param: map[string]any{
				"name": "mike",
			},
			exp: "Name: mike",
		},
		{
			name:     "html/template",
			s:        "Name: {{.name}}",
			renderer: &render.HTMLTemplateRenderer{},
			param: map[string]any{
				"name": "mike",
			},
			exp: "Name: mike",
		},
	}
	for _, d := range data {
		d := d
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			w := &bytes.Buffer{}
			if err := d.renderer.Render(w, d.s, d.param); err != nil {
				if d.isErr {
					return
				}
				t.Fatal(err)
			}
			if d.isErr {
				t.Fatal("error must be returned")
			}
			s := w.String()
			if diff := cmp.Diff(d.exp, s); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
