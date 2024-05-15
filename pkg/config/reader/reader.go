package reader

import (
	"errors"
	"fmt"
	"os"

	gojsonnet "github.com/google/go-jsonnet"
	"github.com/lintnet/lintnet/pkg/config"
	"github.com/lintnet/lintnet/pkg/jsonnet"
	"github.com/spf13/afero"
)

type Reader struct {
	fs       afero.Fs
	importer gojsonnet.Importer
}

func New(fs afero.Fs, importer gojsonnet.Importer) *Reader {
	return &Reader{
		fs:       fs,
		importer: importer,
	}
}

func (r *Reader) read(p string, cfg *config.RawConfig) error {
	return jsonnet.Read(r.fs, p, "{}", r.importer, cfg) //nolint:wrapcheck
}

// Read finds and reads a configuration file.
// A configuration file is a Jsonnet file.
// p is a configuration file path.
// If p is an empty string, it tries to find a configuration file from the default paths.
// The default paths are "lintnet.jsonnet" and ".lintnet.jsonnet".
// The configuration file is loaded to cfg.
// If it fails to find, read, or parse a configuration file, it returns an error.
func (r *Reader) Read(p string, cfg *config.RawConfig) error {
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
