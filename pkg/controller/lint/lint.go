package lint

import (
	"context"
	"encoding/json"
	"io"

	"github.com/google/go-jsonnet/ast"
)

type (
	ParamLint struct {
		RuleBaseDir string
		FilePaths   []string
	}

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

	NewDecoder func(io.Reader) decoder
	decoder    interface {
		Decode(dest interface{}) error
	}
)

func (c *Controller) Lint(_ context.Context, param *ParamLint) error {
	filePaths, err := c.findJsonnet(param.RuleBaseDir)
	if err != nil {
		return err
	}
	jsonnetAsts, err := c.readJsonnets(filePaths)
	if err != nil {
		return err
	}

	results := make(map[string]*FileResult, len(param.FilePaths))
	for _, filePath := range param.FilePaths {
		rs, err := c.lint(filePath, jsonnetAsts)
		if err != nil {
			results[filePath] = &FileResult{
				Error: err.Error(),
			}
			continue
		}
		results[filePath] = &FileResult{
			Results: rs,
		}
	}
	return c.Output(results)
}

func (c *Controller) lint(filePath string, jsonnetAsts map[string]ast.Node) (map[string]*Result, error) {
	input, fileType, err := c.parse(filePath)
	if err != nil {
		return nil, err
	}

	results := c.evaluate(input, fileType, filePath, jsonnetAsts)

	for _, result := range results {
		c.parseResult(result)
	}
	return results, nil
}

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

func (c *Controller) evaluate(input []byte, filePath, fileType string, jsonnetAsts map[string]ast.Node) map[string]*Result {
	vm := newVM(input, filePath, fileType)

	results := make(map[string]*Result, len(jsonnetAsts))
	for k, ja := range jsonnetAsts {
		result, err := vm.Evaluate(ja)
		if err != nil {
			results[k] = &Result{
				Error: err.Error(),
			}
			continue
		}
		results[k] = &Result{
			RawOutput: result,
		}
	}
	return results
}
