package config

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Config struct {
	ErrorLevel string    `json:"error_level" yaml:"error_level"`
	Targets    []*Target `json:"targets"`
	Outputs    []*Output `json:"outputs,omitempty"`
}

type Output struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Renderer string `json:"renderer"`
	Path     string `json:"path"`
	Template string `json:"template"`
}

type Target struct {
	LintFiles []*Module `json:"lint_files" yaml:"lint_files"`
	Modules   []*Module `json:"modules"`
	DataFiles []string  `json:"data_files" yaml:"data_files"`
}

type Module struct {
	Path  string                 `json:"path"`
	Param map[string]interface{} `json:"param"`
}

func (m *Module) UnmarshalJSON(b []byte) error {
	var a interface{}
	if err := json.Unmarshal(b, &a); err != nil {
		return fmt.Errorf("unmarshal as JSON: %w", err)
	}
	switch c := a.(type) {
	case string:
		m.Path = c
		return nil
	case map[string]interface{}:
		p, ok := c["path"]
		if !ok {
			return errors.New("path is required")
		}
		a, ok := p.(string)
		if !ok {
			return errors.New("path must be a string")
		}
		m.Path = a

		param, ok := c["param"]
		if ok {
			a, ok := param.(map[string]interface{})
			if !ok {
				return errors.New("param must be a map[string]interface{}")
			}
			m.Param = a
		}
		return nil
	}
	return errors.New("module must be either string or map[string]interface{}")
}
