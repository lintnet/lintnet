package testcmd_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/google/go-jsonnet"
	"github.com/lintnet/lintnet/pkg/controller/testcmd"
	"github.com/lintnet/lintnet/pkg/testutil"
	"github.com/sirupsen/logrus"
)

func TestController_Test(t *testing.T) { //nolint:funlen
	t.Parallel()
	data := []struct {
		name     string
		paramC   *testcmd.ParamController
		param    *testcmd.ParamTest
		files    map[string]string
		contents map[string]string
		dirs     []string
		isErr    bool
	}{
		{
			name: "normal",
			paramC: &testcmd.ParamController{
				Version: "0.3.0",
			},
			param: &testcmd.ParamTest{
				RootDir:        "/home/foo/.local/share/lintnet",
				DataRootDir:    "/home/foo/workspace",
				ConfigFilePath: "",
				PWD:            "/home/foo/workspace",
			},
			files: map[string]string{
				"lintnet.jsonnet":                        "testdata/lintnet.jsonnet",
				"hello.jsonnet":                          "testdata/hello.jsonnet",
				"hello_test.jsonnet":                     "testdata/hello_test.jsonnet",
				"testdata/pass.json":                     "testdata/pass.json",
				"testdata/fail.json":                     "testdata/fail.json",
				"/home/foo/workspace/hello.jsonnet":      "testdata/hello.jsonnet",
				"/home/foo/workspace/hello_test.jsonnet": "testdata/hello_test.jsonnet",
				"/home/foo/workspace/testdata/pass.json": "testdata/pass.json",
				"/home/foo/workspace/testdata/fail.json": "testdata/fail.json",
			},
			contents: map[string]string{},
		},
	}
	for _, d := range data {
		d := d
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
			ctrl := testcmd.NewController(d.paramC, fs, stdout, importer)
			logger := logrus.New()
			logger.SetLevel(logrus.DebugLevel)
			if err := ctrl.Test(context.Background(), logrus.NewEntry(logger), d.param); err != nil {
				if d.isErr {
					return
				}
				t.Fatal(err)
			}
			if d.isErr {
				t.Fatal("error must be returned")
			}
		})
	}
}
