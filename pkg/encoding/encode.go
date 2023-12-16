package encoding

import (
	"errors"
	"io"
	"path/filepath"
)

type (
	NewDecoder func(io.Reader) Decoder
	Decoder    interface {
		Decode() (interface{}, error)
	}
)

func GetNewDecoder(fileName string) (NewDecoder, string, error) {
	ext := filepath.Ext(fileName)
	switch ext {
	case ".csv":
		return newCSVDecoder, "csv", nil
	case ".json":
		return newJSONDecoder, "json", nil
	case ".toml":
		return newTOMLDecoder, "toml", nil
	case ".tsv":
		return newTSVDecoder, "tsv", nil
	case ".yml", ".yaml":
		return newYAMLDecoder, "yaml", nil
	default:
		return nil, "", errors.New("this format is unsupported")
	}
}
