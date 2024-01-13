package lint

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/lintnet/lintnet/pkg/config"
	"github.com/lintnet/lintnet/pkg/errlevel"
	"github.com/lintnet/lintnet/pkg/log"
	"github.com/lintnet/lintnet/pkg/module"
	"github.com/sirupsen/logrus"
)

type ParamLint struct {
	ErrorLevel      string
	ShownErrorLevel string
	RootDir         string
	DataRootDir     string
	ConfigFilePath  string
	TargetID        string
	FilePaths       []string
	Output          string
	OutputSuccess   bool
	PWD             string
}

type DataSet struct {
	File  *Path
	Files []*Path
}

type Paths []*Path

func (ps Paths) Raw() []string {
	arr := make([]string, len(ps))
	for i, p := range ps {
		arr[i] = p.Raw
	}
	return arr
}

func (c *Controller) Lint(ctx context.Context, logE *logrus.Entry, param *ParamLint) error { //nolint:cyclop,funlen
	rawCfg := &config.RawConfig{}
	if err := c.findAndReadConfig(param.ConfigFilePath, rawCfg); err != nil {
		return err
	}

	if param.TargetID != "" {
		target, err := getTarget(rawCfg.Targets, param.TargetID)
		if err != nil {
			return err
		}
		rawCfg.Targets = []*config.RawTarget{target}
	}

	cfg, err := rawCfg.Parse()
	if err != nil {
		return fmt.Errorf("parse a configuration file: %w", err)
	}

	cfgDir := filepath.Dir(rawCfg.FilePath)
	if !filepath.IsAbs(cfgDir) {
		cfgDir = filepath.Join(param.PWD, cfgDir)
	}
	outputter, err := c.getOutputter(cfg.Outputs, param, cfgDir)
	if err != nil {
		return err
	}

	errLevel, err := getErrorLevel(param.ErrorLevel, cfg.ErrorLevel)
	if err != nil {
		return err
	}

	shownErrLevel, err := getErrorLevel(param.ShownErrorLevel, cfg.ShownErrorLevel)
	if err != nil {
		return err
	}

	if err := c.installModules(ctx, logE, &module.ParamInstall{
		BaseDir: param.RootDir,
	}, cfg.ModuleArchives); err != nil {
		return err
	}

	targets, err := c.findFiles(logE, cfg, param.RootDir, cfgDir)
	if err != nil {
		return err
	}

	if len(param.FilePaths) > 0 {
		logE.Debug("filtering targets by given files")
		if param.TargetID != "" {
			arr := make([]*Path, len(param.FilePaths))
			for i, filePath := range param.FilePaths {
				p := &Path{
					Abs: filePath,
					Raw: filePath,
				}
				if !filepath.IsAbs(filePath) {
					p.Abs = filepath.Join(param.PWD, filePath)
				}
				arr[i] = p
			}
			targets[0].DataFiles = arr
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

	return c.Output(logE, errLevel, shownErrLevel, results, []Outputter{outputter}, param.OutputSuccess)
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

func (c *Controller) lintCombineFiles(target *Target, combineFiles []*Node) ([]*Result, error) {
	rs, err := c.lint(&DataSet{
		Files: target.DataFiles,
	}, combineFiles)
	if err != nil {
		return nil, err
	}
	for _, r := range rs {
		arr := make([]string, len(target.DataFiles))
		for i, dataFile := range target.DataFiles {
			arr[i] = dataFile.Raw
		}
		r.DataFiles = arr
	}
	return rs, nil
}

func (c *Controller) lintNonCombineFiles(target *Target, nonCombineFiles []*Node) []*Result {
	results := make([]*Result, 0, len(target.DataFiles))
	for _, dataFile := range target.DataFiles {
		rs, err := c.lint(&DataSet{
			File: dataFile,
		}, nonCombineFiles)
		if err != nil {
			results = append(results, &Result{
				DataFile: dataFile.Raw,
				Error:    err.Error(),
			})
			continue
		}
		for _, r := range rs {
			r.DataFile = dataFile.Raw
		}
		results = append(results, rs...)
	}
	return results
}

func (c *Controller) lintTarget(target *Target) ([]*Result, error) {
	lintFiles, err := c.parseLintFiles(target.LintFiles)
	if err != nil {
		return nil, err
	}

	combineFiles := []*Node{}
	nonCombineFiles := []*Node{}
	for _, lintFile := range lintFiles {
		if lintFile.Combine {
			combineFiles = append(combineFiles, lintFile)
			continue
		}
		nonCombineFiles = append(nonCombineFiles, lintFile)
	}

	results := c.lintNonCombineFiles(target, nonCombineFiles)

	if len(combineFiles) > 0 {
		rs, err := c.lintCombineFiles(target, combineFiles)
		if err != nil {
			return nil, err
		}
		return append(results, rs...), nil
	}
	return results, nil
}

func (c *Controller) getTLA(dataSet *DataSet) (*TopLevelArgment, error) {
	if dataSet.File != nil {
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

func getErrorLevel(errLevel string, defaultErrorLevel errlevel.Level) (errlevel.Level, error) {
	if errLevel == "" {
		return defaultErrorLevel, nil
	}
	ll, err := errlevel.New(errLevel)
	if err != nil {
		return ll, err //nolint:wrapcheck
	}
	return ll, nil
}

func getTarget(targets []*config.RawTarget, targetID string) (*config.RawTarget, error) {
	for _, target := range targets {
		if target.ID == targetID {
			return target, nil
		}
	}
	return nil, errors.New("target isn't found")
}
