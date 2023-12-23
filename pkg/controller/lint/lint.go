package lint

import (
	"context"
	"errors"
	"os"

	"github.com/google/go-jsonnet/ast"
	"github.com/lintnet/lintnet/pkg/config"
	"github.com/sirupsen/logrus"
)

type ParamLint struct {
	RuleBaseDir    string
	ErrorLevel     string
	RootDir        string
	ConfigFilePath string
	FilePaths      []string
	Outputs        []string
	OutputSuccess  bool
}

func (c *Controller) Lint(ctx context.Context, logE *logrus.Entry, param *ParamLint) error {
	cfg := &config.Config{}
	if err := c.findAndReadConfig(param.ConfigFilePath, cfg); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}

	modulesList, modMap, err := c.listModules(logE, cfg)
	if err != nil {
		return err
	}
	if err := c.downloadModules(ctx, logE, &ParamDownloadModule{
		BaseDir: param.RootDir,
	}, modMap); err != nil {
		return err
	}

	targets, err := c.findFiles(logE, cfg, modulesList, param.RuleBaseDir, param.FilePaths, param.RootDir)
	if err != nil {
		return err
	}

	errLevel := errorLevel
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
			rs, err := c.lint(dataFile, jsonnetAsts)
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

	return c.Output(logE, cfg, errLevel, results, param.Outputs, param.OutputSuccess)
}

func (c *Controller) lint(dataFile string, jsonnetAsts map[string]ast.Node) (map[string]*Result, error) {
	data, err := c.parse(dataFile)
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
