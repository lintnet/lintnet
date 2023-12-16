package lint

import (
	"context"

	"github.com/google/go-jsonnet/ast"
	"github.com/sirupsen/logrus"
)

type (
	ParamLint struct {
		RuleBaseDir string
		FilePaths   []string
	}
)

func (c *Controller) Lint(_ context.Context, _ *logrus.Entry, param *ParamLint) error {
	filePaths, err := c.findJsonnet(param.RuleBaseDir)
	if err != nil {
		return err
	}
	jsonnetAsts, err := c.readJsonnets(filePaths)
	if err != nil {
		return err
	}

	logLevel := infoLevel

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
	return c.Output(logLevel, results)
}

func (c *Controller) lint(filePath string, jsonnetAsts map[string]ast.Node) (map[string]*Result, error) {
	data, err := c.parse(filePath)
	if err != nil {
		return nil, err
	}

	results := c.evaluate(data, jsonnetAsts)
	rs := make(map[string]*Result, len(results))

	for k, result := range results {
		rs[k] = c.parseResult(result)
	}
	return rs, nil
}
