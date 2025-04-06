package cli

import (
	"context"
	"io"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/urfave-cli-v3-util/helpall"
	"github.com/suzuki-shunsuke/urfave-cli-v3-util/vcmd"
	"github.com/urfave/cli/v3"
)

type Runner struct {
	Stdin   io.Reader
	Stdout  io.Writer
	Stderr  io.Writer
	LDFlags *LDFlags
	LogE    *logrus.Entry
}

type LDFlags struct {
	Version string
	Commit  string
	Date    string
}

func (r *Runner) Run(ctx context.Context, args ...string) error {
	return helpall.With(&cli.Command{ //nolint:wrapcheck
		Name:    "lintnet",
		Usage:   "Powerful, Secure, Shareable Linter Powered by Jsonnet. https://lintnet.github.io/",
		Version: r.LDFlags.Version + " (" + r.LDFlags.Commit + ")",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "log-level",
				Usage:   "log level",
				Sources: cli.EnvVars("LINTNET_LOG_LEVEL"),
			},
			&cli.StringFlag{
				Name:    "log-color",
				Usage:   "Log color. One of 'auto' (default), 'always', 'never'",
				Sources: cli.EnvVars("LINTNET_LOG_COLOR"),
			},
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Configuration file path",
				Sources: cli.EnvVars("LINTNET_CONFIG"),
			},
		},
		EnableShellCompletion: true,
		Commands: []*cli.Command{
			(&lintCommand{
				logE:    r.LogE,
				version: r.LDFlags.Version,
			}).command(),
			(&infoCommand{
				logE:    r.LogE,
				version: r.LDFlags.Version,
				commit:  r.LDFlags.Commit,
			}).command(),
			(&initCommand{
				logE: r.LogE,
			}).command(),
			(&testCommand{
				logE:    r.LogE,
				version: r.LDFlags.Version,
			}).command(),
			(&newCommand{
				logE: r.LogE,
			}).command(),
			(&completionCommand{
				logE:   r.LogE,
				stdout: r.Stdout,
			}).command(),
			vcmd.New(&vcmd.Command{
				Name:    "lintnet",
				Version: r.LDFlags.Version,
				SHA:     r.LDFlags.Commit,
			}),
		},
	}, nil).Run(ctx, args)
}
