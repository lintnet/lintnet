package config

import (
	"encoding/json"
	"path/filepath"
	"strings"
)

type LintFile struct {
	ID     string         `json:"id,omitempty"`
	Path   string         `json:"path,omitempty"`
	Config map[string]any `json:"config,omitempty"`
}

type LintGlob struct {
	// Glob is either an absolute path or a relative path from configuration file path
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
	excluded := p != lg.Glob
	p = filepath.Clean(p)
	return &ModuleGlob{
		SlashPath: p,
		Config:    lg.Config,
		Excluded:  excluded,
	}
}
