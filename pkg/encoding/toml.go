package encoding

import (
	"fmt"
	"io"

	"github.com/BurntSushi/toml"
)

type tomlDecoder struct {
	decoder *toml.Decoder
}

func (d *tomlDecoder) Decode() (interface{}, error) {
	var v interface{}
	_, err := d.decoder.Decode(&v)
	if err != nil {
		return nil, fmt.Errorf("parse a file as TOML: %w", err)
	}
	return v, nil
}

func newTOMLDecoder(r io.Reader) Decoder {
	return &tomlDecoder{
		decoder: toml.NewDecoder(r),
	}
}
