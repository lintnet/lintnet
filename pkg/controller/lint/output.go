package lint

import (
	"encoding/json"
	"errors"
	"fmt"
)

func (c *Controller) Output(results map[string]*FileResult) error {
	if !checkFailed(results) {
		return nil
	}
	encoder := json.NewEncoder(c.stdout)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(results); err != nil {
		return fmt.Errorf("encode the result as JSON: %w", err)
	}
	return errors.New("lint failed")
}
