package cli

import (
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

func (lc *newCommand) action(c *cli.Context) error {
	ctrl := newcmd.NewController(afero.NewOsFs())
	logE := lc.logE
	log.SetLevel(c.String("log-level"), logE)
	log.SetColor(c.String("log-color"), logE)
	fileName := "main.jsonnet"
	if c.Args().Len() > 0 {
		fileName = c.Args().First()
	}
	return ctrl.New(c.Context, logE, fileName) //nolint:wrapcheck
}
