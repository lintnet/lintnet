package cli

import (
	"github.com/lintnet/lintnet/pkg/controller/initcmd"
	"github.com/lintnet/lintnet/pkg/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"
)

type initCommand struct {
	logE *logrus.Entry
}

func (lc *initCommand) command() *cli.Command {
	return &cli.Command{
		Name:      "init",
		Usage:     "Scaffold configuration file",
		UsageText: "lintnet init",
		Description: `Scaffold configuration file.

$ lintnet init

This command generates lintnet.jsonnet.
If the file already exists, this command does nothing.
`,
		Action: lc.action,
	}
}

func (lc *initCommand) action(c *cli.Context) error {
	ctrl := initcmd.NewController(afero.NewOsFs())
	logE := lc.logE
	log.SetLevel(c.String("log-level"), logE)
	log.SetColor(c.String("log-color"), logE)
	return ctrl.Init(c.Context, logE) //nolint:wrapcheck
}
