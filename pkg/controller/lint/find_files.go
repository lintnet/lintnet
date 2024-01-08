package lint

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/lintnet/lintnet/pkg/config"
	"github.com/lintnet/lintnet/pkg/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"golang.org/x/exp/maps"
)

var ignoreDirs = map[string]struct{}{
	"node_modules": {},
	".git":         {},
}

type LintFile struct { //nolint:revive
	Path       string
	ModulePath string
	Config     map[string]any
}

func (c *Controller) findTarget(logE *logrus.Entry, target *config.Target, rootDir string) (*Target, error) {
	lintFiles, err := c.findFilesFromModules(target.LintFiles, "")
	if err != nil {
		return nil, err
	}
	logE.WithFields(logrus.Fields{
		"lint_globs": log.JSON(target.LintFiles),
		"lint_files": log.JSON(lintFiles),
	}).Debug("found lint files")
	dataFiles, err := c.findFilesFromPaths(target.DataFiles)
	if err != nil {
		return nil, err
	}
	logE.WithFields(logrus.Fields{
		"data_globs": log.JSON(target.DataFiles),
		"data_files": log.JSON(dataFiles),
	}).Debug("found data files")
	modules, err := c.findFilesFromModules(target.Modules, rootDir)
	if err != nil {
		return nil, err
	}
	logE.WithFields(logrus.Fields{
		"module_globs": log.JSON(target.Modules),
		"modules":      log.JSON(modules),
	}).Debug("found modules")
	lintFiles = append(lintFiles, modules...)
	return &Target{
		Combine:   target.Combine,
		LintFiles: lintFiles,
		DataFiles: dataFiles,
	}, nil
}

type Target struct {
	ID        string
	Combine   bool
	LintFiles []*config.LintFile
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

func filterTarget(target *Target, filePaths []string) *Target {
	newTarget := &Target{}
	for _, lintFile := range target.LintFiles {
		for _, filePath := range filePaths {
			if lintFile.Path == filePath {
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

func (c *Controller) findFiles(logE *logrus.Entry, cfg *config.Config, rootDir string) ([]*Target, error) {
	if len(cfg.Targets) == 0 {
		return nil, nil
	}

	targets := make([]*Target, len(cfg.Targets))
	for i, target := range cfg.Targets {
		t, err := c.findTarget(logE, target, rootDir)
		if err != nil {
			return nil, err
		}
		t.ID = target.ID
		targets[i] = t
	}
	return targets, nil
}

func (c *Controller) findFilesFromModule(m *config.ModuleGlob, rootDir string, matchFiles map[string][]*config.LintFile) error {
	if m.Excluded {
		pattern := filepath.Join(rootDir, filepath.FromSlash(m.SlashPath))
		for file := range matchFiles {
			matched, err := doublestar.Match(pattern, file)
			if err != nil {
				return fmt.Errorf("check file match: %w", err)
			}
			if matched {
				delete(matchFiles, file)
			}
		}
		return nil
	}
	matches := map[string]struct{}{}
	if err := doublestar.GlobWalk(afero.NewIOFS(c.fs), filepath.Join(rootDir, filepath.FromSlash(m.SlashPath)), func(path string, d fs.DirEntry) error {
		if _, ok := ignoreDirs[d.Name()]; ok {
			return fs.SkipDir
		}
		if strings.HasSuffix(d.Name(), "_test.jsonnet") {
			return nil
		}
		matches[path] = struct{}{}
		return nil
	}, doublestar.WithNoFollow(), doublestar.WithFilesOnly()); err != nil {
		return fmt.Errorf("search files: %w", err)
	}
	for file := range matches {
		relPath, err := filepath.Rel(rootDir, file)
		if err != nil {
			return fmt.Errorf("get a relative path from the root directory to a module: %w", err)
		}
		var id string
		if m.Archive == nil {
			id = filepath.ToSlash(file)
		} else {
			id = filepath.ToSlash(relPath) // TODO add tag
		}
		matchFiles[file] = append(matchFiles[file], &config.LintFile{
			ID:     id,
			Path:   file,
			Config: m.Config,
		})
	}
	return nil
}

func (c *Controller) findFilesFromModules(modules []*config.ModuleGlob, rootDir string) ([]*config.LintFile, error) {
	matchFiles := map[string][]*config.LintFile{}
	for _, m := range modules {
		if err := c.findFilesFromModule(m, rootDir, matchFiles); err != nil {
			return nil, err
		}
	}
	arr := []*config.LintFile{}
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
		if err := doublestar.GlobWalk(afero.NewIOFS(c.fs), line, func(path string, d fs.DirEntry) error {
			if _, ok := ignoreDirs[d.Name()]; ok {
				return fs.SkipDir
			}
			matchFiles[path] = struct{}{}
			return nil
		}, doublestar.WithNoFollow(), doublestar.WithFilesOnly()); err != nil {
			return nil, fmt.Errorf("search files: %w", err)
		}
	}
	return maps.Keys(matchFiles), nil
}
