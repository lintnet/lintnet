package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/lintnet/lintnet/pkg/errlevel"
)

type RawConfig struct {
	ErrorLevel string       `json:"error_level"`
	Targets    []*RawTarget `json:"targets"`
	Outputs    []*Output    `json:"outputs,omitempty"`
}

func (rc *RawConfig) Parse() (*Config, error) {
	cfg := &Config{
		ErrorLevel: errlevel.Error,
		Targets:    make([]*Target, len(rc.Targets)),
		Outputs:    rc.Outputs,
	}
	if rc.ErrorLevel != "" {
		level, err := errlevel.New(rc.ErrorLevel)
		if err != nil {
			return nil, fmt.Errorf("parse the error level: %w", err)
		}
		cfg.ErrorLevel = level
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
	cfg.ModuleArchives = moduleArchives
	return cfg, nil
}

type Config struct {
	ErrorLevel     errlevel.Level
	Targets        []*Target
	Outputs        []*Output
	ModuleArchives map[string]*ModuleArchive
}

type Output struct {
	ID string `json:"id"`
	// text/template, html/template, jsonnet
	Renderer string `json:"renderer"`
	// path to a template file
	Template string `json:"template"`
	// parameter
	Config map[string]any `json:"config"`
}

type Target struct {
	ID             string
	LintFiles      []*ModuleGlob
	Modules        []*ModuleGlob
	ModuleArchives map[string]*ModuleArchive
	DataFiles      []string
	Combine        bool
}

type RawTarget struct {
	ID        string       `json:"id,omitempty"`
	LintGlobs []*LintGlob  `json:"lint_files"`
	Modules   []*RawModule `json:"modules"`
	DataFiles []string     `json:"data_files"`
	Combine   bool         `json:"combine,omitempty"`
}

type LintGlob struct {
	Glob   string         `json:"path"`
	Config map[string]any `json:"config"`
}

func (lg *LintGlob) UnmarshalJSON(b []byte) error {
	rm := &RawModule{}
	if err := json.Unmarshal(b, rm); err != nil {
		return err //nolint:wrapcheck
	}
	lg.Config = rm.Config
	lg.Glob = rm.Glob
	return nil
}

func (lg *LintGlob) ToModule() *ModuleGlob {
	p := strings.TrimPrefix(lg.Glob, "!")
	return &ModuleGlob{
		ID:        p,
		SlashPath: p,
		Config:    lg.Config,
		Excluded:  p != lg.Glob,
	}
}

func (rt *RawTarget) Parse() (*Target, error) {
	lintFiles := make([]*ModuleGlob, len(rt.LintGlobs))
	for i, lintGlob := range rt.LintGlobs {
		lintFiles[i] = lintGlob.ToModule()
	}
	target := &Target{
		ID:        rt.ID,
		Combine:   rt.Combine,
		LintFiles: lintFiles,
		Modules:   make([]*ModuleGlob, len(rt.Modules)),
		DataFiles: rt.DataFiles,
	}
	archives := make(map[string]*ModuleArchive, len(rt.Modules))
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
	Glob   string         `json:"path"`
	Config map[string]any `json:"config"`
}

func (rm *RawModule) Parse() (*ModuleGlob, error) {
	m, err := ParseModuleLine(rm.Glob)
	if err != nil {
		return nil, fmt.Errorf("parse a module path: %w", err)
	}
	m.Config = rm.Config
	return m, nil
}

type LintFile struct {
	ID     string
	Path   string         `json:"path"`
	Config map[string]any `json:"config"`
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
			rm.Config = a
		}
		return nil
	}
	return errors.New("module must be either string or map[string]any")
}
