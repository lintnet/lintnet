package output

import (
	"github.com/lintnet/lintnet/pkg/domain"
	"github.com/lintnet/lintnet/pkg/errlevel"
	"github.com/sirupsen/logrus"
)

type Output struct {
	LintnetVersion string          `json:"lintnet_version"`
	Env            string          `json:"env"`
	Errors         []*domain.Error `json:"errors,omitempty"`
	Config         map[string]any  `json:"config,omitempty"`
}

func FormatResults(logE *logrus.Entry, results []*domain.Result, errLevel errlevel.Level) []*domain.Error {
	list := make([]*domain.Error, 0, len(results))
	for _, result := range results {
		for _, fe := range result.FlatErrors() {
			el := errlevel.Error
			invalid := false
			if fe.Level != "" {
				e, err := errlevel.New(fe.Level)
				if err != nil {
					logE.WithError(err).WithFields(logrus.Fields{
						"lint_file":   fe.LintFile,
						"error_level": fe.Level,
					}).Warn("error level is invalid")
					invalid = true
				} else {
					el = e
				}
			}
			if invalid || el >= errLevel {
				list = append(list, fe)
			}
		}
	}
	return list
}
