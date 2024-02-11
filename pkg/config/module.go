package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

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
	p := filepath.Clean(m.SlashPath)
	if strings.HasPrefix(p, "..") {
		return nil, fmt.Errorf("'..' is forbidden: %w", err)
	}
	m.SlashPath = p
	return m, nil
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

type Module struct {
	Archive   *ModuleArchive
	SlashPath string
	Config    map[string]interface{}
}

type ModuleGlob struct {
	SlashPath string                 `json:"slash_path,omitempty"`
	Archive   *ModuleArchive         `json:"archive,omitempty"`
	Config    map[string]interface{} `json:"config,omitempty"`
	Excluded  bool                   `json:"excluded,omitempty"`
}

func (m *Module) FilePath() string {
	if m.Archive == nil {
		return filepath.FromSlash(m.SlashPath)
	}
	return filepath.Join(m.Archive.FilePath(), filepath.FromSlash(m.SlashPath))
}

type ModuleArchive struct {
	Type      string `json:"type,omitempty"`
	Host      string `json:"host,omitempty"`
	RepoOwner string `json:"repo_owner,omitempty"`
	RepoName  string `json:"repo_name,omitempty"`
	Ref       string `json:"ref,omitempty"`
	Tag       string `json:"tag,omitempty"`
}

// String returns a human readable string.
// This is different from file path.
// This is used for log.
func (m *ModuleArchive) String() string {
	a := fmt.Sprintf("%s/%s/%s/%s/%s", m.Type, m.Host, m.RepoOwner, m.RepoName, m.Ref)
	if m.Tag != "" {
		a = fmt.Sprintf("%s:%s", a, m.Tag)
	}
	return a
}

var fullCommitHashPattern = regexp.MustCompile("[a-fA-F0-9]{40}")

func validateRef(ref string) error {
	if fullCommitHashPattern.MatchString(ref) {
		return nil
	}
	return errors.New("ref must be full commit hash")
}

func (m *ModuleArchive) FilePath() string {
	return filepath.Join(m.Type, m.Host, m.RepoOwner, m.RepoName, m.Ref)
}

func ParseImport(line string) (*Module, error) {
	mg, err := ParseModuleLine(line)
	if err != nil {
		return nil, err
	}
	return &Module{
		Archive:   mg.Archive,
		SlashPath: mg.SlashPath,
	}, nil
}

func ParseModuleLine(line string) (*ModuleGlob, error) {
	// <type>/github.com/<repo owner>/<repo name>/<path>@<commit hash>[:<tag>]
	line = strings.TrimSpace(line)
	excluded := false
	if l := strings.TrimPrefix(line, "!"); l != line {
		excluded = true
		line = strings.TrimSpace(l)
	}
	elems := strings.Split(line, "/")
	if len(elems) < 5 { //nolint:gomnd
		return nil, errors.New("line is invalid")
	}
	if elems[0] != "github_archive" {
		return nil, errors.New("unsupported module type")
	}
	if elems[1] != "github.com" {
		return nil, errors.New("module host must be 'github.com'")
	}
	pathAndRefAndTag := strings.Join(elems[4:], "/")
	path, refAndTag, ok := strings.Cut(pathAndRefAndTag, "@")
	if !ok {
		return nil, errors.New("ref is required")
	}
	ref, tag, _ := strings.Cut(refAndTag, ":")
	if err := validateRef(ref); err != nil {
		return nil, err
	}
	return &ModuleGlob{
		SlashPath: strings.Join(append(elems[:4], ref, path), "/"),
		Archive: &ModuleArchive{
			Type:      "github_archive",
			Host:      "github.com",
			RepoOwner: elems[2],
			RepoName:  elems[3],
			Ref:       ref,
			Tag:       tag,
		},
		Excluded: excluded,
	}, nil
}
