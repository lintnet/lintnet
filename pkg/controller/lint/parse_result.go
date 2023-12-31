package lint

import (
	"encoding/json"
	"fmt"
)

type (
	// return of vm.Evaluate()
	JsonnetEvaluateResult struct {
		// Key    string
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

	Result struct {
		TargetID  string           `json:"target_id,omitempty"`
		LintFile  string           `json:"lint_file,omitempty"`
		DataFile  string           `json:"data_file,omitempty"`
		DataFiles []string         `json:"data_files,omitempty"`
		RawResult []*JsonnetResult `json:"-"`
		RawOutput string           `json:"-"`
		Interface any              `json:"result,omitempty"`
		Error     string           `json:"error,omitempty"`
	}
)

func (result *Result) FlatErrors() []*FlatError {
	fes := make([]*FlatError, 0, len(result.RawResult))
	for _, r := range result.RawResult {
		if r.Excluded {
			continue
		}
		fes = append(fes, &FlatError{
			Rule:     r.Name,
			Level:    r.Level,
			Message:  r.Message,
			LintFile: result.LintFile,
			DataFile: result.DataFile,
			// DataFilePaths: result.DataFiles,
			TargetID: result.TargetID,
			Location: r.Location,
			Custom:   r.Custom,
		})
	}
	return fes
}

func (c *Controller) parseResult(result []byte) ([]*JsonnetResult, any, error) {
	var rs any
	if err := json.Unmarshal(result, &rs); err != nil {
		return nil, nil, fmt.Errorf("unmarshal the result as JSON: %w", err)
	}

	out := []*JsonnetResult{}
	if err := json.Unmarshal(result, &out); err != nil {
		return nil, rs, fmt.Errorf("unmarshal the result as JSON: %w", err)
	}
	return out, rs, nil
}
