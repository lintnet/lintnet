package filefilter_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/lintnet/lintnet/pkg/config"
	"github.com/lintnet/lintnet/pkg/domain"
	"github.com/lintnet/lintnet/pkg/filefilter"
	"github.com/lintnet/lintnet/pkg/filefind"
	"github.com/sirupsen/logrus"
)

func TestFilterTargetsByFilePaths(t *testing.T) { //nolint:funlen
	t.Parallel()
	data := []struct {
		name    string
		param   *filefilter.Param
		targets []*filefind.Target
		exp     []*filefind.Target
	}{
		{
			name: "normal",
			param: &filefilter.Param{
				DataRootDir: "/home/foo/workspace",
				PWD:         "/home/foo/workspace",
				FilePaths: []string{
					"foo.json",
				},
			},
			targets: []*filefind.Target{
				{
					LintFiles: []*config.LintFile{
						{
							ID:   "hello.jsonnet",
							Path: "hello.jsonnet",
						},
					},
					DataFiles: domain.Paths{
						{
							Raw: "foo.json",
							Abs: "/home/foo/workspace/foo.json",
						},
						{
							Raw: "bar.json",
							Abs: "/home/foo/workspace/bar.json",
						},
					},
				},
			},
			exp: []*filefind.Target{
				{
					LintFiles: []*config.LintFile{
						{
							ID:   "hello.jsonnet",
							Path: "hello.jsonnet",
						},
					},
					DataFiles: domain.Paths{
						{
							Raw: "foo.json",
							Abs: "/home/foo/workspace/foo.json",
						},
					},
				},
			},
		},
	}
	for _, d := range data {
		d := d
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			targets := filefilter.FilterTargetsByFilePaths(d.param, d.targets)
			if diff := cmp.Diff(d.exp, targets); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestFilterTargetsByDataRootDir(t *testing.T) { //nolint:funlen
	t.Parallel()
	data := []struct {
		name    string
		param   *filefilter.Param
		targets []*filefind.Target
		exp     []*filefind.Target
		isErr   bool
	}{
		{
			name: "normal",
			param: &filefilter.Param{
				DataRootDir: "/home/foo/workspace",
				PWD:         "/home/foo/workspace",
			},
			targets: []*filefind.Target{
				{
					LintFiles: []*config.LintFile{
						{
							ID:   "hello.jsonnet",
							Path: "hello.jsonnet",
						},
					},
					DataFiles: domain.Paths{
						{
							Raw: "foo.json",
							Abs: "/home/foo/workspace/foo.json",
						},
						{
							Raw: "bar.json",
							Abs: "/home/bar/workspace/bar.json",
						},
					},
				},
				{
					LintFiles: []*config.LintFile{
						{
							ID:   "foo.jsonnet",
							Path: "foo.jsonnet",
						},
					},
					DataFiles: domain.Paths{
						{
							Raw: "foo.json",
							Abs: "/home/bar/workspace/foo.json",
						},
						{
							Raw: "bar.json",
							Abs: "/home/bar/workspace/bar.json",
						},
					},
				},
			},
			exp: []*filefind.Target{
				{
					LintFiles: []*config.LintFile{
						{
							ID:   "hello.jsonnet",
							Path: "hello.jsonnet",
						},
					},
					DataFiles: domain.Paths{
						{
							Raw: "foo.json",
							Abs: "/home/foo/workspace/foo.json",
						},
					},
				},
			},
		},
	}
	for _, d := range data {
		d := d
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			targets, err := filefilter.FilterTargetsByDataRootDir(logrus.NewEntry(logrus.New()), d.param, d.targets)
			if err != nil {
				if d.isErr {
					return
				}
				t.Fatal(err)
			}
			if d.isErr {
				t.Fatal("error must be returned")
			}
			if diff := cmp.Diff(d.exp, targets); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
