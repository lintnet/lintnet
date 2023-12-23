package lint

import (
	"fmt"

	"github.com/lintnet/lintnet/pkg/errlevel"
)

func isFailed(results []*FlatError, errLevel errlevel.Level) (bool, error) {
	for _, result := range results {
		e := result.Level
		if e == "" {
			e = "error"
		}
		feErrLevel, err := errlevel.New(e)
		if err != nil {
			return false, fmt.Errorf("verify the error level of a result: %w", err)
		}
		if feErrLevel >= errLevel {
			return true, nil
		}
	}
	return false, nil
}
