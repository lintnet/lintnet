package config

import (
	"path/filepath"
)

type DataFile struct {
	Path     string
	Excluded bool
}

func NewDataFile(s string) *DataFile {
	a, excluded := parseNegationOperator(s)
	return &DataFile{
		Path:     filepath.Clean(a),
		Excluded: excluded,
	}
}
