package lint

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/lintnet/lintnet/pkg/config"
	"github.com/lintnet/lintnet/pkg/jsonnet"
)

func (c *Controller) readConfig(p string, cfg *config.Config) error {
	vm := jsonnet.NewVM("{}", c.importer)
	node, err := jsonnet.Read(c.fs, p)
	if err != nil {
		return fmt.Errorf("parse a configuration file as Jsonnet: %w", err)
	}
	result, err := vm.Evaluate(node)
	if err != nil {
		return fmt.Errorf("evaluate a configuration file as Jsonnet: %w", err)
	}
	if err := json.Unmarshal([]byte(result), cfg); err != nil {
		return fmt.Errorf("unmarshal config as JSON: %w", err)
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
		return nil
	}
	return fmt.Errorf("config file isn't found: %w", os.ErrNotExist)
}
