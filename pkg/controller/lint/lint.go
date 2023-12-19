package lint

import (
	"context"
	"errors"
	"os"

	"github.com/google/go-jsonnet/ast"
	"github.com/sirupsen/logrus"
	"github.com/lintnet/lintnet/pkg/config"
)

type ParamLint struct {
	RuleBaseDir    string
	ErrorLevel     string
	ConfigFilePath string
	FilePaths      []string
}

func (c *Controller) Lint(_ context.Context, logE *logrus.Entry, param *ParamLint) error {
	cfg := &config.Config{}
	if err := c.findAndReadConfig(param.ConfigFilePath, cfg); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}

	// TODO download modules
	// TODO search lint files
	// TODO search data files

	targets, err := c.findFiles(logE, cfg, param.RuleBaseDir, param.FilePaths)
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

	results := make(map[string]*FileResult, len(targets))
	for _, target := range targets {
		jsonnetAsts, err := c.readJsonnets(target.LintFiles)
		if err != nil {
			return err
		}
		for _, dataFile := range target.DataFiles {
			rs, err := c.lint(dataFile, jsonnetAsts, nil)
			if err != nil {
				results[dataFile] = &FileResult{
					Error: err.Error(),
				}
				continue
			}
			results[dataFile] = &FileResult{
				Results: rs,
			}
		}
	}

	return c.Output(logE, cfg, errLevel, results)
}

func (c *Controller) lint(dataFile string, jsonnetAsts map[string]ast.Node, libs map[string]string) (map[string]*Result, error) {
	data, err := c.parse(dataFile)
	if err != nil {
		return nil, err
	}

	results := c.evaluate(data, jsonnetAsts, libs)
	rs := make(map[string]*Result, len(results))

	for k, result := range results {
		rs[k] = c.parseResult(result)
	}
	return rs, nil
}
