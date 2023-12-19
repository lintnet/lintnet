package lint

import (
	"errors"
	"fmt"
	"os"

	"github.com/lintnet/lintnet/pkg/config"
	"gopkg.in/yaml.v3"
)

func (c *Controller) readConfig(p string, cfg *config.Config) error {
	cfgFile, err := c.fs.Open(p)
	if err != nil {
		return fmt.Errorf("open a config file: %w", err)
	}
	defer cfgFile.Close()
	if err := yaml.NewDecoder(cfgFile).Decode(cfg); err != nil {
		return fmt.Errorf("read a config file as YAML: %w", err)
	}
	return nil
}

func (c *Controller) findAndReadConfig(p string, cfg *config.Config) error {
	if p != "" {
		if err := c.readConfig(p, cfg); err != nil {
			return fmt.Errorf("read a config file: %w", err)
		}
		return nil
	}
	paths := []string{
		"lintnet.yaml",
		"lintnet.yml",
		".lintnet.yaml",
		".lintnet.yml",
	}
	for _, p := range paths {
		if err := c.readConfig(p, cfg); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return fmt.Errorf("read a config file: %w", err)
		}
		return nil
	}
	return fmt.Errorf("config file isn't found: %w", os.ErrNotExist)
}
