package cli

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/lintnet/lintnet/pkg/config"
	"github.com/lintnet/lintnet/pkg/controller/testcmd"
	"github.com/lintnet/lintnet/pkg/github"
	"github.com/lintnet/lintnet/pkg/jsonnet"
	"github.com/lintnet/lintnet/pkg/module"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/slog-error/slogerr"
	"github.com/suzuki-shunsuke/slog-util/slogutil"
	"github.com/suzuki-shunsuke/urfave-cli-v3-util/urfave"
	"github.com/urfave/cli/v3"
)

type testCommand struct {
	version string
}

func (tc *testCommand) command(logger *slogutil.Logger) *cli.Command {
	return &cli.Command{
		Name:      "test",
		Aliases:   []string{"t"},
		Usage:     "Test lint files",
		ArgsUsage: "[<lint file, test file, or directory> ...]",
		Description: `Test lint files.

If you run "lintnet test" without any argument,
lintnet searches lint files using a configuration file and tests all lint files having test files.
Lint files without test files are ignored.
You can test only specific files by specifying files as arguments.
If you specify files explicitly, a configuration file is unnecessary.
This means when you develop modules, you don't have to prepare a configuration file.
If you specify directories, lint files in those directories and subdirectories are tested.
For example, "lintnet test ." searches files matching the glob pattern "**/*.jsonnet",
and "lintnet test foo" search files matching "foo/**/*.jsonnet".
If a configuration file isn't specified and isn't found, "lintnet test" works as "lintnet test .".

You can test only a specific target with -target option.
`,
		Action: urfave.Action(tc.action, logger),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "target",
				Aliases: []string{"t"},
				Usage:   "Target ID",
			},
		},
	}
}

func (tc *testCommand) action(ctx context.Context, cmd *cli.Command, logger *slogutil.Logger) error {
	fs := afero.NewOsFs()
	if err := logger.SetLevel(cmd.String("log-level")); err != nil {
		return fmt.Errorf("set log level: %w", err)
	}
	rootDir := os.Getenv("LINTNET_ROOT_DIR")
	if rootDir == "" {
		dir, err := config.GetRootDir()
		if err != nil {
			slogerr.WithError(logger.Logger, err).Warn("get the root directory")
		}
		rootDir = dir
	}
	modInstaller := module.NewInstaller(fs, github.New(ctx), http.DefaultClient)
	importer := jsonnet.NewImporter(ctx, logger.Logger, &module.ParamInstall{
		BaseDir: rootDir,
	}, &jsonnet.FileImporter{
		JPaths: []string{rootDir},
	}, modInstaller)
	param := &testcmd.ParamController{
		Version: tc.version,
	}
	ctrl := testcmd.NewController(param, fs, os.Stdout, importer)
	pwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get the current directory: %w", err)
	}
	return ctrl.Test(ctx, logger.Logger, &testcmd.ParamTest{ //nolint:wrapcheck
		FilePaths:      cmd.Args().Slice(),
		ConfigFilePath: cmd.String("config"),
		TargetID:       cmd.String("target"),
		RootDir:        rootDir,
		PWD:            pwd,
	})
}
