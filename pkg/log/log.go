package log

import (
	"bytes"
	"encoding/json"
)

func JSON(data any) any {
	return &jsonData{
		data: data,
	}
}

type jsonData struct {
	data any
}

func (j *jsonData) String() string {
	buf := &bytes.Buffer{}
	encoder := json.NewEncoder(buf)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(j.data); err != nil {
		return err.Error()
	}
	return buf.String()
}
