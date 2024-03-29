package lint

import (
	"testing"

	"github.com/lintnet/lintnet/pkg/domain"
	"github.com/lintnet/lintnet/pkg/errlevel"
)

func Test_isFailed(t *testing.T) {
	t.Parallel()
	data := []struct {
		name     string
		results  []*domain.Error
		errLevel errlevel.Level
		exp      bool
		isErr    bool
	}{
		{
			name: "false",
			results: []*domain.Error{
				{
					Level: "warn",
				},
			},
			errLevel: errlevel.Error,
			exp:      false,
		},
		{
			name: "true",
			results: []*domain.Error{
				{
					Level: "warn",
				},
				{
					Level: "error",
				},
			},
			errLevel: errlevel.Error,
			exp:      true,
		},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			f, err := isFailed(d.results, d.errLevel)
			if err != nil {
				if d.isErr {
					return
				}
				t.Fatal(err)
			}
			if d.isErr {
				t.Fatal("error must be returned")
			}
			if f != d.exp {
				t.Fatalf("got %v, wanted %v", f, d.exp)
			}
		})
	}
}
