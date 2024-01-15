package config_test

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/lintnet/lintnet/pkg/config"
)

func TestRawModule_UnmarshalJSON(t *testing.T) {
	t.Parallel()
	data := []struct {
		name  string
		s     string
		exp   *config.RawModule
		isErr bool
	}{
		{
			name: "string",
			s:    `"github_archive/github.com/suzuki-shunsuke/example-lintnet-modules/ghalint/job_secrets/main.jsonnet@696511bac987973002692e733735650f86b9c59e:v0.1.3"`,
			exp: &config.RawModule{
				Glob: "github_archive/github.com/suzuki-shunsuke/example-lintnet-modules/ghalint/job_secrets/main.jsonnet@696511bac987973002692e733735650f86b9c59e:v0.1.3",
			},
		},
		{
			name: "string",
			s: `{
				"path": "github_archive/github.com/suzuki-shunsuke/example-lintnet-modules/ghalint/job_secrets/main.jsonnet@696511bac987973002692e733735650f86b9c59e:v0.1.3",
				"param": {
					"excludes": ["foo"]
				}
				}`,
			exp: &config.RawModule{
				Glob: "github_archive/github.com/suzuki-shunsuke/example-lintnet-modules/ghalint/job_secrets/main.jsonnet@696511bac987973002692e733735650f86b9c59e:v0.1.3",
				Config: map[string]any{
					"excludes": []any{"foo"},
				},
			},
		},
	}
	for _, d := range data {
		d := d
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			m := &config.RawModule{}
			if err := json.Unmarshal([]byte(d.s), m); err != nil {
				if d.isErr {
					return
				}
				t.Fatal(err)
			}
			if d.isErr {
				t.Fatal("error must be returned")
			}
			if diff := cmp.Diff(d.exp, m); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
