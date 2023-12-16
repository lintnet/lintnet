package lint

import (
	"encoding/json"
	"errors"
	"fmt"
)

func (c *Controller) Output(results map[string]*FileResult) error {
	if !isFailed(results) {
		return nil
	}
	return c.outputJSON(c.formatResultToOutput(results))
}

type FlatError struct {
	DataFilePath string      `json:"data_file,omitempty"`
	LintFilePath string      `json:"lint_file,omitempty"`
	Location     interface{} `json:"location,omitempty"`
	RuleName     string      `json:"rule,omitempty"`
	Error        string      `json:"error,omitempty"`
}

func (c *Controller) formatResultToOutput(results map[string]*FileResult) []*FlatError {
	list := make([]*FlatError, 0, len(results))
	for dataFilePath, fileResult := range results {
		list = append(list, fileResult.flattenError(dataFilePath)...)
	}
	if len(list) == 0 {
		return nil
	}
	return list
}

func (c *Controller) outputJSON(result interface{}) error {
	encoder := json.NewEncoder(c.stdout)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(result); err != nil {
		return fmt.Errorf("encode the result as JSON: %w", err)
	}
	return errors.New("lint failed")
}
