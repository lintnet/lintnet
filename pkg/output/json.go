package output

import (
	"encoding/json"
	"fmt"
	"io"
)

type jsonOutputter struct {
	stdout io.Writer
}

func (o *jsonOutputter) Output(result *Output) error {
	return outputJSON(o.stdout, result)
}

func outputJSON(w io.Writer, result any) error {
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(result); err != nil {
		return fmt.Errorf("encode the result as JSON: %w", err)
	}
	return nil
}
