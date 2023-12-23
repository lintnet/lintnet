package initcmd

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

//go:embed init.yaml
var cfgTemplate []byte

func (c *Controller) Init(_ context.Context, _ *logrus.Entry) error {
	paths := []string{
		"lintnet.yaml",
		"lintnet.yml",
		".lintnet.yaml",
		".lintnet.yml",
	}
	for _, path := range paths {
		f, err := afero.Exists(c.fs, path)
		if err != nil {
			return fmt.Errorf("check if a configuration file exists: %w", err)
		}
		if f {
			return nil
		}
	}
	f, err := c.fs.Create("lintnet.yaml")
	if err != nil {
		return fmt.Errorf("create a configuration file: %w", err)
	}
	defer f.Close()
	if _, err := f.Write(cfgTemplate); err != nil {
		return fmt.Errorf("write a configuration file: %w", err)
	}
	return nil
}
