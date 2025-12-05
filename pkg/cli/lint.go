package cli

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"runtime"

	"github.com/lintnet/lintnet/pkg/config"
	"github.com/lintnet/lintnet/pkg/controller/lint"
	"github.com/lintnet/lintnet/pkg/github"
	"github.com/lintnet/lintnet/pkg/jsonnet"
	"github.com/lintnet/lintnet/pkg/module"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/slog-error/slogerr"
	"github.com/suzuki-shunsuke/slog-util/slogutil"
	"github.com/urfave/cli/v3"
)

type lintCommand struct {
	logger      *slog.Logger
	logLevelVar *slog.LevelVar
	version     string
}

func (lc *lintCommand) command() *cli.Command {
	return &cli.Command{
		Name:      "lint",
		Aliases:   []string{"l"},
		Usage:     "Lint files",
		UsageText: "lintnet lint [command options] [lint file paths and data file paths]",
		Description: `Lint files

$ lintnet lint

You can lint only specific files.

$ lintnet lint [lint file paths and data file paths]

You can also lint only a specific target.

$ lintnet lint -target [target id]

By default, lintnet outputs nothing when the lint succeeds.
You can output JSON even if the lint succeeds. This is useful if you pass the output to other program such as jq.

$ lintnet lint -output-success
`,
		Action: lc.action,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "You can customize the output format. You can specify an output id",
			},
			&cli.StringFlag{
				Name:    "target",
				Aliases: []string{"t"},
				Usage:   "Lint only a specific target. You can specify a target id",
			},
			&cli.StringFlag{
				Name:    "error-level",
				Aliases: []string{"e"},
				Usage:   "Set the error level",
				Sources: cli.EnvVars("LINTNET_ERROR_LEVEL"),
			},
			&cli.StringFlag{
				Name:    "shown-error-level",
				Usage:   "Set the shown error level",
				Sources: cli.EnvVars("LINTNET_SHOWN_ERROR_LEVEL"),
			},
			&cli.BoolFlag{
				Name:    "output-success",
				Usage:   "Output the result even if the lint succeeds",
				Sources: cli.EnvVars("LINTNET_OUTPUT_SUCCESS"),
			},
		},
	}
}

func (lc *lintCommand) action(ctx context.Context, cmd *cli.Command) error {
	fs := afero.NewOsFs()
	logger := lc.logger
	if err := slogutil.SetLevel(lc.logLevelVar, cmd.String("log-level")); err != nil {
		return fmt.Errorf("set log level: %w", err)
	}
	rootDir := os.Getenv("LINTNET_ROOT_DIR")
	if rootDir == "" {
		dir, err := config.GetRootDir()
		if err != nil {
			slogerr.WithError(logger, err).Warn("get the root directory")
		}
		rootDir = dir
	}
	modInstaller := module.NewInstaller(fs, github.New(ctx), http.DefaultClient)
	importer := jsonnet.NewImporter(ctx, logger, &module.ParamInstall{
		BaseDir: rootDir,
	}, &jsonnet.FileImporter{
		JPaths: []string{rootDir},
	}, modInstaller)
	param := &lint.ParamController{
		Version: lc.version,
		Env:     fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
	ctrl := lint.NewController(param, fs, os.Stdout, modInstaller, importer)
	pwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get the current directory: %w", err)
	}
	return ctrl.Lint(ctx, logger, &lint.ParamLint{ //nolint:wrapcheck
		FilePaths:       cmd.Args().Slice(),
		ErrorLevel:      cmd.String("error-level"),
		ShownErrorLevel: cmd.String("shown-error-level"),
		ConfigFilePath:  cmd.String("config"),
		TargetID:        cmd.String("target"),
		OutputSuccess:   cmd.Bool("output-success"),
		Output:          cmd.String("output"),
		RootDir:         rootDir,
		DataRootDir:     pwd,
		PWD:             pwd,
	})
}
