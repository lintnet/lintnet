package output

import (
	"fmt"
	"runtime"

	"github.com/lintnet/lintnet/pkg/domain"
	"github.com/lintnet/lintnet/pkg/errlevel"
)

type Output struct {
	LintnetVersion string              `json:"lintnet_version"`
	Env            string              `json:"env"`
	Errors         []*domain.FlatError `json:"errors,omitempty"`
	Config         map[string]any      `json:"config,omitempty"`
}

type Formatter struct {
	version string
}

func (f *Formatter) Format(results []*domain.Result, errLevel errlevel.Level) *Output {
	list := make([]*domain.FlatError, 0, len(results))
	for _, result := range results {
		for _, fe := range result.FlatErrors() {
			el, err := errlevel.New(fe.Level) // TODO output warning
			if err != nil || el >= errLevel {
				list = append(list, fe)
			}
		}
	}
	return &Output{
		LintnetVersion: f.version,
		Env:            fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		Errors:         list,
	}
}
