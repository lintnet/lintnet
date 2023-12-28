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
	Param      map[string]any
}

func (c *Controller) findTarget(target *config.Target, modules []*module.Module, rootDir string) (*Target, error) {
	lintFiles, err := c.findFilesFromModules(target.LintFiles)
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
			Path:  b.Path,
			Param: b.Param,
		})
	}
	for _, mod := range modules {
		a = append(a, &LintFile{
			ModulePath: path.Join(mod.ID(), mod.Path),
			Path:       filepath.Join(rootDir, filepath.FromSlash(mod.ID()), filepath.FromSlash(mod.Path)),
			Param:      mod.Param,
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

func filterTargets(targets []*Target, filePaths []string) []*Target {
	newTargets := make([]*Target, 0, len(targets))
	for _, target := range targets {
		newTarget := filterTarget(target, filePaths)
		if len(newTarget.LintFiles) > 0 {
			newTargets = append(newTargets, newTarget)
		}
	}
	return newTargets
}

func filterTarget(target *Target, filePaths []string) *Target { //nolint:cyclop
	newTarget := &Target{}
	for _, lintFile := range target.LintFiles {
		for _, filePath := range filePaths {
			if lintFile.Path == filePath || lintFile.ModulePath == filePath {
				newTarget.LintFiles = append(newTarget.LintFiles, lintFile)
				break
			}
		}
	}
	lintChanged := false
	if len(newTarget.LintFiles) > 0 {
		newTarget.DataFiles = target.DataFiles
		lintChanged = true
	}
	dataChanged := false
	for _, dataFile := range target.DataFiles {
		for _, filePath := range filePaths {
			if dataFile == filePath {
				dataChanged = true
				if !lintChanged {
					newTarget.DataFiles = append(newTarget.DataFiles, dataFile)
				}
			}
		}
	}
	if dataChanged {
		newTarget.LintFiles = target.LintFiles
	}
	return newTarget
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

func (c *Controller) findFilesFromModules(modules []*config.Module) ([]*config.Module, error) {
	matchFiles := map[string][]*config.Module{}
	for _, m := range modules {
		if pattern := strings.TrimPrefix(m.Path, "!"); pattern != m.Path {
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
		matches, err := doublestar.Glob(afero.NewIOFS(c.fs), m.Path, doublestar.WithFilesOnly())
		if err != nil {
			return nil, fmt.Errorf("search files: %w", err)
		}
		for _, file := range matches {
			matchFiles[file] = append(matchFiles[file], &config.Module{
				Path:  file,
				Param: m.Param,
			})
		}
	}
	arr := []*config.Module{}
	for _, m := range matchFiles {
		arr = append(arr, m...)
	}
	return arr, nil
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
