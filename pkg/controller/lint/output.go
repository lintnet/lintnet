package lint

import (
	"errors"
	"fmt"
	"runtime"

	"github.com/lintnet/lintnet/pkg/domain"
	"github.com/lintnet/lintnet/pkg/errlevel"
	"github.com/lintnet/lintnet/pkg/output"
	"github.com/sirupsen/logrus"
)

type Outputter interface {
	Output(result *output.Output) error
}

func (c *Controller) Output(logE *logrus.Entry, errLevel, shownErrLevel errlevel.Level, results []*domain.Result, outputters []Outputter, outputSuccess bool) error {
	fes := &output.Output{
		Errors:         output.FormatResults(logE, results, shownErrLevel),
		LintnetVersion: c.param.Version,
		Env:            fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
	failed, err := isFailed(fes.Errors, errLevel)
	if err != nil {
		return err
	}
	if !outputSuccess && len(fes.Errors) == 0 {
		return nil
	}
	for _, outputter := range outputters {
		if err := outputter.Output(fes); err != nil {
			logE.WithError(err).Error("output errors")
		}
	}
	if failed {
		return errors.New("lint failed")
	}
	return nil
}
