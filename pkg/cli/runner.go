package cli

import (
	"context"
	"io"
	"log/slog"

	"github.com/suzuki-shunsuke/urfave-cli-v3-util/helpall"
	"github.com/suzuki-shunsuke/urfave-cli-v3-util/vcmd"
	"github.com/urfave/cli/v3"
)

type Runner struct {
	Stdin       io.Reader
	Stdout      io.Writer
	Stderr      io.Writer
	LDFlags     *LDFlags
	Logger      *slog.Logger
	LogLevelVar *slog.LevelVar
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
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Configuration file path",
				Sources: cli.EnvVars("LINTNET_CONFIG"),
			},
		},
		EnableShellCompletion: true,
		Commands: []*cli.Command{
			(&lintCommand{
				logger:      r.Logger,
				logLevelVar: r.LogLevelVar,
				version:     r.LDFlags.Version,
			}).command(),
			(&infoCommand{
				logger:  r.Logger,
				version: r.LDFlags.Version,
				commit:  r.LDFlags.Commit,
			}).command(),
			(&initCommand{
				logger:      r.Logger,
				logLevelVar: r.LogLevelVar,
			}).command(),
			(&testCommand{
				logger:      r.Logger,
				logLevelVar: r.LogLevelVar,
				version:     r.LDFlags.Version,
			}).command(),
			(&newCommand{
				logger:      r.Logger,
				logLevelVar: r.LogLevelVar,
			}).command(),
			(&completionCommand{
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
