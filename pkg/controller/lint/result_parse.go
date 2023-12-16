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
		Name        string      `json:"name,omitempty"`
		Description string      `json:"description,omitempty"`
		Message     string      `json:"message,omitempty"`
		Level       string      `json:"level,omitempty"`
		Location    *Location   `json:"location,omitempty"`
		Metadata    interface{} `json:"metadata,omitempty"`
		Failed      bool        `json:"failed,omitempty"`
	}

	Location struct {
		S   string
		Raw interface{}
	}

	// Aggregate results

	FileResult struct {
		// lint file -> result
		Results map[string]*Result `json:"results,omitempty"`
		Error   string             `json:"error,omitempty"`
	}

	Result struct {
		RawResult []*JsonnetResult `json:"-"`
		RawOutput string           `json:"-"`
		Interface interface{}      `json:"result,omitempty"`
		Error     string           `json:"error,omitempty"`
	}
)

func (r *FileResult) flattenError(p string) []*FlatError {
	if r.Error != "" {
		return []*FlatError{
			{
				LintFilePath: p,
				Error:        r.Error,
			},
		}
	}
	list := make([]*FlatError, 0, len(r.Results))
	for lintFilePath, result := range r.Results {
		list = append(list, result.flattenError(p, lintFilePath)...)
	}
	if len(list) == 0 {
		return nil
	}
	return list
}

func (r *Result) flattenError(dataFilePath, lintFilePath string) []*FlatError {
	if r.Error != "" {
		return []*FlatError{
			{
				DataFilePath: dataFilePath,
				LintFilePath: lintFilePath,
				Error:        r.Error,
			},
		}
	}
	arr := make([]*FlatError, 0, len(r.RawResult))
	for _, result := range r.RawResult {
		arr = append(arr, result.flattenError(dataFilePath, lintFilePath)...)
	}
	return arr
}

func (r *JsonnetResult) flattenError(dataFilePath, lintFilePath string) []*FlatError {
	if r.Failed {
		return []*FlatError{
			{
				DataFilePath: dataFilePath,
				LintFilePath: lintFilePath,
				RuleName:     r.Name,
				Error:        r.Message,
				Location:     r.Location,
			},
		}
	}
	return nil
}

func (r *FileResult) isFailed() bool {
	if r.Error != "" {
		return true
	}
	for _, r := range r.Results {
		if r.Error != "" {
			return true
		}
		for _, result := range r.RawResult {
			if result.isFailed() {
				return true
			}
		}
	}
	return false
}

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

	out := []*JsonnetResult{}
	if err := json.Unmarshal(rb, &out); err != nil {
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
