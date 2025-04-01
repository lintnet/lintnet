package cli

import (
	"github.com/urfave/cli/v3"
)

type versionCommand struct{}

func (vc *versionCommand) command() *cli.Command {
	return &cli.Command{
		Name:   "version",
		Usage:  "Show version",
		Action: vc.action,
	}
}

func (vc *versionCommand) action(c *cli.Context) error {
	cli.ShowVersion(c)
	return nil
}
