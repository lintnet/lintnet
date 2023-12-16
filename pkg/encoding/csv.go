package encoding

import (
	"encoding/csv"
	"fmt"
	"io"
)

type csvDecoder struct {
	reader *csv.Reader
}

func newCSVDecoder(r io.Reader) Decoder {
	return &csvDecoder{
		reader: csv.NewReader(r),
	}
}

func newTSVDecoder(r io.Reader) Decoder {
	reader := csv.NewReader(r)
	reader.Comma = '	'
	return &csvDecoder{
		reader: reader,
	}
}

func (c *csvDecoder) Decode() (interface{}, error) {
	records, err := c.reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("parse a file as CSV: %w", err)
	}
	return records, nil
}
