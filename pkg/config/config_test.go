package config_test

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/lintnet/lintnet/pkg/config"
)

func TestModule_UnmarshalJSON(t *testing.T) {
	t.Parallel()
	data := []struct {
		name  string
		s     string
		exp   *config.Module
		isErr bool
	}{
		{
			name: "string",
			s:    `"github.com/suzuki-shunsuke/example-lintnet-modules/ghalint/job_secrets/main.jsonnet@696511bac987973002692e733735650f86b9c59e:v0.1.3"`,
			exp: &config.Module{
				Path: "github.com/suzuki-shunsuke/example-lintnet-modules/ghalint/job_secrets/main.jsonnet@696511bac987973002692e733735650f86b9c59e:v0.1.3",
			},
		},
		{
			name: "string",
			s: `{
				"path": "github.com/suzuki-shunsuke/example-lintnet-modules/ghalint/job_secrets/main.jsonnet@696511bac987973002692e733735650f86b9c59e:v0.1.3",
				"param": {
					"excludes": ["foo"]
				}
				}`,
			exp: &config.Module{
				Path: "github.com/suzuki-shunsuke/example-lintnet-modules/ghalint/job_secrets/main.jsonnet@696511bac987973002692e733735650f86b9c59e:v0.1.3",
				Param: map[string]any{
					"excludes": []any{"foo"},
				},
			},
		},
	}
	for _, d := range data {
		d := d
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			m := &config.Module{}
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
