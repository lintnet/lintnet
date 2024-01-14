package filefind

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/lintnet/lintnet/pkg/config"
	"github.com/lintnet/lintnet/pkg/domain"
	"github.com/lintnet/lintnet/pkg/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"golang.org/x/exp/maps"
)

type FileFinder struct {
	fs afero.Fs
}

func NewFileFinder(fs afero.Fs) *FileFinder {
	return &FileFinder{
		fs: fs,
	}
}

func (f *FileFinder) Find(logE *logrus.Entry, cfg *config.Config, rootDir, cfgDir string) ([]*domain.Target, error) {
	if len(cfg.Targets) == 0 {
		return nil, nil
	}

	targets := make([]*domain.Target, len(cfg.Targets))
	for i, target := range cfg.Targets {
		t, err := f.findTarget(logE, target, rootDir, cfgDir, cfg.IgnoredPatterns)
		if err != nil {
			return nil, err
		}
		t.ID = target.ID
		targets[i] = t
	}
	return targets, nil
}

func (f *FileFinder) findTarget(logE *logrus.Entry, target *config.Target, rootDir, cfgDir string, ignorePatterns []string) (*domain.Target, error) {
	lintFiles, err := f.findFilesFromModules(target.LintFiles, "", ignorePatterns)
	if err != nil {
		return nil, err
	}
	logE.WithFields(logrus.Fields{
		"lint_globs": log.JSON(target.LintFiles),
		"lint_files": log.JSON(lintFiles),
	}).Debug("found lint files")
	dataFiles, err := f.findFilesFromPaths(target.DataFiles, cfgDir, ignorePatterns)
	if err != nil {
		return nil, err
	}
	logE.WithFields(logrus.Fields{
		"data_globs": log.JSON(target.DataFiles),
		"data_files": log.JSON(dataFiles),
	}).Debug("found data files")
	modules, err := f.findFilesFromModules(target.Modules, rootDir, ignorePatterns)
	if err != nil {
		return nil, err
	}
	logE.WithFields(logrus.Fields{
		"module_globs": log.JSON(target.Modules),
		"modules":      log.JSON(modules),
	}).Debug("found modules")
	lintFiles = append(lintFiles, modules...)
	return &domain.Target{
		LintFiles: lintFiles,
		DataFiles: dataFiles,
	}, nil
}

func (f *FileFinder) findFilesFromModule(m *config.ModuleGlob, rootDir string, matchFiles map[string][]*config.LintFile, ignorePatterns []string) error { //nolint:cyclop
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
	if err := doublestar.GlobWalk(afero.NewIOFS(f.fs), filepath.Join(rootDir, filepath.FromSlash(m.SlashPath)), func(path string, d fs.DirEntry) error {
		if err := ignorePath(path, ignorePatterns); err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if strings.HasSuffix(d.Name(), "_test.jsonnet") {
			return nil
		}
		matches[path] = struct{}{}
		return nil
	}, doublestar.WithNoFollow()); err != nil {
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

func (f *FileFinder) findFilesFromModules(modules []*config.ModuleGlob, rootDir string, ignorePatterns []string) ([]*config.LintFile, error) {
	matchFiles := map[string][]*config.LintFile{}
	for _, m := range modules {
		if err := f.findFilesFromModule(m, rootDir, matchFiles, ignorePatterns); err != nil {
			return nil, err
		}
	}
	arr := []*config.LintFile{}
	for _, m := range matchFiles {
		arr = append(arr, m...)
	}
	return arr, nil
}

func (f *FileFinder) excludeFiles(pattern, cfgDir string, matchFiles map[string]*domain.Path) error {
	for file := range matchFiles {
		if !filepath.IsAbs(pattern) {
			pattern = filepath.Join(cfgDir, pattern)
		}
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

func (f *FileFinder) findFilesFromPath(line, cfgDir string, matchFiles map[string]*domain.Path, ignoredPatterns []string) error {
	if pattern := strings.TrimPrefix(line, "!"); pattern != line {
		if err := f.excludeFiles(pattern, cfgDir, matchFiles); err != nil {
			return err
		}
		return nil
	}
	isAbs := filepath.IsAbs(line)
	if !isAbs {
		line = filepath.Join(cfgDir, line)
	}
	if err := doublestar.GlobWalk(afero.NewIOFS(f.fs), line, func(path string, d fs.DirEntry) error {
		p := &domain.Path{
			Raw: path,
			Abs: path,
		}
		if !isAbs {
			a, err := filepath.Rel(cfgDir, path)
			if err == nil {
				p.Raw = a
			}
		}
		if err := ignorePath(p.Raw, ignoredPatterns); err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		matchFiles[path] = p
		return nil
	}, doublestar.WithNoFollow()); err != nil {
		return fmt.Errorf("search files: %w", err)
	}
	return nil
}

func (f *FileFinder) findFilesFromPaths(files []string, cfgDir string, ignoredPatterns []string) ([]*domain.Path, error) {
	matchFiles := map[string]*domain.Path{}
	for _, line := range files {
		if err := f.findFilesFromPath(line, cfgDir, matchFiles, ignoredPatterns); err != nil {
			return nil, err
		}
	}
	return maps.Values(matchFiles), nil
}

func ignorePath(path string, ignorePatterns []string) error {
	for _, ignoredDir := range ignorePatterns {
		f, err := doublestar.PathMatch(ignoredDir, path)
		if err != nil {
			return fmt.Errorf("check if a path is included in a ignored directory: %w", err)
		}
		if f {
			return fs.SkipDir
		}
	}
	return nil
}