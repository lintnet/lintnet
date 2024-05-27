package encoding

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/suzuki-shunsuke/go-convmap/convmap"
	"gopkg.in/yaml.v3"
)

type yamlUnmarshaler struct{}

func (d *yamlUnmarshaler) Unmarshal(b []byte) (any, error) {
	// Treat YAML as multiple documents
	var arr []any
	dec := yaml.NewDecoder(bytes.NewReader(b))
	for {
		var value any
		err := dec.Decode(&value)
		if errors.Is(err, io.EOF) {
			return arr, nil
		}
		if err != nil {
			return nil, fmt.Errorf("decode YAML: %w", err)
		}

		// https://github.com/lintnet/lintnet/issues/437
		v, err := convmap.Convert(value, nil)
		if err != nil {
			return nil, fmt.Errorf("convert map keys to string: %w", err)
		}

		arr = append(arr, v)
	}
}
