package cli

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/lintnet/lintnet/pkg/controller/newcmd"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/slog-util/slogutil"
	"github.com/urfave/cli/v3"
)

type newCommand struct {
	logger      *slog.Logger
	logLevelVar *slog.LevelVar
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
	if err := slogutil.SetLevel(lc.logLevelVar, cmd.String("log-level")); err != nil {
		return fmt.Errorf("set log level: %w", err)
	}
	fileName := "main.jsonnet"
	if cmd.Args().Len() > 0 {
		fileName = cmd.Args().First()
	}
	return ctrl.New(ctx, lc.logger, fileName) //nolint:wrapcheck
}
