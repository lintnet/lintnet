package lint

import (
	"context"

	"github.com/google/go-jsonnet/ast"
)

type (
	ParamLint struct {
		RuleBaseDir string
		FilePaths   []string
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
