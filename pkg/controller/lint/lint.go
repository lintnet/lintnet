package lint

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/lintnet/lintnet/pkg/config"
	"github.com/lintnet/lintnet/pkg/errlevel"
	"github.com/lintnet/lintnet/pkg/jsonnet"
	"github.com/lintnet/lintnet/pkg/module"
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

	modulesList, modMap, err := module.ListModules(cfg)
	if err != nil {
		return fmt.Errorf("list modules: %w", err)
	}
	if err := c.installModules(ctx, logE, &module.ParamInstall{
		BaseDir: param.RootDir,
	}, modMap); err != nil {
		return err
	}

	targets, err := c.findFiles(logE, cfg, modulesList, param.RuleBaseDir, param.FilePaths, param.RootDir)
	if err != nil {
		return err
	}

	errLevel, err := c.getErrorLevel(param.ErrorLevel)
	if err != nil {
		return err
	}

	results, err := c.getResults(targets)
	if err != nil {
		return err
	}

	return c.Output(logE, cfg, errLevel, results, param.Outputs, param.OutputSuccess)
}

func (c *Controller) getResults(targets []*Target) (map[string]*FileResult, error) {
	results := make(map[string]*FileResult, len(targets))
	for _, target := range targets {
		if err := c.lintTarget(target, results); err != nil {
			return nil, err
		}
	}
	return results, nil
}

func (c *Controller) lintTarget(target *Target, results map[string]*FileResult) error {
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
	return nil
}

func (c *Controller) getErrorLevel(errorLevel string) (errlevel.Level, error) {
	if errorLevel == "" {
		return errlevel.Error, nil
	}
	ll, err := errlevel.New(errorLevel)
	if err != nil {
		return ll, err //nolint:wrapcheck
	}
	return ll, nil
}

func (c *Controller) lint(dataFile string, jsonnetAsts map[string]jsonnet.Node) (map[string]*Result, error) {
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
