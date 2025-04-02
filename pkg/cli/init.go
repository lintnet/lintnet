package cli

import (
	"context"

	"github.com/lintnet/lintnet/pkg/controller/initcmd"
	"github.com/lintnet/lintnet/pkg/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/urfave/cli/v3"
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

func (lc *initCommand) action(ctx context.Context, cmd *cli.Command) error {
	ctrl := initcmd.NewController(afero.NewOsFs())
	logE := lc.logE
	log.SetLevel(cmd.String("log-level"), logE)
	log.SetColor(cmd.String("log-color"), logE)
	return ctrl.Init(ctx, logE) //nolint:wrapcheck
}
