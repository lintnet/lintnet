package initcmd_test

import (
	"context"
	"testing"

	"github.com/lintnet/lintnet/pkg/controller/initcmd"
	"github.com/lintnet/lintnet/pkg/testutil"
	"github.com/sirupsen/logrus"
)

func TestController_Init(t *testing.T) {
	t.Parallel()
	data := []struct {
		name  string
		files map[string]string
		dirs  []string
		isErr bool
	}{
		{
			name: "normal",
		},
		{
			name: "already exists",
			files: map[string]string{
				"lintnet.jsonnet": "testdata/lintnet.jsonnet",
			},
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
			ctrl := initcmd.NewController(fs)
			if err := ctrl.Init(context.Background(), logrus.NewEntry(logrus.New())); err != nil {
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
