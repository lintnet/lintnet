package lint

import (
	"context"
	"errors"
	"os"

	"github.com/google/go-jsonnet/ast"
	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/lintnet/pkg/config"
)

type ParamLint struct {
	RuleBaseDir    string
	ErrorLevel     string
	ConfigFilePath string
	FilePaths      []string
}

func (c *Controller) Lint(_ context.Context, _ *logrus.Entry, param *ParamLint) error {
	cfg := &config.Config{}
	if err := c.findAndReadConfig(param.ConfigFilePath, cfg); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}

	// TODO download modules
	// TODO search lint files
	// TODO search data files

	filePaths, err := c.findJsonnet(param.RuleBaseDir)
	if err != nil {
		return err
	}
	jsonnetAsts, err := c.readJsonnets(filePaths)
	if err != nil {
		return err
	}

	errLevel := infoLevel
	if param.ErrorLevel != "" {
		ll, err := newErrorLevel(param.ErrorLevel)
		if err != nil {
			return err
		}
		errLevel = ll
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
	return c.Output(errLevel, results)
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
