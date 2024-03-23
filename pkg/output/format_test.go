package output_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/lintnet/lintnet/pkg/domain"
	"github.com/lintnet/lintnet/pkg/errlevel"
	"github.com/lintnet/lintnet/pkg/output"
	"github.com/sirupsen/logrus"
)

func TestFormatResults(t *testing.T) {
	t.Parallel()
	data := []struct {
		name     string
		results  []*domain.Result
		errLevel errlevel.Level
		exp      []*domain.Error
	}{
		{
			name: "normal",
			results: []*domain.Result{
				{
					LintFile: "hello.jsonnet",
					DataFile: "hello.json",
					RawResult: []*domain.JsonnetResult{
						{
							Name: "description is required",
						},
						{
							Name:  "description is required",
							Level: "info",
						},
					},
				},
			},
			errLevel: errlevel.Error,
			exp: []*domain.Error{
				{
					Name:     "description is required",
					LintFile: "hello.jsonnet",
					DataFile: "hello.json",
				},
			},
		},
	}
	logE := logrus.NewEntry(logrus.New())
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			errs := output.FormatResults(logE, d.results, d.errLevel)
			if diff := cmp.Diff(d.exp, errs); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
