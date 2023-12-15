package lint

import (
	"encoding/json"
)

type (
	FileResult struct {
		Results map[string]*Result `json:"results,omitempty"`
		Error   string             `json:"error,omitempty"`
	}
	Result struct {
		Output    *Output     `json:"-"`
		RawOutput string      `json:"-"`
		RawResult interface{} `json:"result,omitempty"`
		Error     string      `json:"error,omitempty"`
	}
	Output struct {
		Name        string  `json:"name,omitempty"`
		Description string  `json:"description,omitempty"`
		Rules       []*Rule `json:"rules,omitempty"`
	}
	Rule struct {
		Name        string   `json:"name,omitempty"`
		Description string   `json:"description,omitempty"`
		Errors      []*Error `json:"errors,omitempty"`
	}
	Error struct {
		Message string `json:"message,omitempty"`
	}
)

func (c *Controller) parseResult(result *Result) {
	if result.Error != "" {
		return
	}
	rb := []byte(result.RawOutput)

	var rs interface{}
	if err := json.Unmarshal(rb, &rs); err != nil {
		result.Error = err.Error()
		return
	}
	result.RawResult = rs

	out := &Output{}
	if err := json.Unmarshal(rb, out); err != nil {
		result.Error = err.Error()
		return
	}
	result.Output = out
}
