package lint

import (
	"errors"
	"fmt"
	"os"

	"github.com/lintnet/lintnet/pkg/config"
	"github.com/lintnet/lintnet/pkg/jsonnet"
)

func (c *Controller) readConfig(p string, cfg *config.RawConfig) error {
	return jsonnet.Read(c.fs, p, "{}", c.importer, cfg) //nolint:wrapcheck
}

func (c *Controller) findAndReadConfig(p string, cfg *config.RawConfig) error {
	if p != "" {
		if err := c.readConfig(p, cfg); err != nil {
			return fmt.Errorf("read a config file: %w", err)
		}
		cfg.FilePath = p
		return nil
	}
	paths := []string{
		"lintnet.jsonnet",
		".lintnet.jsonnet",
	}
	for _, p := range paths {
		if err := c.readConfig(p, cfg); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return fmt.Errorf("read a config file: %w", err)
		}
		cfg.FilePath = p
		return nil
	}
	return fmt.Errorf("config file isn't found: %w", os.ErrNotExist)
}
