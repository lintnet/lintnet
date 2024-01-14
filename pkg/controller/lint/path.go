package lint

import (
	"path/filepath"
)

// slash, path
// rel, abs
// glob, path

func Abs(base, p string) string {
	if filepath.IsAbs(p) {
		return p
	}
	return filepath.Join(base, p)
}
