package cli

import (
	"fmt"
	"os"

	"github.com/lintnet/lintnet/pkg/config"
	"github.com/lintnet/lintnet/pkg/controller/info"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
	"github.com/urfave/cli/v2"
)

type infoCommand struct {
	logE    *logrus.Entry
	version string
}

func (lc *infoCommand) command() *cli.Command {
	return &cli.Command{
		Name:   "info",
		Usage:  "Show information",
		Action: lc.action,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name: "module-root-dir",
			},
			&cli.BoolFlag{
				Name: "mask-user",
			},
		},
	}
}

func (lc *infoCommand) action(c *cli.Context) error {
	fs := afero.NewOsFs()
	logE := lc.logE
	rootDir := os.Getenv("LINTNET_ROOT_DIR")
	if rootDir == "" {
		dir, err := config.GetRootDir()
		if err != nil {
			logerr.WithError(logE, err).Warn("get the root directory")
		}
		rootDir = dir
	}
	param := &info.ParamController{
		Version: lc.version,
	}
	ctrl := info.NewController(param, fs, os.Stdout)
	pwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get the current directory: %w", err)
	}
	return ctrl.Info(c.Context, &info.ParamInfo{ //nolint:wrapcheck
		ConfigFilePath: c.String("config"),
		RootDir:        rootDir,
		PWD:            pwd,
		DataRootDir:    pwd,
		ModuleRootDir:  c.Bool("module-root-dir"),
		MaskUser:       c.Bool("mask-user"),
	})
}
