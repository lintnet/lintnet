package module_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/lintnet/lintnet/pkg/config"
	"github.com/lintnet/lintnet/pkg/module"
)

func TestListModules(t *testing.T) { //nolint:funlen
	t.Parallel()
	data := []struct {
		name    string
		cfg     *config.Config
		modList [][]*module.Module
		modMap  map[string]*module.Module
		isErr   bool
	}{
		{
			name: "normal",
			cfg: &config.Config{
				Targets: []*config.Target{
					{
						Modules: []*config.Module{
							{
								Path: "github.com/suzuki-shunsuke/example-lintnet-modules/ghalint/job_secrets/main.jsonnet@696511bac987973002692e733735650f86b9c59e:v0.1.3",
							},
						},
					},
				},
			},
			modList: [][]*module.Module{
				{
					{
						Type:      "github",
						Host:      "github.com",
						RepoOwner: "suzuki-shunsuke",
						RepoName:  "example-lintnet-modules",
						Path:      "ghalint/job_secrets/main.jsonnet",
						Ref:       "696511bac987973002692e733735650f86b9c59e",
						Tag:       "v0.1.3",
					},
				},
			},
			modMap: map[string]*module.Module{
				"github.com/suzuki-shunsuke/example-lintnet-modules/696511bac987973002692e733735650f86b9c59e": {
					Type:      "github",
					Host:      "github.com",
					RepoOwner: "suzuki-shunsuke",
					RepoName:  "example-lintnet-modules",
					Path:      "ghalint/job_secrets/main.jsonnet",
					Ref:       "696511bac987973002692e733735650f86b9c59e",
					Tag:       "v0.1.3",
				},
			},
		},
	}
	for _, d := range data {
		d := d
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			modList, modMap, err := module.ListModules(d.cfg)
			if err != nil {
				if d.isErr {
					return
				}
				t.Fatal(err)
			}
			if d.isErr {
				t.Fatal("error must be returned")
			}
			if diff := cmp.Diff(d.modList, modList); diff != "" {
				t.Fatal(diff)
			}
			if diff := cmp.Diff(d.modMap, modMap); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
