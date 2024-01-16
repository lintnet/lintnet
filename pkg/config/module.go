package config

import (
	"fmt"
	"path/filepath"

	"github.com/lintnet/lintnet/pkg/domain"
)

type Module struct {
	ID        string
	Archive   *ModuleArchive
	SlashPath string
	Config    map[string]interface{}
}

type ModuleGlob struct {
	Path     *domain.Path           `json:"path,omitempty"`
	Archive  *ModuleArchive         `json:"archive,omitempty"`
	Config   map[string]interface{} `json:"config,omitempty"`
	Excluded bool                   `json:"excluded,omitempty"`
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

func (m *ModuleArchive) FilePath() string {
	return filepath.Join(m.Type, m.Host, m.RepoOwner, m.RepoName, m.Ref)
}
