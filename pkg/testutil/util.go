package testutil

import (
	"fmt"
	"os"

	"github.com/lintnet/lintnet/pkg/osfile"
	"github.com/spf13/afero"
)

func NewFs(files map[string]string, dirs ...string) (afero.Fs, error) {
	fs := afero.NewMemMapFs()
	for _, dir := range dirs {
		if err := osfile.MkdirAll(fs, dir); err != nil {
			return nil, err //nolint:wrapcheck
		}
	}
	for name, body := range files {
		if err := afero.WriteFile(fs, name, []byte(body), osfile.FilePermission); err != nil {
			return nil, err //nolint:wrapcheck
		}
	}
	return fs, nil
}

func ReadFs(files map[string]string, dirs ...string) (afero.Fs, error) {
	m := make(map[string]string, len(files))
	for k, v := range files {
		b, err := os.ReadFile(v)
		if err != nil {
			return nil, fmt.Errorf("read a test data file: %w", err)
		}
		m[k] = string(b)
	}
	return NewFs(m, dirs...)
}
