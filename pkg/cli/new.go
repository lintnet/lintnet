package cli

import (
	"context"
	"fmt"

	"github.com/lintnet/lintnet/pkg/controller/newcmd"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/slog-util/slogutil"
	"github.com/suzuki-shunsuke/urfave-cli-v3-util/urfave"
	"github.com/urfave/cli/v3"
)

type newCommand struct{}

func (lc *newCommand) command(logger *slogutil.Logger) *cli.Command {
	return &cli.Command{
		Name:      "new",
		Usage:     "Create a lint file and a test file",
		UsageText: "lintnet new [<lint file|main.jsonnet>]",
		Description: `Create a lint file and a test file.

$ lintnet new [<lint file|main.jsonnet>]

This command creates a lint file and a test file.
If the argument is not given, the lint file is created as "main.jsonnet".
`,
		Action: urfave.Action(lc.action, logger),
	}
}

func (lc *newCommand) action(ctx context.Context, cmd *cli.Command, logger *slogutil.Logger) error {
	ctrl := newcmd.NewController(afero.NewOsFs())
	if err := logger.SetLevel(cmd.String("log-level")); err != nil {
		return fmt.Errorf("set log level: %w", err)
	}
	fileName := "main.jsonnet"
	if cmd.Args().Len() > 0 {
		fileName = cmd.Args().First()
	}
	return ctrl.New(ctx, logger.Logger, fileName) //nolint:wrapcheck
}
