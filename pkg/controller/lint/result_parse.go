package lint

import (
	"encoding/json"
	"fmt"
)

type (
	// Process Jsonnet

	// return of vm.Evaluate()
	JsonnetEvaluateResult struct {
		Result string
		Error  string
	}
	// unmarshal Jsonnet as JSON
	JsonnetResult struct {
		Name             string           `json:"name,omitempty"`
		ID               string           `json:"id,omitempty"`
		ShortDescription string           `json:"short_description,omitempty"`
		Description      string           `json:"description,omitempty"`
		Error            string           `json:"error,omitempty"`
		Level            string           `json:"level,omitempty"`
		Errors           []*Error         `json:"errors,omitempty"`
		Locations        []*Location      `json:"locations,omitempty"`
		SubRules         []*JsonnetResult `json:"sub_rules,omitempty"`
		Failed           bool             `json:"failed,omitempty"`
	}

	Location struct {
		S   string
		Raw interface{}
	}

	Error struct {
		Error    string    `json:"error,omitempty"`
		Level    string    `json:"level,omitempty"`
		FilePath string    `json:"file_path,omitempty"`
		Location *Location `json:"location,omitempty"`
	}

	// Aggregate results

	FileResult struct {
		// lint file -> result
		Results map[string]*Result `json:"results,omitempty"`
		Error   string             `json:"error,omitempty"`
	}

	Result struct {
		RawResult *JsonnetResult `json:"-"`
		RawOutput string         `json:"-"`
		Interface interface{}    `json:"result,omitempty"`
		Error     string         `json:"error,omitempty"`
	}

	// Format result to output
)

func (l *Location) MarshalJSON() ([]byte, error) {
	if l.S != "" {
		return []byte(l.S), nil
	}
	return json.Marshal(l.Raw) //nolint:wrapcheck
}

func (l *Location) UnmarshalJSON(b []byte) error {
	var a interface{}
	if err := json.Unmarshal(b, &a); err != nil {
		return err //nolint:wrapcheck
	}
	l.Raw = a
	if v, ok := a.(string); ok {
		l.S = v
		return nil
	}
	return nil
}

func (c *Controller) parseResult(result *JsonnetEvaluateResult) *Result {
	if result.Error != "" {
		return &Result{
			RawOutput: result.Result,
			Error:     result.Error,
		}
	}
	rb := []byte(result.Result)

	var rs interface{}
	if err := json.Unmarshal(rb, &rs); err != nil {
		return &Result{
			RawOutput: result.Result,
			Error:     result.Error,
		}
	}

	out := &JsonnetResult{}
	if err := json.Unmarshal(rb, out); err != nil {
		return &Result{
			RawOutput: result.Result,
			Interface: rs,
			Error:     fmt.Errorf("unmarshal the result as JSON: %w", err).Error(),
		}
	}
	return &Result{
		RawOutput: result.Result,
		RawResult: out,
		Interface: rs,
	}
}
