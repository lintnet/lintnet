package config

import (
	"errors"
	"path/filepath"
	"regexp"
	"strings"
)

type Module struct {
	ID        string
	Archive   *ModuleArchive
	SlashPath string
	Config    map[string]interface{}
}

type ModuleGlob struct {
	ID        string
	SlashPath string
	Archive   *ModuleArchive
	Config    map[string]interface{}
	Excluded  bool
}

func (m *Module) FilePath() string {
	if m.Archive == nil {
		return filepath.FromSlash(m.SlashPath)
	}
	return filepath.Join(m.Archive.FilePath(), filepath.FromSlash(m.SlashPath))
}

type ModuleArchive struct {
	ID        string
	Type      string
	Host      string
	RepoOwner string
	RepoName  string
	Ref       string
	Tag       string
}

var fullCommitHashPattern = regexp.MustCompile("[a-fA-F0-9]{40}")

func validateRef(ref string) error {
	if fullCommitHashPattern.MatchString(ref) {
		return nil
	}
	return errors.New("ref must be full commit hash")
}

func (m *ModuleArchive) FilePath() string {
	return filepath.Join(m.Host, m.RepoOwner, m.RepoName, m.Ref)
}

func ParseImport(line string) (*Module, error) {
	mg, err := ParseModuleLine(line)
	if err != nil {
		return nil, err
	}
	return &Module{
		ID:        mg.ID,
		Archive:   mg.Archive,
		SlashPath: mg.SlashPath,
	}, nil
}

func ParseModuleLine(line string) (*ModuleGlob, error) {
	// github.com/<repo owner>/<repo name>/<path>@<commit hash>[:<tag>]
	line = strings.TrimSpace(line)
	excluded := false
	if l := strings.TrimPrefix(line, "!"); l != line {
		excluded = true
		line = strings.TrimSpace(l)
	}
	elems := strings.Split(line, "/")
	if len(elems) < 4 { //nolint:gomnd
		return nil, errors.New("line is invalid")
	}
	if elems[0] != "github.com" {
		return nil, errors.New("module must start with 'github.com/'")
	}
	pathAndRefAndTag := strings.Join(elems[3:], "/")
	path, refAndTag, ok := strings.Cut(pathAndRefAndTag, "@")
	if !ok {
		return nil, errors.New("ref is required")
	}
	ref, tag, _ := strings.Cut(refAndTag, ":")
	if err := validateRef(ref); err != nil {
		return nil, err
	}
	return &ModuleGlob{
		ID:        line,
		SlashPath: strings.Join(append(elems[:3], ref, path), "/"),
		Archive: &ModuleArchive{
			ID:        strings.Join(append(elems[:3], refAndTag), "/"),
			Type:      "github",
			Host:      "github.com",
			RepoOwner: elems[1],
			RepoName:  elems[2],
			Ref:       ref,
			Tag:       tag,
		},
		Excluded: excluded,
	}, nil
}
