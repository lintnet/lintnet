package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

type RawModule struct {
	Glob   string          `json:"path"`
	Files  []*LintGlobFile `json:"files,omitempty"`
	Config map[string]any  `json:"config"`
}

func (rm *RawModule) Parse() (*ModuleGlob, error) {
	m, err := ParseModuleLine(rm.Glob)
	if err != nil {
		return nil, fmt.Errorf("parse a module path: %w", err)
	}
	m.Config = rm.Config
	p := path.Clean(m.SlashPath)
	if strings.HasPrefix(p, "..") {
		return nil, fmt.Errorf("'..' is forbidden: %w", err)
	}
	m.SlashPath = p
	for _, f := range rm.Files {
		f.Clean()
		if strings.HasPrefix(f.Path, "..") {
			return nil, fmt.Errorf("'..' is forbidden: %w", err)
		}
	}
	m.Files = rm.Files
	return m, nil
}

func (rm *RawModule) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err == nil {
		rm.Glob = s
		return nil
	}
	a := struct {
		Path   string          `json:"path"`
		Config map[string]any  `json:"config"`
		Files  []*LintGlobFile `json:"files"`
	}{}
	if err := json.Unmarshal(b, &a); err != nil {
		return err
	}
	rm.Glob = a.Path
	rm.Config = a.Config
	rm.Files = a.Files
	return nil
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
	Files     []*LintGlobFile        `json:"files,omitempty"`
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
	line, excluded := parseNegationOperator(strings.TrimSpace(line))
	elems := strings.Split(line, "/")
	if len(elems) < 3 { //nolint:mnd
		return nil, errors.New("line is invalid")
	}
	moduleType, host, repoOwner, repoName := elems[0], elems[1], elems[2], elems[3]
	if moduleType != "github_archive" {
		return nil, errors.New("unsupported module type")
	}
	if host != "github.com" {
		return nil, errors.New("module host must be 'github.com'")
	}
	if len(elems) == 4 { //nolint:mnd
		// <type>/github.com/<repo owner>/<repo name>@<commit hash>[:<tag>]
		repoName, refAndTag, ok := strings.Cut(repoName, "@")
		if !ok {
			return nil, errors.New("ref is required")
		}
		ref, tag, err := parseRefAndTag(refAndTag)
		if err != nil {
			return nil, err
		}
		return &ModuleGlob{
			SlashPath: strings.Join(append(elems[:3], repoName, ref), "/"),
			Archive: &ModuleArchive{
				Type:      moduleType,
				Host:      host,
				RepoOwner: repoOwner,
				RepoName:  repoName,
				Ref:       ref,
				Tag:       tag,
			},
			Excluded: excluded,
		}, nil
	}
	pathAndRefAndTag := strings.Join(elems[4:], "/")
	path, refAndTag, ok := strings.Cut(pathAndRefAndTag, "@")
	if !ok {
		return nil, errors.New("ref is required")
	}
	ref, tag, err := parseRefAndTag(refAndTag)
	if err != nil {
		return nil, err
	}
	return &ModuleGlob{
		SlashPath: strings.Join(append(elems[:4], ref, path), "/"),
		Archive: &ModuleArchive{
			Type:      moduleType,
			Host:      host,
			RepoOwner: repoOwner,
			RepoName:  repoName,
			Ref:       ref,
			Tag:       tag,
		},
		Excluded: excluded,
	}, nil
}

func parseRefAndTag(refAndTag string) (string, string, error) {
	ref, tag, _ := strings.Cut(refAndTag, ":")
	return ref, tag, validateRef(ref)
}

type LintGlobFile struct {
	Path     string         `json:"path"`
	Config   map[string]any `json:"config,omitempty"`
	Excluded bool           `json:"-"`
}

func (lg *LintGlobFile) UnmarshalJSON(b []byte) error {
	var p string
	if err := json.Unmarshal(b, &p); err == nil {
		lg.Path = p
		return nil
	}
	rm := struct {
		Path   string         `json:"path"`
		Config map[string]any `json:"config,omitempty"`
	}{}

	if err := json.Unmarshal(b, &rm); err != nil {
		return err //nolint:wrapcheck
	}
	lg.Config = rm.Config
	lg.Path = rm.Path
	return nil
}

func (lg *LintGlobFile) Clean() {
	p, excluded := parseNegationOperator(lg.Path)
	lg.Excluded = excluded
	lg.Path = path.Clean(p)
}
