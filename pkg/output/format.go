package output

import (
	"log/slog"

	"github.com/lintnet/lintnet/pkg/domain"
	"github.com/lintnet/lintnet/pkg/errlevel"
	"github.com/suzuki-shunsuke/slog-error/slogerr"
)

type Output struct {
	LintnetVersion string          `json:"lintnet_version"`
	Env            string          `json:"env"`
	Errors         []*domain.Error `json:"errors,omitempty"`
	Config         map[string]any  `json:"config,omitempty"`
}

func FormatResults(logger *slog.Logger, results []*domain.Result, errLevel errlevel.Level) []*domain.Error {
	list := make([]*domain.Error, 0, len(results))
	for _, result := range results {
		for _, fe := range result.FlatErrors() {
			el := errlevel.Error
			invalid := false
			if fe.Level != "" {
				e, err := errlevel.New(fe.Level)
				if err != nil {
					slogerr.WithError(logger, err).Warn("error level is invalid", "lint_file", fe.LintFile, "error_level", fe.Level)
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
