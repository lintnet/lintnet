package encoding_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/lintnet/lintnet/pkg/encoding"
)

func TestNewUnmarshaler(t *testing.T) { //nolint:funlen
	t.Parallel()
	data := []struct {
		name     string
		fileName string
		data     string
		exp      any
		fileType string
		isErr    bool
	}{
		{
			name:     "csv",
			fileName: "hello.csv",
			data: `mike,20
`,
			exp:      [][]string{{"mike", "20"}},
			fileType: "csv",
		},
		{
			name:     "hcl2",
			fileName: "hello.hcl",
			data: `resource "null_resource" "hello" {}
`,
			exp: map[string]any{
				"resource": map[string]any{
					"null_resource": map[string]any{
						"hello": []any{
							map[string]any{},
						},
					},
				},
			},
			fileType: "hcl2",
		},
		{
			name:     "json",
			fileName: "hello.json",
			data: `{"name": "hello"}
`,
			exp: map[string]any{
				"name": "hello",
			},
			fileType: "json",
		},
		{
			name:     "plain",
			fileName: "hello.txt",
			data:     `Hello`,
			fileType: "plain_text",
		},
		{
			name:     "toml",
			fileName: "hello.toml",
			data: `name = "hello"
`,
			exp: map[string]any{
				"name": "hello",
			},
			fileType: "toml",
		},
		{
			name:     "yaml",
			fileName: "hello.yaml",
			data: `name: hello
`,
			exp: []any{
				map[string]any{
					"name": "hello",
				},
			},
			fileType: "yaml",
		},
	}
	for _, d := range data {
		d := d
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			unmarshaler, fileType, err := encoding.NewUnmarshaler(d.fileName)
			if err != nil {
				if d.isErr {
					return
				}
				t.Fatal(err)
			}
			if d.isErr {
				t.Fatal("error must be returned")
			}
			if fileType != d.fileType {
				t.Fatalf("got %s, wanted %s", fileType, d.fileType)
			}
			data, err := unmarshaler.Unmarshal([]byte(d.data))
			if err != nil {
				if d.isErr {
					return
				}
				t.Fatal(err)
			}
			if d.isErr {
				t.Fatal("error must be returned")
			}
			if diff := cmp.Diff(d.exp, data); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
