package lint_test

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-jsonnet"
	"github.com/lintnet/lintnet/pkg/controller/lint"
	"github.com/lintnet/lintnet/pkg/testutil"
	"github.com/sirupsen/logrus"
)

func TestController_Lint(t *testing.T) { //nolint:funlen,gocognit,cyclop
	t.Parallel()
	data := []struct {
		name     string
		isErr    bool
		param    *lint.ParamLint
		paramC   *lint.ParamController
		files    map[string]string
		dirs     []string
		contents map[string]string
		exp      string
	}{
		{
			name: "normal",
			param: &lint.ParamLint{
				RootDir:        "/home/foo/.local/share/lintnet",
				DataRootDir:    "/home/foo/workspace",
				ConfigFilePath: "",
				PWD:            "/home/foo/workspace",
			},
			paramC: &lint.ParamController{
				Version: "v0.3.0",
				Env:     "darwin/arm64",
			},
			files: map[string]string{
				"lintnet.jsonnet":                   "testdata/lintnet.jsonnet",
				"/home/foo/workspace/foo.json":      "testdata/foo.json",
				"/home/foo/workspace/hello.jsonnet": "testdata/hello.jsonnet",
			},
			dirs:     []string{},
			contents: map[string]string{},
			exp:      "testdata/result.json",
		},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			fs, err := testutil.ReadFs(d.files, d.dirs...)
			if err != nil {
				t.Fatal(err)
			}
			stdout := &bytes.Buffer{}
			data := make(map[string]jsonnet.Contents, len(d.contents))
			for k, v := range d.contents {
				data[k] = jsonnet.MakeContents(v)
			}
			importer := &jsonnet.MemoryImporter{
				Data: data,
			}
			ctrl := lint.NewController(d.paramC, fs, stdout, &lint.MockModuleInstaller{}, importer)
			ctx := context.Background()
			logE := logrus.NewEntry(logrus.New())
			var exp any
			if d.exp != "" {
				b, err := os.ReadFile(d.exp)
				if err != nil {
					t.Fatal(err)
				}
				if err := json.Unmarshal(b, &exp); err != nil {
					t.Fatal(err)
				}
			}
			if err := ctrl.Lint(ctx, logE, d.param); err != nil {
				if d.isErr {
					return
				}
			}
			if d.isErr {
				t.Fatal("error must be returned")
			}
			if d.exp != "" {
				var result any
				if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
					t.Fatal(err)
				}
				if diff := cmp.Diff(exp, result); diff != "" {
					t.Fatal(diff)
				}
			}
		})
	}
}
