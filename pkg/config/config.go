package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/lintnet/lintnet/pkg/errlevel"
)

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

func (rc *RawConfig) Parse() (*Config, error) { //nolint:cyclop
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
		cfg.ShownErrorLevel = cfg.ErrorLevel
	}

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
	for _, output := range rc.Outputs {
		if strings.HasPrefix(output.Template, "github_archive/github.com/") {
			m, err := ParseImport(output.Template)
			if err != nil {
				return nil, fmt.Errorf("parse a module path: %w", err)
			}
			output.TemplateModule = m
			moduleArchives[m.Archive.String()] = m.Archive
		}
		if strings.HasPrefix(output.Transform, "github_archive/github.com/") {
			m, err := ParseImport(output.Transform)
			if err != nil {
				return nil, fmt.Errorf("parse a module path: %w", err)
			}
			output.TransformModule = m
			moduleArchives[m.Archive.String()] = m.Archive
		}
	}
	cfg.ModuleArchives = moduleArchives
	return cfg, nil
}
