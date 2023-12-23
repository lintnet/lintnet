package cli

import (
	"net/http"
	"os"

	"github.com/lintnet/lintnet/pkg/config"
	"github.com/lintnet/lintnet/pkg/controller/lint"
	"github.com/lintnet/lintnet/pkg/github"
	"github.com/lintnet/lintnet/pkg/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
	"github.com/urfave/cli/v2"
)

type lintCommand struct {
	logE *logrus.Entry
}

func (lc *lintCommand) command() *cli.Command {
	return &cli.Command{
		Name:   "lint",
		Usage:  "Lint files",
		Action: lc.action,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "rule-base-dir",
				Aliases: []string{"d"},
			},
			&cli.StringFlag{
				Name:    "error-level",
				Aliases: []string{"e"},
				EnvVars: []string{"LINTNET_ERROR_LEVEL"},
			},
			&cli.BoolFlag{
				Name:    "output-success",
				EnvVars: []string{"LINTNET_OUTPUT_SUCCESS"},
			},
		},
	}
}

func (lc *lintCommand) action(c *cli.Context) error {
	gh := github.New(c.Context)
	ctrl := lint.NewController(afero.NewOsFs(), os.Stdout, gh, http.DefaultClient)
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
	return ctrl.Lint(c.Context, logE, &lint.ParamLint{ //nolint:wrapcheck
		FilePaths:      c.Args().Slice(),
		RuleBaseDir:    c.String("rule-base-dir"),
		ErrorLevel:     c.String("error-level"),
		ConfigFilePath: c.String("config"),
		OutputSuccess:  c.Bool("output-success"),
		RootDir:        rootDir,
	})
}
