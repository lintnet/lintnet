package config_test

import (
	"testing"

	"github.com/lintnet/lintnet/pkg/config"
)

func TestModule_FilePath(t *testing.T) {
	t.Parallel()
	data := []struct {
		name string
		mod  *config.Module
		exp  string
	}{
		{
			name: "not module",
			mod: &config.Module{
				SlashPath: "foo/bar.jsonnet",
			},
			exp: "foo/bar.jsonnet",
		},
		{
			name: "module",
			mod: &config.Module{
				SlashPath: "foo/bar.jsonnet",
				Archive: &config.ModuleArchive{
					Type:      "github_archive",
					Host:      "github.com",
					RepoOwner: "suzuki-shunsuke",
					RepoName:  "example-lintnet-modules",
					Ref:       "0ed62adf055a4fbd7ef7ebe304f01794508ed325",
				},
			},
			exp: "github_archive/github.com/suzuki-shunsuke/example-lintnet-modules/0ed62adf055a4fbd7ef7ebe304f01794508ed325/foo/bar.jsonnet",
		},
	}
	for _, d := range data {
		d := d
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			s := d.mod.FilePath()
			if s != d.exp {
				t.Fatalf("got %s, wanted %s", s, d.exp)
			}
		})
	}
}
