package cli

import (
	"context"

	"github.com/lintnet/lintnet/pkg/controller/newcmd"
	"github.com/lintnet/lintnet/pkg/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/urfave/cli/v3"
)

type newCommand struct {
	logE *logrus.Entry
}

func (lc *newCommand) command() *cli.Command {
	return &cli.Command{
		Name:      "new",
		Usage:     "Create a lint file and a test file",
		UsageText: "lintnet new [<lint file|main.jsonnet>]",
		Description: `Create a lint file and a test file.

$ lintnet new [<lint file|main.jsonnet>]

This command creates a lint file and a test file.
If the argument is not given, the lint file is created as "main.jsonnet".
`,
		Action: lc.action,
	}
}

func (lc *newCommand) action(ctx context.Context, cmd *cli.Command) error {
	ctrl := newcmd.NewController(afero.NewOsFs())
	logE := lc.logE
	log.SetLevel(cmd.String("log-level"), logE)
	log.SetColor(cmd.String("log-color"), logE)
	fileName := "main.jsonnet"
	if cmd.Args().Len() > 0 {
		fileName = cmd.Args().First()
	}
	return ctrl.New(ctx, logE, fileName) //nolint:wrapcheck
}
