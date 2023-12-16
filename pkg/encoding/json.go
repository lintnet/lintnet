package encoding

import (
	"encoding/json"
	"fmt"
)

type jsonUnmarshaler struct{}

func (d *jsonUnmarshaler) Unmarshal(b []byte) (interface{}, error) {
	var dest interface{}
	if err := json.Unmarshal(b, &dest); err != nil {
		return nil, fmt.Errorf("parse a file as JSON: %w", err)
	}
	return dest, nil
}
