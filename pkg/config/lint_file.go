package config

import (
	"encoding/json"
	"path"
	"strings"
)

type LintFile struct {
	ID     string         `json:"id,omitempty"`
	Path   string         `json:"path,omitempty"`
	Config map[string]any `json:"config,omitempty"`
}

type LintGlob struct {
	// Glob is either an absolute path or a relative path from configuration file path
	Glob     string         `json:"path"`
	Config   map[string]any `json:"config,omitempty"`
	Excluded bool           `json:"-"`
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

func (lg *LintGlob) Clean() {
	p, excluded := parseNegationOperator(lg.Glob)
	lg.Excluded = excluded
	lg.Glob = path.Clean(p)
}

func parseNegationOperator(p string) (string, bool) {
	a := strings.TrimPrefix(p, "!")
	return strings.TrimSpace(a), a != p
}
