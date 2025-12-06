package cli

import (
	"context"
	"fmt"

	"github.com/lintnet/lintnet/pkg/controller/initcmd"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/slog-util/slogutil"
	"github.com/suzuki-shunsuke/urfave-cli-v3-util/urfave"
	"github.com/urfave/cli/v3"
)

type initCommand struct{}

func (lc *initCommand) command(logger *slogutil.Logger) *cli.Command {
	return &cli.Command{
		Name:      "init",
		Usage:     "Scaffold configuration file",
		UsageText: "lintnet init",
		Description: `Scaffold configuration file.

$ lintnet init

This command generates lintnet.jsonnet.
If the file already exists, this command does nothing.
`,
		Action: urfave.Action(lc.action, logger),
	}
}

func (lc *initCommand) action(ctx context.Context, cmd *cli.Command, logger *slogutil.Logger) error {
	ctrl := initcmd.NewController(afero.NewOsFs())
	if err := logger.SetLevel(cmd.String("log-level")); err != nil {
		return fmt.Errorf("set log level: %w", err)
	}
	return ctrl.Init(ctx, logger.Logger) //nolint:wrapcheck
}
