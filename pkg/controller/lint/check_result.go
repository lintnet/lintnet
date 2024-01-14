package lint

import (
	"fmt"

	"github.com/lintnet/lintnet/pkg/domain"
	"github.com/lintnet/lintnet/pkg/errlevel"
)

func isFailed(results []*domain.FlatError, errLevel errlevel.Level) (bool, error) {
	for _, result := range results {
		f, err := result.Failed(errLevel)
		if err != nil {
			return false, fmt.Errorf("check if the command failed: %w", err)
		}
		if f {
			return true, nil
		}
	}
	return false, nil
}
