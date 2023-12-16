package cli

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/lintnet/pkg/controller/lint"
	"github.com/suzuki-shunsuke/lintnet/pkg/log"
	"github.com/urfave/cli/v2"
)

type lintCommand struct {
	logE *logrus.Entry
}

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
			&cli.StringFlag{
				Name:    "error-level",
				Aliases: []string{"e"},
				EnvVars: []string{"LINTNET_ERROR_LEVEL"},
			},
		},
	}
}

func (lc *lintCommand) action(c *cli.Context) error {
	ctrl := lint.NewController(afero.NewOsFs(), os.Stdout)
	logE := lc.logE
	log.SetLevel(c.String("log-level"), logE)
	log.SetColor(c.String("log-color"), logE)
	return ctrl.Lint(c.Context, logE, &lint.ParamLint{ //nolint:wrapcheck
		FilePaths:   c.Args().Slice(),
		RuleBaseDir: c.String("rule-base-dir"),
		ErrorLevel:  c.String("error-level"),
	})
}
