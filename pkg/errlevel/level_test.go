package errlevel_test

import (
	"testing"

	"github.com/lintnet/lintnet/pkg/errlevel"
)

func TestNew(t *testing.T) {
	t.Parallel()
	data := []struct {
		name  string
		input string
		level errlevel.Level
		isErr bool
	}{
		{
			name:  "info",
			input: "info",
			level: errlevel.Info,
		},
		{
			name:  "invalid",
			input: "invalid",
			isErr: true,
		},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			l, err := errlevel.New(d.input)
			if err != nil {
				if d.isErr {
					return
				}
				t.Fatal(err)
			}
			if d.isErr {
				t.Fatal("error must be returned")
			}
			if d.level != l {
				t.Fatalf("got %v, wanted %v", l, d.level)
			}
		})
	}
}
