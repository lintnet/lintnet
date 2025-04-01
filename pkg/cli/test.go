package cli

import (
	"fmt"
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
	"github.com/urfave/cli/v3"
)

type testCommand struct {
	logE    *logrus.Entry
	version string
}

func (tc *testCommand) command() *cli.Command {
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
		Action: tc.action,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "target",
				Aliases: []string{"t"},
				Usage:   "Target ID",
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
	ctrl := testcmd.NewController(param, fs, os.Stdout, importer)
	pwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get the current directory: %w", err)
	}
	return ctrl.Test(c.Context, logE, &testcmd.ParamTest{ //nolint:wrapcheck
		FilePaths:      c.Args().Slice(),
		ConfigFilePath: c.String("config"),
		TargetID:       c.String("target"),
		RootDir:        rootDir,
		PWD:            pwd,
	})
}
