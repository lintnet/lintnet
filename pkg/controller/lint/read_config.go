package lint

import (
	"errors"
	"fmt"
	"os"

	"github.com/lintnet/lintnet/pkg/config"
	"github.com/lintnet/lintnet/pkg/jsonnet"
	"github.com/spf13/afero"
)

type ConfigReader struct {
	fs       afero.Fs
	importer *jsonnet.Importer
}

func (r *ConfigReader) read(p string, cfg *config.RawConfig) error {
	return jsonnet.Read(r.fs, p, "{}", r.importer, cfg) //nolint:wrapcheck
}

func (r *ConfigReader) Read(p string, cfg *config.RawConfig) error {
	if p != "" {
		if err := r.read(p, cfg); err != nil {
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
		if err := r.read(p, cfg); err != nil {
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
