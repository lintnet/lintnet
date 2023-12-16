package encoding

import (
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
)

type yamlDecoder struct {
	decoder *yaml.Decoder
}

func (d *yamlDecoder) Decode() (interface{}, error) {
	var dest interface{}
	if err := d.decoder.Decode(&dest); err != nil {
		return nil, fmt.Errorf("parse a file as YAML: %w", err)
	}
	return dest, nil
}

func newYAMLDecoder(r io.Reader) Decoder {
	return &yamlDecoder{
		decoder: yaml.NewDecoder(r),
	}
}
