package lint

import (
	"errors"

	"github.com/lintnet/lintnet/pkg/domain"
	"github.com/lintnet/lintnet/pkg/errlevel"
	"github.com/lintnet/lintnet/pkg/output"
	"github.com/sirupsen/logrus"
)

type Outputter interface {
	Output(result *output.Output) error
}

func (c *Controller) Output(logE *logrus.Entry, errLevel, shownErrLevel errlevel.Level, results []*domain.Result, outputters []Outputter, outputSuccess bool) error {
	formatter := &output.Formatter{}
	fes := formatter.Format(results, shownErrLevel)
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
