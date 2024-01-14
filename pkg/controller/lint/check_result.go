package lint

import (
	"github.com/lintnet/lintnet/pkg/domain"
	"github.com/lintnet/lintnet/pkg/errlevel"
)

func isFailed(results []*domain.FlatError, errLevel errlevel.Level) (bool, error) {
	for _, result := range results {
		f, err := result.Failed(errLevel)
		if err != nil {
			return false, err
		}
		if f {
			return true, nil
		}
	}
	return false, nil
}
