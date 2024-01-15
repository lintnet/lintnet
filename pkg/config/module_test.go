package config_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
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
					Host:      "github.com",
					RepoOwner: "suzuki-shunsuke",
					RepoName:  "example-lintnet-modules",
					Ref:       "0ed62adf055a4fbd7ef7ebe304f01794508ed325",
				},
			},
			exp: "github.com/suzuki-shunsuke/example-lintnet-modules/0ed62adf055a4fbd7ef7ebe304f01794508ed325/foo/bar.jsonnet",
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

func TestParseModuleLine(t *testing.T) {
	t.Parallel()
	data := []struct {
		name  string
		line  string
		mod   *config.ModuleGlob
		isErr bool
	}{
		{
			name: "module",
			line: "github.com/suzuki-shunsuke/example-lintnet-modules/foo/bar.jsonnet@0ed62adf055a4fbd7ef7ebe304f01794508ed325:v0.1.3",
			mod: &config.ModuleGlob{
				ID:        "github.com/suzuki-shunsuke/example-lintnet-modules/foo/bar.jsonnet@0ed62adf055a4fbd7ef7ebe304f01794508ed325:v0.1.3",
				SlashPath: "github.com/suzuki-shunsuke/example-lintnet-modules/0ed62adf055a4fbd7ef7ebe304f01794508ed325/foo/bar.jsonnet",
				Archive: &config.ModuleArchive{
					Type:      "github",
					Host:      "github.com",
					RepoOwner: "suzuki-shunsuke",
					RepoName:  "example-lintnet-modules",
					Ref:       "0ed62adf055a4fbd7ef7ebe304f01794508ed325",
					Tag:       "v0.1.3",
				},
			},
		},
	}
	for _, d := range data {
		d := d
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			mod, err := config.ParseModuleLine(d.line)
			if err != nil {
				if d.isErr {
					return
				}
				t.Fatal(err)
			}
			if d.isErr {
				t.Fatal("error must be returned")
			}
			if diff := cmp.Diff(d.mod, mod); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
