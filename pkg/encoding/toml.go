package encoding

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type tomlUnmarshaler struct{}

func (d *tomlUnmarshaler) Unmarshal(b []byte) (any, error) {
	var v any
	if err := toml.Unmarshal(b, &v); err != nil {
		return nil, fmt.Errorf("parse a file as TOML: %w", err)
	}
	return v, nil
}
