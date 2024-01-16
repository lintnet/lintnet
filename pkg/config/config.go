package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/lintnet/lintnet/pkg/domain"
	"github.com/lintnet/lintnet/pkg/errlevel"
)

type RawConfig struct {
	FilePath        string       `json:"-"`
	ErrorLevel      string       `json:"error_level,omitempty"`
	ShownErrorLevel string       `json:"shown_error_level,omitempty"`
	IgnoredDirs     []string     `json:"ignored_dirs,omitempty"`
	Targets         []*RawTarget `json:"targets"`
	Outputs         []*RawOutput `json:"outputs,omitempty"`
}

func (rc *RawConfig) GetTarget(targetID string) (*RawTarget, error) {
	for _, target := range rc.Targets {
		if target.ID == targetID {
			return target, nil
		}
	}
	return nil, errors.New("target isn't found")
}

type Config struct {
	ErrorLevel      errlevel.Level            `json:"error_level,omitempty"`
	ShownErrorLevel errlevel.Level            `json:"shown_error_level,omitempty"`
	Targets         []*Target                 `json:"targets,omitempty"`
	Outputs         []*Output                 `json:"outputs,omitempty"`
	ModuleArchives  map[string]*ModuleArchive `json:"module_archives,omitempty"`
	IgnoredPatterns []string                  `json:"ignore_patterns,omitempty"`
}

type RawOutput struct {
	ID string `json:"id"`
	// text/template, html/template, jsonnet
	Renderer string `json:"renderer"`
	// path to a template file
	Template *RawModule `json:"template"`
	// parameter
	Config map[string]any `json:"config"`
	// Transform parameter
	Transform       string  `json:"transform"`
	TemplateModule  *Module `json:"-"`
	TransformModule *Module `json:"-"`
}

type Output struct {
	ID string `json:"id"`
	// text/template, html/template, jsonnet
	Renderer string `json:"renderer"`
	// path to a template file
	Template *Module `json:"template"`
	// parameter
	Config map[string]any `json:"config"`
	// Transform parameter
	Transform       string  `json:"transform"`
	TemplateModule  *Module `json:"-"`
	TransformModule *Module `json:"-"`
}

type Target struct {
	ID             string                    `json:"id,omitempty"`
	LintFiles      []*ModuleGlob             `json:"lint_files,omitempty"`
	Modules        []*ModuleGlob             `json:"modules,omitempty"`
	ModuleArchives map[string]*ModuleArchive `json:"module_archives,omitempty"`
	DataFiles      []string                  `json:"data_files,omitempty"`
}

type RawTarget struct {
	ID        string           `json:"id,omitempty"`
	LintGlobs []*LintGlob      `json:"lint_files"`
	Modules   []*RawModuleGlob `json:"modules"`
	DataFiles []string         `json:"data_files"`
}

type LintGlob struct {
	Glob   string         `json:"path"`
	Config map[string]any `json:"config"`
}

func (lg *LintGlob) UnmarshalJSON(b []byte) error {
	rm := &RawModuleGlob{}
	if err := json.Unmarshal(b, rm); err != nil {
		return err //nolint:wrapcheck
	}
	lg.Config = rm.Config
	lg.Glob = rm.Glob
	return nil
}

func (lg *LintGlob) ToModule(cfgDir string) *ModuleGlob {
	p := strings.TrimPrefix(lg.Glob, "!")
	a := filepath.FromSlash(p)
	if !filepath.IsAbs(a) {
		a = filepath.Join(cfgDir, a)
	}
	return &ModuleGlob{
		Path: &domain.Path{
			Raw: p,
			Abs: a,
		},
		Config:   lg.Config,
		Excluded: p != lg.Glob,
	}
}

type RawModule struct {
	Path   string         `json:"path"`
	Config map[string]any `json:"config"`
}

func (rm *RawModule) ToModule() *Module {
	return &Module{
		ID:        rm.Path,
		SlashPath: rm.Path,
		Config:    rm.Config,
	}
}

type RawModuleGlob struct {
	Glob   string         `json:"path"`
	Config map[string]any `json:"config"`
}

type LintFile struct {
	ID     string         `json:"id,omitempty"`
	Path   string         `json:"path,omitempty"`
	Config map[string]any `json:"config,omitempty"`
}

func (rm *RawModuleGlob) UnmarshalJSON(b []byte) error {
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
