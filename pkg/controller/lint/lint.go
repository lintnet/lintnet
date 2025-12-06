package lint

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"

	"github.com/lintnet/lintnet/pkg/config"
	"github.com/lintnet/lintnet/pkg/errlevel"
	"github.com/lintnet/lintnet/pkg/filefilter"
	"github.com/lintnet/lintnet/pkg/log"
	"github.com/lintnet/lintnet/pkg/module"
	"github.com/lintnet/lintnet/pkg/output"
)

type ParamLint struct {
	ErrorLevel      string   `json:"error_level,omitempty"`
	ShownErrorLevel string   `json:"shown_error_level,omitempty"`
	RootDir         string   `json:"root_dir,omitempty"`
	DataRootDir     string   `json:"data_root_dir,omitempty"`
	ConfigFilePath  string   `json:"config_file_path,omitempty"`
	TargetID        string   `json:"target_id,omitempty"`
	FilePaths       []string `json:"file_paths,omitempty"`
	Output          string   `json:"output,omitempty"`
	OutputSuccess   bool     `json:"output_success,omitempty"`
	PWD             string   `json:"pwd,omitempty"`
}

func (p *ParamLint) FilterParam() *filefilter.Param {
	return &filefilter.Param{
		DataRootDir: p.DataRootDir,
		TargetID:    p.TargetID,
		FilePaths:   p.FilePaths,
		PWD:         p.PWD,
	}
}

func (p *ParamLint) OutputterParam() *output.ParamGet {
	return &output.ParamGet{
		RootDir: p.RootDir,
		Output:  p.Output,
	}
}

// Lint lints files.
func (c *Controller) Lint(ctx context.Context, logger *slog.Logger, param *ParamLint) error { //nolint:cyclop,funlen
	logger.Debug("parameter", "param", log.JSON(param))
	// Find and read a configuration file.
	rawCfg := &config.RawConfig{}
	if err := c.configReader.Read(param.ConfigFilePath, rawCfg); err != nil {
		return fmt.Errorf("read a configuration file: %w", err)
	}

	logger.Debug("read config", "config", log.JSON(rawCfg))

	if param.TargetID != "" {
		// If a target id is specified, gets a target from the configuration file by the target id.
		target, err := rawCfg.GetTarget(param.TargetID)
		if err != nil {
			return fmt.Errorf("get a target from configuration file by target id: %w", err)
		}
		rawCfg.Targets = []*config.RawTarget{target}
	}

	// Parse the configuration file.
	cfg, err := rawCfg.Parse()
	if err != nil {
		return fmt.Errorf("parse a configuration file: %w", err)
	}

	logger.Debug("parse config", "config", log.JSON(cfg), "raw_config", log.JSON(rawCfg))

	// Get a directory of the configuration file.
	cfgDir := filepath.Dir(rawCfg.FilePath)
	if !filepath.IsAbs(cfgDir) {
		cfgDir = filepath.Join(param.PWD, cfgDir)
	}
	cfgDir = filepath.Clean(cfgDir)

	// Get an outputter.
	outputter, err := c.outputGetter.Get(cfg.Outputs, param.OutputterParam(), cfgDir)
	if err != nil {
		return fmt.Errorf("get an outputter: %w", err)
	}

	errLevel, err := getErrorLevel(param.ErrorLevel, cfg.ErrorLevel)
	if err != nil {
		return err
	}

	shownErrLevel, err := getErrorLevel(param.ShownErrorLevel, cfg.ShownErrorLevel)
	if err != nil {
		return err
	}

	modRootDir := filepath.Join(param.RootDir, "modules")

	// Install modules.
	if err := c.moduleInstaller.Installs(ctx, logger, &module.ParamInstall{
		BaseDir: modRootDir,
	}, cfg.ModuleArchives); err != nil {
		return fmt.Errorf("install modules: %w", err)
	}

	// Find targets, which are pairs of lint files and data files.
	targets, err := c.fileFinder.Find(logger, cfg, modRootDir, cfgDir)
	if err != nil {
		return fmt.Errorf("find files: %w", err)
	}

	logger.Debug("found files", "targets", log.JSON(targets))

	filterParam := param.FilterParam()

	if len(param.FilePaths) > 0 {
		// If files are specified, filters targets by the files.
		targets = filefilter.FilterTargetsByFilePaths(filterParam, targets)
		logger.Debug("filtered targets by given files", "filter_param", log.JSON(filterParam), "targets", log.JSON(targets))
	}

	// Lint targets.
	results, err := c.linter.Lint(targets)
	if err != nil {
		return fmt.Errorf("lint targets: %w", err)
	}
	logger.Debug("linted", "config", log.JSON(cfg), "results", log.JSON(results), "targets", log.JSON(targets))

	// Output results.
	return c.Output(logger, errLevel, shownErrLevel, results, []Outputter{outputter}, param.OutputSuccess)
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
