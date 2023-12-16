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
		Location    interface{} `json:"location,omitempty"`
		Metadata    interface{} `json:"metadata,omitempty"`
		Failed      bool        `json:"failed,omitempty"`
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

func (r *FileResult) flattenError(logLevel ErrorLevel, dataFilePath string) []*FlatError {
	if r.Error != "" {
		return []*FlatError{
			{
				DataFilePath: dataFilePath,
				Message:      r.Error,
			},
		}
	}
	list := make([]*FlatError, 0, len(r.Results))
	for lintFilePath, result := range r.Results {
		list = append(list, result.flattenError(logLevel, dataFilePath, lintFilePath)...)
	}
	if len(list) == 0 {
		return nil
	}
	return list
}

func (r *Result) flattenError(logLevel ErrorLevel, dataFilePath, lintFilePath string) []*FlatError {
	if r.Error != "" {
		return []*FlatError{
			{
				DataFilePath: dataFilePath,
				LintFilePath: lintFilePath,
				Message:      r.Error,
			},
		}
	}
	arr := make([]*FlatError, 0, len(r.RawResult))
	for _, result := range r.RawResult {
		arr = append(arr, result.flattenError(logLevel, dataFilePath, lintFilePath)...)
	}
	return arr
}

func (r *JsonnetResult) flattenError(logLevel ErrorLevel, dataFilePath, lintFilePath string) []*FlatError {
	if !r.Failed {
		return nil
	}
	fe := &FlatError{
		DataFilePath: dataFilePath,
		LintFilePath: lintFilePath,
		RuleName:     r.Name,
		Message:      r.Message,
		Location:     r.Location,
		Level:        r.Level,
	}
	if r.Level == "" {
		return []*FlatError{fe}
	}
	ll, err := newErrorLevel(r.Level)
	if err != nil {
		return []*FlatError{
			fe,
			{
				DataFilePath: dataFilePath,
				LintFilePath: lintFilePath,
				RuleName:     r.Name,
				Message:      err.Error(),
				Location:     r.Location,
			},
		}
	}
	if ll < logLevel {
		return nil
	}
	return []*FlatError{fe}
}

func (r *FileResult) isFailed() bool {
	if r.Error != "" {
		return true
	}
	for _, r := range r.Results {
		if r.isFailed() {
			return true
		}
	}
	return false
}

func (r *Result) isFailed() bool {
	if r.Error != "" {
		return true
	}
	for _, result := range r.RawResult {
		if result.isFailed() {
			return true
		}
	}
	return false
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
