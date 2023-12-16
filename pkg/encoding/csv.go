package encoding

import (
	"encoding/csv"
	"fmt"
	"strings"
)

type csvUnmarshaler struct {
	TSV bool
}

func (c *csvUnmarshaler) Unmarshal(b []byte) (interface{}, error) {
	reader := csv.NewReader(strings.NewReader(string(b)))
	if c.TSV {
		reader.Comma = '	'
	}
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("parse a file as CSV: %w", err)
	}
	return records, nil
}
