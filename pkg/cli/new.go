package cli

import (
	"context"
	"fmt"

	"github.com/lintnet/lintnet/pkg/controller/newcmd"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/slog-util/slogutil"
	"github.com/urfave/cli/v3"
)

type newCommand struct{}

type NewArgs struct {
	*GlobalFlags

	FileName string
}

func (lc *newCommand) command(logger *slogutil.Logger, gFlags *GlobalFlags) *cli.Command {
	args := &NewArgs{
		GlobalFlags: gFlags,
	}
	return &cli.Command{
		Name:      "new",
		Usage:     "Create a lint file and a test file",
		UsageText: "lintnet new [<lint file|main.jsonnet>]",
		Description: `Create a lint file and a test file.

$ lintnet new [<lint file|main.jsonnet>]

This command creates a lint file and a test file.
If the argument is not given, the lint file is created as "main.jsonnet".
`,
		Action: func(ctx context.Context, _ *cli.Command) error {
			return lc.action(ctx, logger, args)
		},
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name:        "file",
				Value:       "main.jsonnet",
				Destination: &args.FileName,
			},
		},
	}
}

func (lc *newCommand) action(ctx context.Context, logger *slogutil.Logger, args *NewArgs) error {
	ctrl := newcmd.NewController(afero.NewOsFs())
	if err := logger.SetLevel(args.LogLevel); err != nil {
		return fmt.Errorf("set log level: %w", err)
	}
	return ctrl.New(ctx, logger.Logger, args.FileName) //nolint:wrapcheck
}
