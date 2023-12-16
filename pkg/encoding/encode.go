package encoding

import (
	"path/filepath"
)

type Unmarshaler interface {
	Unmarshal(b []byte) (interface{}, error)
}

func NewUnmarshaler(fileName string) (Unmarshaler, string, error) {
	ext := filepath.Ext(fileName)
	switch ext {
	case ".csv":
		return &csvUnmarshaler{}, "csv", nil
	case ".json":
		return &jsonUnmarshaler{}, "json", nil
	case ".toml":
		return &tomlUnmarshaler{}, "toml", nil
	case ".tsv":
		return &csvUnmarshaler{
			TSV: true,
		}, "tsv", nil
	case ".yml", ".yaml":
		return &yamlUnmarshaler{}, "yaml", nil
	default:
		return &plainUnmarshaler{}, "plain_text", nil
	}
}
