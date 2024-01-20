package info_test

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/lintnet/lintnet/pkg/controller/info"
	"github.com/lintnet/lintnet/pkg/testutil"
)

func TestController_Info(t *testing.T) {
	t.Parallel()
	data := []struct {
		name   string
		isErr  bool
		paramC *info.ParamController
		param  *info.ParamInfo
		exp    string
		files  map[string]string
		dirs   []string
	}{
		{
			name: "normal",
			paramC: &info.ParamController{
				Version: "0.3.0-4",
				Commit:  "a17e309c9d93daa83b546df47ed49c5a56b5250b",
				Env:     "darwin/arm64",
			},
			param: &info.ParamInfo{
				RootDir:        "/home/foo/.local/share/lintnet",
				DataRootDir:    "/home/foo/workspace",
				ConfigFilePath: "",
				PWD:            "/home/foo/workspace",
			},
			exp: "testdata/info.json",
			files: map[string]string{
				"lintnet.jsonnet": "testdata/lintnet.jsonnet",
			},
		},
	}
	for _, d := range data {
		d := d
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			stdout := &bytes.Buffer{}
			fs, err := testutil.ReadFs(d.files, d.dirs...)
			if err != nil {
				t.Fatal(err)
			}
			ctrl := info.NewController(d.paramC, fs, stdout)
			if err := ctrl.Info(context.Background(), d.param); err != nil {
				if d.isErr {
					return
				}
				t.Fatal(err)
			}
			if d.isErr {
				t.Fatal("error must be returned")
			}
			if d.exp != "" {
				b, err := os.ReadFile(d.exp)
				if err != nil {
					t.Fatal(err)
				}
				var exp any
				if err := json.Unmarshal(b, &exp); err != nil {
					t.Fatal(err)
				}
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
