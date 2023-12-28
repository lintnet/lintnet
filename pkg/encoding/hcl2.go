package encoding

import (
	"encoding/json"
	"fmt"

	"github.com/tmccombs/hcl2json/convert"
)

type hcl2Unmarshaler struct{}

func (d *hcl2Unmarshaler) Unmarshal(b []byte) (any, error) {
	hclBytes, err := convert.Bytes(b, "", convert.Options{})
	if err != nil {
		return nil, fmt.Errorf("convert hcl2 to JSON: %w", err)
	}
	var dest any
	if err := json.Unmarshal(hclBytes, &dest); err != nil {
		return nil, fmt.Errorf("unmarshal hcl2: %w", err)
	}
	return dest, nil
}
