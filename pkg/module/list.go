package module

import (
	"errors"
	"path/filepath"
	"regexp"
	"strings"
)

type Module struct {
	ID      string
	Archive *Archive
	Path    string
	Param   map[string]interface{}
}

type Glob struct {
	ID        string
	SlashPath string
	Archive   *Archive
	Glob      string
	Param     map[string]interface{}
	Excluded  bool
}

func (m *Module) FilePath() string {
	if m.Archive == nil {
		return filepath.FromSlash(m.Path)
	}
	return filepath.Join(m.Archive.FilePath(), filepath.FromSlash(m.Path))
}

type Archive struct {
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

func (m *Archive) FilePath() string {
	return filepath.Join(m.Host, m.RepoOwner, m.RepoName, m.Ref)
}

func ParseModuleLine(line string) (*Glob, error) {
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
	return &Glob{
		ID:        line,
		SlashPath: strings.Join(append(elems[:3], ref, path), "/"),
		Archive: &Archive{
			ID:        strings.Join(append(elems[:3], refAndTag), "/"),
			Type:      "github",
			Host:      "github.com",
			RepoOwner: elems[1],
			RepoName:  elems[2],
			Ref:       ref,
			Tag:       tag,
		},
		Glob:     path,
		Excluded: excluded,
	}, nil
}
