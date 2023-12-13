package cli

import (
	"os"

	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/lintnet/pkg/controller/lint"
	"github.com/urfave/cli/v2"
)

type lintCommand struct{}

func (lc *lintCommand) command() *cli.Command {
	return &cli.Command{
		Name:   "lint",
		Usage:  "Lint files",
		Action: lc.action,
	}
}

func (lc *lintCommand) action(c *cli.Context) error {
	ctrl := lint.NewController(afero.NewOsFs(), os.Stdout)
	return ctrl.Lint(c.Context, c.Args().Slice()...) //nolint:wrapcheck
}
