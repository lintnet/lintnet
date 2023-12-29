package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/lintnet/lintnet/pkg/errlevel"
	"github.com/lintnet/lintnet/pkg/module"
)

type RawConfig struct {
	ErrorLevel string       `json:"error_level" yaml:"error_level"`
	Targets    []*RawTarget `json:"targets"`
	Outputs    []*RawOutput `json:"outputs,omitempty"`
}

func (rc *RawConfig) Parse() (*Config, error) {
	cfg := &Config{
		ErrorLevel: errlevel.Error,
		Targets:    make([]*Target, len(rc.Targets)),
	}
	if rc.ErrorLevel != "" {
		level, err := errlevel.New(rc.ErrorLevel)
		if err != nil {
			return nil, fmt.Errorf("parse the error level: %w", err)
		}
		cfg.ErrorLevel = level
	}
	moduleArchives := map[string]*module.Archive{}
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
	cfg.ModuleArchives = moduleArchives
	return cfg, nil
}

type Config struct {
	ErrorLevel     errlevel.Level
	Targets        []*Target
	Outputs        []*Output
	ModuleArchives map[string]*module.Archive
}

type RawOutput struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Renderer  string `json:"renderer"`
	SlashPath string `json:"path"`
	Template  string `json:"template"`
}

type Output struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Renderer string `json:"renderer"`
	Path     string `json:"path"`
	Template string `json:"template"`
}

type Target struct {
	LintFiles      []*module.Glob
	Modules        []*module.Glob
	ModuleArchives map[string]*module.Archive
	DataFiles      []string
}

type RawTarget struct {
	LintGlobs []*LintGlob  `json:"lint_files" yaml:"lint_files"`
	Modules   []*RawModule `json:"modules"`
	DataFiles []string     `json:"data_files" yaml:"data_files"`
}

type LintGlob struct {
	Glob  string         `json:"path"`
	Param map[string]any `json:"param"`
}

func (lg *LintGlob) UnmarshalJSON(b []byte) error {
	rm := &RawModule{}
	if err := json.Unmarshal(b, rm); err != nil {
		return err //nolint:wrapcheck
	}
	lg.Param = rm.Param
	lg.Glob = rm.Glob
	return nil
}

func (lg *LintGlob) ToModule() *module.Glob {
	p := strings.TrimPrefix(lg.Glob, "!")
	return &module.Glob{
		ID:        p,
		SlashPath: p,
		Param:     lg.Param,
		Excluded:  p != lg.Glob,
	}
}

func (rt *RawTarget) Parse() (*Target, error) {
	lintFiles := make([]*module.Glob, len(rt.LintGlobs))
	for i, lintGlob := range rt.LintGlobs {
		lintFiles[i] = lintGlob.ToModule()
	}
	target := &Target{
		LintFiles: lintFiles,
		Modules:   make([]*module.Glob, len(rt.Modules)),
		DataFiles: rt.DataFiles,
	}
	archives := make(map[string]*module.Archive, len(rt.Modules))
	for i, m := range rt.Modules {
		a, err := m.Parse()
		if err != nil {
			return nil, err
		}
		target.Modules[i] = a
		archives[a.Archive.ID] = a.Archive
	}
	target.ModuleArchives = archives
	return target, nil
}

type RawModule struct {
	Glob  string         `json:"path"`
	Param map[string]any `json:"param"`
}

func (rm *RawModule) Parse() (*module.Glob, error) {
	m, err := module.ParseModuleLine(rm.Glob)
	if err != nil {
		return nil, fmt.Errorf("parse a module path: %w", err)
	}
	m.Param = rm.Param
	return m, nil
}

type LintFile struct {
	ID    string
	Path  string         `json:"path"`
	Param map[string]any `json:"param"`
}

func (rm *RawModule) UnmarshalJSON(b []byte) error {
	var a any
	if err := json.Unmarshal(b, &a); err != nil {
		return fmt.Errorf("unmarshal as JSON: %w", err)
	}
	switch c := a.(type) {
	case string:
		rm.Glob = c
		return nil
	case map[string]any:
		p, ok := c["path"]
		if !ok {
			return errors.New("path is required")
		}
		a, ok := p.(string)
		if !ok {
			return errors.New("path must be a string")
		}
		rm.Glob = a

		param, ok := c["param"]
		if ok {
			a, ok := param.(map[string]any)
			if !ok {
				return errors.New("param must be a map[string]any")
			}
			rm.Param = a
		}
		return nil
	}
	return errors.New("module must be either string or map[string]any")
}
