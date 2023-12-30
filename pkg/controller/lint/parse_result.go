package lint

import (
	"encoding/json"
	"fmt"
)

type (
	// Process Jsonnet

	// return of vm.Evaluate()
	JsonnetEvaluateResult struct {
		Key    string
		Result string
		Error  string
	}
	// unmarshal Jsonnet as JSON
	JsonnetResult struct {
		Name        string `json:"name,omitempty"`
		Description string `json:"description,omitempty"`
		Message     string `json:"message,omitempty"`
		Level       string `json:"level,omitempty"`
		Location    any    `json:"location,omitempty"`
		Custom      any    `json:"custom,omitempty"`
		Excluded    bool   `json:"excluded,omitempty"`
	}

	// Aggregate results

	FileResult struct {
		// lint file -> result
		Results []*Result `json:"results,omitempty"`
		Error   string    `json:"error,omitempty"`
	}

	Result struct {
		Key       string           `json:"-"`
		RawResult []*JsonnetResult `json:"-"`
		RawOutput string           `json:"-"`
		Interface any              `json:"result,omitempty"`
		Error     string           `json:"error,omitempty"`
	}
)

func (r *FileResult) FlattenError(dataFilePath string) []*FlatError {
	if r.Error != "" {
		return []*FlatError{
			{
				DataFilePath: dataFilePath,
				Message:      r.Error,
			},
		}
	}
	list := make([]*FlatError, 0, len(r.Results))
	for _, result := range r.Results {
		list = append(list, result.flattenError(dataFilePath, result.Key)...)
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
				Message:      r.Error,
				Custom: map[string]any{
					"result": r.Interface,
				},
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
	if r.Excluded {
		return nil
	}
	return []*FlatError{
		{
			DataFilePath: dataFilePath,
			LintFilePath: lintFilePath,
			RuleName:     r.Name,
			Message:      r.Message,
			Location:     r.Location,
			Level:        r.Level,
			Custom:       r.Custom,
		},
	}
}

func (c *Controller) parseResult(result *JsonnetEvaluateResult) *Result {
	// jsonnet VM returns the result as JSON string,
	// so parse the JSON string to structured data
	if result.Error != "" {
		return &Result{
			Key:       result.Key,
			RawOutput: result.Result,
			Error:     result.Error,
		}
	}
	rb := []byte(result.Result)

	var rs any
	if err := json.Unmarshal(rb, &rs); err != nil {
		return &Result{
			Key:       result.Key,
			RawOutput: result.Result,
			Error:     result.Error,
		}
	}

	out := []*JsonnetResult{}
	if err := json.Unmarshal(rb, &out); err != nil {
		return &Result{
			Key:       result.Key,
			RawOutput: result.Result,
			Interface: rs,
			Error:     fmt.Errorf("unmarshal the result as JSON: %w", err).Error(),
		}
	}
	return &Result{
		Key:       result.Key,
		RawOutput: result.Result,
		RawResult: out,
		Interface: rs,
	}
}
