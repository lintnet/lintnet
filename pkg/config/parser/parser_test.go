package parser_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/lintnet/lintnet/pkg/config"
	"github.com/lintnet/lintnet/pkg/config/parser"
	"github.com/lintnet/lintnet/pkg/errlevel"
)

func TestParse(t *testing.T) { //nolint:funlen
	t.Parallel()
	data := []struct {
		name   string
		rawCfg *config.RawConfig
		cfg    *config.Config
		isErr  bool
	}{
		{
			name: "normal",
			rawCfg: &config.RawConfig{
				ErrorLevel: "info",
				Targets: []*config.RawTarget{
					{
						LintGlobs: []*config.LintGlob{
							{
								Glob: "yoo/*.jsonnet",
							},
							{
								Glob: "zoo/*.jsonnet",
								Config: map[string]any{
									"limit": 10,
								},
							},
						},
						Modules: []*config.RawModule{
							{
								Glob: "github_archive/github.com/suzuki-shunsuke/example-lintnet-modules/foo/*.jsonnet@0ed62adf055a4fbd7ef7ebe304f01794508ed325:v0.1.3",
								Config: map[string]any{
									"excluded": []string{"test"},
								},
							},
							{
								Glob: "github_archive/github.com/suzuki-shunsuke/example-lintnet-modules/bar/*.jsonnet@0ed62adf055a4fbd7ef7ebe304f01794508ed325:v0.1.3",
							},
						},
						DataFiles: []string{
							".github/workflows/*.yml",
							".github/workflows/*.yaml",
							"!.github/workflows/excluded.yaml",
						},
					},
				},
			},
			cfg: &config.Config{
				ErrorLevel:      errlevel.Info,
				ShownErrorLevel: errlevel.Info,
				IgnoredPatterns: []string{
					"**/.git/**",
					"**/node_modules/**",
				},
				Targets: []*config.Target{
					{
						LintFiles: []*config.ModuleGlob{
							{
								ID:        "yoo/*.jsonnet",
								SlashPath: "yoo/*.jsonnet",
							},
							{
								ID:        "zoo/*.jsonnet",
								SlashPath: "zoo/*.jsonnet",
								Config: map[string]any{
									"limit": 10,
								},
							},
						},
						DataFiles: []string{
							".github/workflows/*.yml",
							".github/workflows/*.yaml",
							"!.github/workflows/excluded.yaml",
						},
						Modules: []*config.ModuleGlob{
							{
								ID:        "github_archive/github.com/suzuki-shunsuke/example-lintnet-modules/foo/*.jsonnet@0ed62adf055a4fbd7ef7ebe304f01794508ed325:v0.1.3",
								SlashPath: "github_archive/github.com/suzuki-shunsuke/example-lintnet-modules/0ed62adf055a4fbd7ef7ebe304f01794508ed325/foo/*.jsonnet",
								Archive: &config.ModuleArchive{
									Type:      "github_archive",
									Host:      "github.com",
									RepoOwner: "suzuki-shunsuke",
									RepoName:  "example-lintnet-modules",
									Ref:       "0ed62adf055a4fbd7ef7ebe304f01794508ed325",
									Tag:       "v0.1.3",
								},
								Config: map[string]any{
									"excluded": []string{"test"},
								},
							},
							{
								ID:        "github_archive/github.com/suzuki-shunsuke/example-lintnet-modules/bar/*.jsonnet@0ed62adf055a4fbd7ef7ebe304f01794508ed325:v0.1.3",
								SlashPath: "github_archive/github.com/suzuki-shunsuke/example-lintnet-modules/0ed62adf055a4fbd7ef7ebe304f01794508ed325/bar/*.jsonnet",
								Archive: &config.ModuleArchive{
									Type:      "github_archive",
									Host:      "github.com",
									RepoOwner: "suzuki-shunsuke",
									RepoName:  "example-lintnet-modules",
									Ref:       "0ed62adf055a4fbd7ef7ebe304f01794508ed325",
									Tag:       "v0.1.3",
								},
							},
						},
						ModuleArchives: map[string]*config.ModuleArchive{
							"github_archive/github.com/suzuki-shunsuke/example-lintnet-modules/0ed62adf055a4fbd7ef7ebe304f01794508ed325:v0.1.3": {
								Type:      "github_archive",
								Host:      "github.com",
								RepoOwner: "suzuki-shunsuke",
								RepoName:  "example-lintnet-modules",
								Ref:       "0ed62adf055a4fbd7ef7ebe304f01794508ed325",
								Tag:       "v0.1.3",
							},
						},
					},
				},
				ModuleArchives: map[string]*config.ModuleArchive{
					"github_archive/github.com/suzuki-shunsuke/example-lintnet-modules/0ed62adf055a4fbd7ef7ebe304f01794508ed325:v0.1.3": {
						Type:      "github_archive",
						Host:      "github.com",
						RepoOwner: "suzuki-shunsuke",
						RepoName:  "example-lintnet-modules",
						Ref:       "0ed62adf055a4fbd7ef7ebe304f01794508ed325",
						Tag:       "v0.1.3",
					},
				},
			},
		},
	}
	for _, d := range data {
		d := d
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			cfg, err := parser.Parse(d.rawCfg)
			if err != nil {
				if d.isErr {
					return
				}
				t.Fatal(err)
			}
			if d.isErr {
				t.Fatal("error must be returned")
			}
			if diff := cmp.Diff(d.cfg, cfg); diff != "" {
				t.Fatal(diff)
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
			line: "github_archive/github.com/suzuki-shunsuke/example-lintnet-modules/foo/bar.jsonnet@0ed62adf055a4fbd7ef7ebe304f01794508ed325:v0.1.3",
			mod: &config.ModuleGlob{
				ID:        "github_archive/github.com/suzuki-shunsuke/example-lintnet-modules/foo/bar.jsonnet@0ed62adf055a4fbd7ef7ebe304f01794508ed325:v0.1.3",
				SlashPath: "github_archive/github.com/suzuki-shunsuke/example-lintnet-modules/0ed62adf055a4fbd7ef7ebe304f01794508ed325/foo/bar.jsonnet",
				Archive: &config.ModuleArchive{
					Type:      "github_archive",
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
			mod, err := parser.ParseModuleLine(d.line)
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
