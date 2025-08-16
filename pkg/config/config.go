package config

import (
	"errors"
	"fmt"

	"github.com/lintnet/lintnet/pkg/errlevel"
)

// getIgnoredPatterns returns ignored patterns.
// If ignoredDirs is nil, it returns default ignored patterns.
// An ignored pattern is "**/<ignored dir>/**", which is a pattern of doublestar.
func getIgnoredPatterns(ignoredDirs []string) []string {
	if ignoredDirs == nil {
		ignoredDirs = []string{
			".git",
			"node_modules",
		}
	}
	ignoredPatterns := make([]string, len(ignoredDirs))
	for i, d := range ignoredDirs {
		ignoredPatterns[i] = fmt.Sprintf("**/%s/**", d)
	}
	return ignoredPatterns
}

type Config struct {
	ErrorLevel      errlevel.Level            `json:"error_level,omitempty"`
	ShownErrorLevel errlevel.Level            `json:"shown_error_level,omitempty"`
	Targets         []*Target                 `json:"targets,omitempty"`
	Outputs         Outputs                   `json:"outputs,omitempty"`
	ModuleArchives  map[string]*ModuleArchive `json:"module_archives,omitempty"`
	IgnoredPatterns []string                  `json:"ignore_patterns,omitempty"`
}

func (c *Config) setErrorLevel(errLevel string) error {
	if errLevel == "" {
		c.ErrorLevel = errlevel.Error
		return nil
	}
	level, err := errlevel.New(errLevel)
	if err != nil {
		return fmt.Errorf("parse the error level: %w", err)
	}
	c.ErrorLevel = level
	return nil
}

func (c *Config) setShownErrorLevel(errLevel string) error {
	if errLevel == "" {
		c.ShownErrorLevel = errlevel.Info
		return nil
	}
	level, err := errlevel.New(errLevel)
	if err != nil {
		return fmt.Errorf("parse the error level: %w", err)
	}
	c.ShownErrorLevel = level
	return nil
}

// setIgnoredPatterns sets ignored patterns.
func (c *Config) setIgnoredPatterns(dirs []string) {
	c.IgnoredPatterns = getIgnoredPatterns(dirs)
	if c.IgnoredPatterns == nil {
		c.IgnoredPatterns = []string{
			"node_modules",
			".git",
		}
	}
}

type RawConfig struct {
	FilePath        string       `json:"-"`
	ErrorLevel      string       `json:"error_level,omitempty"`
	ShownErrorLevel string       `json:"shown_error_level,omitempty"`
	IgnoredDirs     []string     `json:"ignored_dirs,omitempty"`
	Targets         []*RawTarget `json:"targets"`
	Outputs         Outputs      `json:"outputs,omitempty"`
}

func (rc *RawConfig) GetTarget(targetID string) (*RawTarget, error) {
	for _, target := range rc.Targets {
		if target.ID == targetID {
			return target, nil
		}
	}
	return nil, errors.New("target isn't found")
}

// Parse processes a raw configuration.
func (rc *RawConfig) Parse() (*Config, error) {
	cfg := &Config{
		Targets: make([]*Target, len(rc.Targets)),
		Outputs: rc.Outputs,
	}
	cfg.setIgnoredPatterns(rc.IgnoredDirs)

	if err := cfg.setErrorLevel(rc.ErrorLevel); err != nil {
		return nil, err
	}

	if err := cfg.setShownErrorLevel(rc.ShownErrorLevel); err != nil {
		return nil, err
	}

	if cfg.ShownErrorLevel > cfg.ErrorLevel {
		// ShownErrorLevel should be lower than or equal to ErrorLevel.
		// If ShownErrorLevel is higher than ErrorLevel, it sets ShownErrorLevel to ErrorLevel.
		cfg.ShownErrorLevel = cfg.ErrorLevel
	}

	// moduleArchives is a map of modules.
	// Extract modules from configuration to install them.
	moduleArchives := map[string]*ModuleArchive{}
	for i, rt := range rc.Targets {
		target, err := rt.Parse()
		if err != nil {
			return nil, err
		}
		cfg.Targets[i] = target
		for k, ma := range target.ModuleArchives {
			moduleArchives[k] = ma
		}
	}
	if err := rc.Outputs.Preprocess(moduleArchives); err != nil {
		return nil, err
	}
	cfg.ModuleArchives = moduleArchives
	return cfg, nil
}
