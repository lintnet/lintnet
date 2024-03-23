package filefind_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/lintnet/lintnet/pkg/config"
	"github.com/lintnet/lintnet/pkg/domain"
	"github.com/lintnet/lintnet/pkg/filefind"
	"github.com/lintnet/lintnet/pkg/testutil"
	"github.com/sirupsen/logrus"
)

func TestFinder_Find(t *testing.T) { //nolint:funlen
	t.Parallel()
	data := []struct {
		name    string
		files   map[string]string
		dirs    []string
		cfg     *config.Config
		rootDir string
		cfgDir  string
		targets []*filefind.Target
		isErr   bool
	}{
		{
			name: "normal",
			files: map[string]string{
				"foo.json":           `{}`,
				"hello.jsonnet":      `{}`,
				"hello_test.jsonnet": `{}`,
			},
			cfg: &config.Config{
				Targets: []*config.Target{
					{
						LintFiles: []*config.LintGlob{
							{
								Glob: "*.jsonnet",
							},
						},
						DataFiles: []*config.DataFile{
							{
								Path: "*.json",
							},
						},
					},
				},
			},
			rootDir: "/home/foo/.local/share/lintnet",
			cfgDir:  "",
			targets: []*filefind.Target{
				{
					DataFiles: domain.Paths{
						{
							Raw: "foo.json",
							Abs: "foo.json",
						},
					},
					LintFiles: []*config.LintFile{
						{
							ID:   "hello.jsonnet",
							Path: "hello.jsonnet",
						},
					},
				},
			},
		},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			fs, err := testutil.NewFs(d.files, d.dirs...)
			if err != nil {
				t.Fatal(err)
			}
			finder := filefind.NewFileFinder(fs)
			targets, err := finder.Find(logrus.NewEntry(logrus.New()), d.cfg, d.rootDir, d.cfgDir)
			if err != nil {
				if d.isErr {
					t.Fatal(err)
				}
			}
			if d.isErr {
				t.Fatal("error must be returned")
			}
			if diff := cmp.Diff(d.targets, targets); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
