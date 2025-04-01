package cli

import (
	"fmt"
	"net/http"
	"os"
	"runtime"

	"github.com/lintnet/lintnet/pkg/config"
	"github.com/lintnet/lintnet/pkg/controller/lint"
	"github.com/lintnet/lintnet/pkg/github"
	"github.com/lintnet/lintnet/pkg/jsonnet"
	"github.com/lintnet/lintnet/pkg/log"
	"github.com/lintnet/lintnet/pkg/module"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
	"github.com/urfave/cli/v3"
)

type lintCommand struct {
	logE    *logrus.Entry
	version string
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
				EnvVars: []string{"LINTNET_ERROR_LEVEL"},
			},
			&cli.StringFlag{
				Name:    "shown-error-level",
				Usage:   "Set the shown error level",
				EnvVars: []string{"LINTNET_SHOWN_ERROR_LEVEL"},
			},
			&cli.BoolFlag{
				Name:    "output-success",
				Usage:   "Output the result even if the lint succeeds",
				EnvVars: []string{"LINTNET_OUTPUT_SUCCESS"},
			},
		},
	}
}

func (lc *lintCommand) action(c *cli.Context) error {
	fs := afero.NewOsFs()
	logE := lc.logE
	log.SetLevel(c.String("log-level"), logE)
	log.SetColor(c.String("log-color"), logE)
	rootDir := os.Getenv("LINTNET_ROOT_DIR")
	if rootDir == "" {
		dir, err := config.GetRootDir()
		if err != nil {
			logerr.WithError(logE, err).Warn("get the root directory")
		}
		rootDir = dir
	}
	modInstaller := module.NewInstaller(fs, github.New(c.Context), http.DefaultClient)
	importer := jsonnet.NewImporter(c.Context, logE, &module.ParamInstall{
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
	return ctrl.Lint(c.Context, logE, &lint.ParamLint{ //nolint:wrapcheck
		FilePaths:       c.Args().Slice(),
		ErrorLevel:      c.String("error-level"),
		ShownErrorLevel: c.String("shown-error-level"),
		ConfigFilePath:  c.String("config"),
		TargetID:        c.String("target"),
		OutputSuccess:   c.Bool("output-success"),
		Output:          c.String("output"),
		RootDir:         rootDir,
		DataRootDir:     pwd,
		PWD:             pwd,
	})
}
