package domain

import (
	"fmt"

	"github.com/google/go-jsonnet/ast"
	"github.com/lintnet/lintnet/pkg/errlevel"
)

type (
	// unmarshal Jsonnet as JSON
	JsonnetResult struct {
		Name        string `json:"name,omitempty"`
		Description string `json:"description,omitempty"`
		Links       Links  `json:"links,omitempty"`
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

func (result *Result) FlatErrors() []*Error {
	if result.Error != "" {
		return []*Error{
			{
				LintFile: result.LintFile,
				DataFile: result.DataFile,
				// DataFilePaths: result.DataFiles,
				TargetID: result.TargetID,
				Message:  result.Error,
			},
		}
	}
	fes := make([]*Error, 0, len(result.RawResult))
	for _, r := range result.RawResult {
		if r.Excluded {
			continue
		}
		fes = append(fes, &Error{
			Name:        r.Name,
			Level:       r.Level,
			Message:     r.Message,
			Description: r.Description,
			LintFile:    result.LintFile,
			DataFile:    result.DataFile,
			Links:       r.Links,
			// DataFilePaths: result.DataFiles,
			TargetID: result.TargetID,
			Location: r.Location,
			Custom:   r.Custom,
		})
	}
	return fes
}

type Error struct {
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	Links       []*Link `json:"links,omitempty"`
	Level       string  `json:"level,omitempty"`
	Message     string  `json:"message,omitempty"`
	LintFile    string  `json:"lint_file,omitempty"`
	DataFile    string  `json:"data_file,omitempty"`
	// DataFilePaths []string `json:"data_files,omitempty"`
	TargetID string `json:"target_id,omitempty"`
	Location any    `json:"location,omitempty"`
	Custom   any    `json:"custom,omitempty"`
}

func (e *Error) Failed(errLevel errlevel.Level) (bool, error) {
	level := errlevel.Error
	if e.Level != "" {
		feErrLevel, err := errlevel.New(e.Level)
		if err != nil {
			return false, fmt.Errorf("verify the error level of a result: %w", err)
		}
		level = feErrLevel
	}
	return level >= errLevel, nil
}

type Node struct {
	Node    ast.Node
	Config  map[string]any
	Key     string
	Combine bool
}
