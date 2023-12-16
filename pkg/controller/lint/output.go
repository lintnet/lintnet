package lint

import (
	"encoding/json"
	"errors"
	"fmt"
)

func (c *Controller) Output(logLevel ErrorLevel, results map[string]*FileResult) error {
	if !isFailed(results) {
		return nil
	}
	fes := c.formatResultToOutput(logLevel, results)
	if len(fes) == 0 {
		return nil
	}
	if err := c.outputJSON(fes); err != nil {
		return fmt.Errorf("lint failed: output errors as JSON: %w", err)
	}
	return errors.New("lint failed")
}

type FlatError struct {
	RuleName     string      `json:"rule,omitempty"`
	Level        string      `json:"level,omitempty"`
	Message      string      `json:"message,omitempty"`
	LintFilePath string      `json:"lint_file,omitempty"`
	DataFilePath string      `json:"data_file,omitempty"`
	Location     interface{} `json:"location,omitempty"`
}

func (c *Controller) formatResultToOutput(logLevel ErrorLevel, results map[string]*FileResult) []*FlatError {
	list := make([]*FlatError, 0, len(results))
	for dataFilePath, fileResult := range results {
		list = append(list, fileResult.flattenError(logLevel, dataFilePath)...)
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
	return nil
}
