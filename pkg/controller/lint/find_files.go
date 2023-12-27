package lint

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/lintnet/lintnet/pkg/config"
	"github.com/lintnet/lintnet/pkg/module"
	"github.com/spf13/afero"
	"golang.org/x/exp/maps"
)

type LintFile struct { //nolint:revive
	Path       string
	ModulePath string
	Param      map[string]interface{}
}

func (c *Controller) findTarget(target *config.Target, modules []*module.Module, rootDir string) (*Target, error) {
	lintFiles, err := c.findFilesFromPaths(target.LintFiles)
	if err != nil {
		return nil, err
	}
	dataFiles, err := c.findFilesFromPaths(target.DataFiles)
	if err != nil {
		return nil, err
	}
	a := make([]*LintFile, 0, len(lintFiles)+len(modules))
	for _, b := range lintFiles {
		a = append(a, &LintFile{
			Path: b,
		})
	}
	for _, mod := range modules {
		a = append(a, &LintFile{
			ModulePath: path.Join(mod.ID(), mod.Path),
			Path:       filepath.Join(rootDir, filepath.FromSlash(mod.ID()), filepath.FromSlash(mod.Path)),
		})
	}
	return &Target{
		LintFiles: a,
		DataFiles: dataFiles,
	}, nil
}

type Target struct {
	LintFiles []*LintFile
	DataFiles []string
}

func (c *Controller) findFiles(cfg *config.Config, modulesList [][]*module.Module, rootDir string) ([]*Target, error) {
	if len(cfg.Targets) == 0 {
		return nil, nil
	}

	targets := make([]*Target, len(cfg.Targets))
	for i, target := range cfg.Targets {
		var modules []*module.Module
		if modulesList != nil {
			modules = modulesList[i]
		}
		t, err := c.findTarget(target, modules, rootDir)
		if err != nil {
			return nil, err
		}
		targets[i] = t
	}
	return targets, nil
}

func (c *Controller) findFilesFromPaths(files []string) ([]string, error) {
	matchFiles := map[string]struct{}{}
	for _, line := range files {
		if pattern := strings.TrimPrefix(line, "!"); pattern != line {
			for file := range matchFiles {
				matched, err := doublestar.Match(pattern, file)
				if err != nil {
					return nil, fmt.Errorf("check file match: %w", err)
				}
				if matched {
					delete(matchFiles, file)
				}
			}
			continue
		}
		matches, err := doublestar.Glob(afero.NewIOFS(c.fs), line, doublestar.WithFilesOnly())
		if err != nil {
			return nil, fmt.Errorf("search files: %w", err)
		}
		for _, file := range matches {
			matchFiles[file] = struct{}{}
		}
	}
	return maps.Keys(matchFiles), nil
}
