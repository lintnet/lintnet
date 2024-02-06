package config_test

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/lintnet/lintnet/pkg/config"
	"github.com/lintnet/lintnet/pkg/errlevel"
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

func TestRawConfig_Parse(t *testing.T) { //nolint:funlen
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
								SlashPath: "yoo/*.jsonnet",
							},
							{
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
			cfg, err := d.rawCfg.Parse()
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
