package cli

import (
	"fmt"
	"os"
	"runtime"

	"github.com/lintnet/lintnet/pkg/config"
	"github.com/lintnet/lintnet/pkg/controller/info"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
	"github.com/urfave/cli/v3"
)

type infoCommand struct {
	logE    *logrus.Entry
	version string
	commit  string
}

func (lc *infoCommand) command() *cli.Command {
	return &cli.Command{
		Name:      "info",
		Usage:     "Output the information regarding lintnet",
		UsageText: "lintnet info [command options]",
		Description: `Output the information regarding lintnet.

$ lintnet info
{
  "version": "v0.3.0",
  "config_file": "lintnet.jsonnet",
  "root_dir": "/Users/foo/Library/Application Support/lintnet",
  "data_root_dir": "/Users/foo/repos/src/github.com/lintnet/lintnet",
  "env": {
	"GITHUB_TOKEN": "(masked)",
	"LINTNET_LOG_LEVEL": "warn"
  }
}

This command is useful for trouble shooting and sharing your environment in GitHub Issues.

You can mask the current user name.

$ lintnet info -mask-user

You can also get the root directory where modules are installed.

$ lintnet info -module-root-dir
`,
		Action: lc.action,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "module-root-dir",
				Usage: "Show only the root directory where modules are installed",
			},
			&cli.BoolFlag{
				Name:  "mask-user",
				Usage: "Mask the current user name",
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
		Commit:  lc.commit,
		Env:     fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
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
		ModuleRootDir:  c.Bool("module-root-dir"),
		MaskUser:       c.Bool("mask-user"),
	})
}
