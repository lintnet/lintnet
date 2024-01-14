package lint

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/lintnet/lintnet/pkg/config"
	"github.com/lintnet/lintnet/pkg/errlevel"
	"github.com/lintnet/lintnet/pkg/filefilter"
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

func (p *ParamLint) FilterParam() *filefilter.Param {
	return &filefilter.Param{
		DataRootDir: p.DataRootDir,
		TargetID:    p.TargetID,
		FilePaths:   p.FilePaths,
		PWD:         p.PWD,
	}
}

func (c *Controller) Lint(ctx context.Context, logE *logrus.Entry, param *ParamLint) error { //nolint:cyclop,funlen
	rawCfg := &config.RawConfig{}
	if err := c.configReader.Read(param.ConfigFilePath, rawCfg); err != nil {
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

	if err := c.moduleInstaller.Installs(ctx, logE, &module.ParamInstall{
		BaseDir: param.RootDir,
	}, cfg.ModuleArchives); err != nil {
		return fmt.Errorf("install modules: %w", err)
	}

	targets, err := c.fileFinder.Find(logE, cfg, param.RootDir, cfgDir)
	if err != nil {
		return fmt.Errorf("find files: %w", err)
	}

	filterParam := param.FilterParam()

	if len(param.FilePaths) > 0 {
		logE.Debug("filtering targets by given files")
		targets = filefilter.FilterTargetsByFilePaths(filterParam, targets)
	}

	if err := filefilter.FilterTargetsByDataRootDir(logE, filterParam, targets); err != nil {
		return fmt.Errorf("filter targets by data root directory: %w", err)
	}

	results, err := c.linter.Lint(targets)
	if err != nil {
		return fmt.Errorf("lint targets: %w", err)
	}
	logE.WithFields(logrus.Fields{
		"config":  log.JSON(cfg),
		"results": log.JSON(results),
		"targets": log.JSON(targets),
	}).Debug("linted")

	return c.Output(logE, errLevel, shownErrLevel, results, []Outputter{outputter}, param.OutputSuccess)
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
