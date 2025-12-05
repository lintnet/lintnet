package cli

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/lintnet/lintnet/pkg/controller/initcmd"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/slog-util/slogutil"
	"github.com/urfave/cli/v3"
)

type initCommand struct {
	logger      *slog.Logger
	logLevelVar *slog.LevelVar
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
	if err := slogutil.SetLevel(lc.logLevelVar, cmd.String("log-level")); err != nil {
		return fmt.Errorf("set log level: %w", err)
	}
	return ctrl.Init(ctx, lc.logger) //nolint:wrapcheck
}
