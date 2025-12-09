package cli

import (
	"context"
	"fmt"

	"github.com/lintnet/lintnet/pkg/controller/initcmd"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/slog-util/slogutil"
	"github.com/urfave/cli/v3"
)

type initCommand struct{}

type InitArgs struct {
	*GlobalFlags
}

func (lc *initCommand) command(logger *slogutil.Logger, gFlags *GlobalFlags) *cli.Command {
	args := &InitArgs{
		GlobalFlags: gFlags,
	}
	return &cli.Command{
		Name:      "init",
		Usage:     "Scaffold configuration file",
		UsageText: "lintnet init",
		Description: `Scaffold configuration file.

$ lintnet init

This command generates lintnet.jsonnet.
If the file already exists, this command does nothing.
`,
		Action: func(ctx context.Context, _ *cli.Command) error {
			return lc.action(ctx, logger, args)
		},
	}
}

func (lc *initCommand) action(ctx context.Context, logger *slogutil.Logger, args *InitArgs) error {
	ctrl := initcmd.NewController(afero.NewOsFs())
	if err := logger.SetLevel(args.LogLevel); err != nil {
		return fmt.Errorf("set log level: %w", err)
	}
	return ctrl.Init(ctx, logger.Logger) //nolint:wrapcheck
}
