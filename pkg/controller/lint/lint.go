package lint

import (
	"context"
	"errors"
	"fmt"

	"github.com/lintnet/lintnet/pkg/config"
	"github.com/lintnet/lintnet/pkg/errlevel"
	"github.com/lintnet/lintnet/pkg/log"
	"github.com/lintnet/lintnet/pkg/module"
	"github.com/sirupsen/logrus"
)

type ParamLint struct {
	ErrorLevel     string
	RootDir        string
	DataRootDir    string
	ConfigFilePath string
	TargetID       string
	FilePaths      []string
	Output         string
	OutputSuccess  bool
	PWD            string
}

func (c *Controller) Lint(ctx context.Context, logE *logrus.Entry, param *ParamLint) error { //nolint:cyclop,funlen
	rawCfg := &config.RawConfig{}
	if err := c.findAndReadConfig(param.ConfigFilePath, rawCfg); err != nil {
		return err
	}

	if param.TargetID != "" {
		target, err := c.getTarget(rawCfg.Targets, param.TargetID)
		if err != nil {
			return err
		}
		rawCfg.Targets = []*config.RawTarget{target}
	}

	cfg, err := rawCfg.Parse()
	if err != nil {
		return fmt.Errorf("parse a configuration file: %w", err)
	}

	outputter, err := c.getOutputter(cfg.Outputs, param.Output, param.RootDir)
	if err != nil {
		return err
	}

	errLevel, err := c.getErrorLevel(cfg, param)
	if err != nil {
		return err
	}

	if err := c.installModules(ctx, logE, &module.ParamInstall{
		BaseDir: param.RootDir,
	}, cfg.ModuleArchives); err != nil {
		return err
	}

	targets, err := c.findFiles(logE, cfg, param.RootDir)
	if err != nil {
		return err
	}

	if len(param.FilePaths) > 0 {
		logE.Debug("filtering targets by given files")
		if param.TargetID != "" {
			targets[0].DataFiles = param.FilePaths
		} else {
			targets = filterTargets(targets, param.FilePaths)
		}
	}

	if err := c.filterTargetsByDataRootDir(logE, param, targets); err != nil {
		return err
	}

	results, err := c.getResults(targets)
	if err != nil {
		return err
	}
	logE.WithFields(logrus.Fields{
		"config":  log.JSON(cfg),
		"results": log.JSON(results),
		"targets": log.JSON(targets),
	}).Debug("linted")

	return c.Output(logE, errLevel, results, []Outputter{outputter}, param.OutputSuccess)
}

func (c *Controller) getTarget(targets []*config.RawTarget, targetID string) (*config.RawTarget, error) {
	for _, target := range targets {
		if target.ID == targetID {
			return target, nil
		}
	}
	return nil, errors.New("target isn't found")
}

func (c *Controller) getResults(targets []*Target) ([]*Result, error) {
	results := make([]*Result, 0, len(targets))
	for _, target := range targets {
		rs, err := c.lintTarget(target)
		if err != nil {
			return nil, err
		}
		for _, r := range rs {
			r.TargetID = target.ID
		}
		results = append(results, rs...)
	}
	return results, nil
}

type DataSet struct {
	File  string
	Files []string
}

func (c *Controller) lintTarget(target *Target) ([]*Result, error) {
	lintFiles, err := c.parseLintFiles(target.LintFiles)
	if err != nil {
		return nil, err
	}
	if target.Combine {
		rs, err := c.lint(&DataSet{
			Files: target.DataFiles,
		}, lintFiles)
		if err != nil {
			return nil, err
		}
		for _, r := range rs {
			r.DataFiles = target.DataFiles
		}
		return rs, nil
	}
	results := make([]*Result, 0, len(target.DataFiles))
	for _, dataFile := range target.DataFiles {
		rs, err := c.lint(&DataSet{
			File: dataFile,
		}, lintFiles)
		if err != nil {
			results = append(results, &Result{
				DataFile: dataFile,
				Error:    err.Error(),
			})
			continue
		}
		for _, r := range rs {
			r.DataFile = dataFile
		}
		results = append(results, rs...)
	}
	return results, nil
}

func (c *Controller) getErrorLevel(cfg *config.Config, param *ParamLint) (errlevel.Level, error) {
	if param.ErrorLevel == "" {
		return cfg.ErrorLevel, nil
	}
	ll, err := errlevel.New(param.ErrorLevel)
	if err != nil {
		return ll, err //nolint:wrapcheck
	}
	return ll, nil
}

func (c *Controller) getTLA(dataSet *DataSet) (*TopLevelArgment, error) {
	if dataSet.File != "" {
		return c.parseDataFile(dataSet.File)
	}
	if len(dataSet.Files) > 0 {
		combinedData := make([]*Data, len(dataSet.Files))
		for i, dataFile := range dataSet.Files {
			data, err := c.parseDataFile(dataFile)
			if err != nil {
				return nil, err
			}
			combinedData[i] = data.Data
		}
		return &TopLevelArgment{
			CombinedData: combinedData,
		}, nil
	}
	return &TopLevelArgment{}, nil
}

func (c *Controller) lint(dataSet *DataSet, lintFiles []*Node) ([]*Result, error) {
	tla, err := c.getTLA(dataSet)
	if err != nil {
		return nil, err
	}
	return c.evaluate(tla, lintFiles), nil
}
