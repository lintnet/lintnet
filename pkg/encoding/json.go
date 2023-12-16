package encoding

import (
	"encoding/json"
	"fmt"
	"io"
)

type jsonDecoder struct {
	decoder *json.Decoder
}

func (d *jsonDecoder) Decode() (interface{}, error) {
	var dest interface{}
	if err := d.decoder.Decode(&dest); err != nil {
		return nil, fmt.Errorf("parse a file as JSON: %w", err)
	}
	return dest, nil
}

func newJSONDecoder(r io.Reader) Decoder {
	return &jsonDecoder{
		decoder: json.NewDecoder(r),
	}
}
