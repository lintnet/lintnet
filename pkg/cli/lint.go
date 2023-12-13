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
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "rule-base-dir",
				Aliases: []string{"d"},
				Value:   "lintnet",
			},
		},
	}
}

func (lc *lintCommand) action(c *cli.Context) error {
	ctrl := lint.NewController(afero.NewOsFs(), os.Stdout)
	return ctrl.Lint(c.Context, &lint.ParamLint{ //nolint:wrapcheck
		FilePaths:   c.Args().Slice(),
		RuleBaseDir: c.String("rule-base-dir"),
	})
}
