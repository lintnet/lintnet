package encoding

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type yamlUnmarshaler struct{}

func (d *yamlUnmarshaler) Unmarshal(b []byte) (any, error) {
	var dest any
	if err := yaml.Unmarshal(b, &dest); err != nil {
		return nil, fmt.Errorf("parse a file as YAML: %w", err)
	}
	return dest, nil
}
