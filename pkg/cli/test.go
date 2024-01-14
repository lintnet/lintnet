package cli

import (
	"net/http"
	"os"

	"github.com/lintnet/lintnet/pkg/config"
	"github.com/lintnet/lintnet/pkg/controller/testcmd"
	"github.com/lintnet/lintnet/pkg/github"
	"github.com/lintnet/lintnet/pkg/jsonnet"
	"github.com/lintnet/lintnet/pkg/log"
	"github.com/lintnet/lintnet/pkg/module"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
	"github.com/urfave/cli/v2"
)

type testCommand struct {
	logE    *logrus.Entry
	version string
}

func (tc *testCommand) command() *cli.Command {
	return &cli.Command{
		Name:   "test",
		Usage:  "Test lint files",
		Action: tc.action,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "output-success",
				EnvVars: []string{"LINTNET_OUTPUT_SUCCESS"},
			},
		},
	}
}

func (tc *testCommand) action(c *cli.Context) error {
	fs := afero.NewOsFs()
	logE := tc.logE
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
	param := &testcmd.ParamController{
		Version: tc.version,
	}
	ctrl := testcmd.NewTestController(param, fs, os.Stdout, importer)
	return ctrl.Test(c.Context, logE, &testcmd.ParamTest{ //nolint:wrapcheck
		FilePaths:      c.Args().Slice(),
		ConfigFilePath: c.String("config"),
		RootDir:        rootDir,
	})
}
