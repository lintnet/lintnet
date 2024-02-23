package config

import (
	"path/filepath"
	"strings"
)

type DataFile struct {
	Path     string
	Excluded bool
}

func NewDataFile(s string) *DataFile {
	a := strings.TrimPrefix(s, "!")
	return &DataFile{
		Path:     filepath.Clean(a),
		Excluded: a != s,
	}
}
