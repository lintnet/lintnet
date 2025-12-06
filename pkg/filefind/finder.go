package filefind

import (
	"fmt"
	"io/fs"
	"log/slog"
	"maps"
	"path/filepath"
	"slices"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/lintnet/lintnet/pkg/config"
	"github.com/lintnet/lintnet/pkg/domain"
	"github.com/lintnet/lintnet/pkg/log"
	"github.com/spf13/afero"
)

type Target struct {
	ID        string             `json:"id,omitempty"`
	LintFiles []*config.LintFile `json:"lint_files,omitempty"`
	DataFiles domain.Paths       `json:"data_files,omitempty"`
}

type FileFinder struct {
	fs afero.Fs
}

func NewFileFinder(fs afero.Fs) *FileFinder {
	return &FileFinder{
		fs: fs,
	}
}

func (f *FileFinder) Find(logger *slog.Logger, cfg *config.Config, rootDir, cfgDir string) ([]*Target, error) {
	if len(cfg.Targets) == 0 {
		return nil, nil
	}

	targets := make([]*Target, 0, len(cfg.Targets))
	for _, target := range cfg.Targets {
		ts, err := f.findTarget(logger, target, rootDir, cfgDir, cfg.IgnoredPatterns)
		if err != nil {
			return nil, err
		}
		for _, t := range ts {
			t.ID = target.ID
		}
		targets = append(targets, ts...)
	}
	return targets, nil
}

func (f *FileFinder) findTarget(logger *slog.Logger, target *config.Target, rootDir, cfgDir string, ignorePatterns []string) ([]*Target, error) {
	lintFiles, err := f.findFilesFromLintFiles(logger, target.LintFiles, cfgDir, ignorePatterns)
	if err != nil {
		return nil, err
	}
	for _, lintFile := range lintFiles {
		if !filepath.IsAbs(lintFile.Path) {
			lintFile.Path = filepath.Join(cfgDir, lintFile.Path)
		}
	}
	logger.Debug("found lint files", "lint_globs", log.JSON(target.LintFiles), "lint_files", log.JSON(lintFiles))

	modules, err := f.findFilesFromModules(logger, target.Modules, rootDir, ignorePatterns)
	if err != nil {
		return nil, err
	}
	logger.Debug("found modules", "module_globs", log.JSON(target.Modules), "modules", log.JSON(modules))
	lintFiles = append(lintFiles, modules...)

	dataFiles, err := f.findDataFiles(target.BaseDataPath, target.DataFiles, cfgDir, ignorePatterns)
	if err != nil {
		return nil, err
	}
	logger.Debug("found data files", "data_globs", log.JSON(target.DataFiles), "data_files", log.JSON(dataFiles))

	targets := make([]*Target, len(dataFiles))
	for i, dataFile := range dataFiles {
		targets[i] = &Target{
			LintFiles: lintFiles,
			DataFiles: dataFile,
		}
	}
	return targets, nil
}

func (f *FileFinder) globModuleFiles(logger *slog.Logger, rootDir, pattern string, m *config.ModuleGlob, file *config.LintGlobFile, matches map[string][]*config.LintFile, ignorePatterns []string) error {
	logger.Debug("search module files", "pattern", pattern)
	if err := doublestar.GlobWalk(afero.NewIOFS(f.fs), pattern, func(path string, d fs.DirEntry) error {
		if err := ignorePath(path, ignorePatterns); err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if strings.HasSuffix(d.Name(), "_test.jsonnet") {
			return nil
		}
		link, err := m.Archive.URL(rootDir, path)
		if err != nil {
			return fmt.Errorf("get a module url: %w", err)
		}
		moduleID, err := getModuleID(rootDir, path, m.Archive.Tag)
		if err != nil {
			return err
		}
		lintFile := &config.LintFile{
			ID:     moduleID,
			Config: m.Config,
			Path:   path,
			Link:   link,
		}
		if file != nil && file.Config != nil {
			lintFile.Config = file.Config
		}
		matches[path] = append(matches[path], lintFile)
		return nil
	}, doublestar.WithNoFollow()); err != nil {
		return fmt.Errorf("search files: %w", err)
	}
	return nil
}

func (f *FileFinder) findFilesFromModule(logger *slog.Logger, m *config.ModuleGlob, rootDir string, matchFiles map[string][]*config.LintFile, ignorePatterns []string) error { //nolint:cyclop
	if len(m.Files) == 0 && m.Excluded {
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
	matches := map[string][]*config.LintFile{}
	pattern := filepath.Join(rootDir, filepath.FromSlash(m.SlashPath))
	if len(m.Files) == 0 {
		if err := f.globModuleFiles(logger, rootDir, pattern, m, nil, matches, ignorePatterns); err != nil {
			return err
		}
	}
	for _, file := range m.Files {
		pattern := filepath.Join(rootDir, filepath.FromSlash(m.SlashPath), filepath.FromSlash(file.Path))
		if file.Excluded {
			logger.Debug("check excluded files", "pattern", pattern, "files", slices.Collect(maps.Keys(matches)))
			for matchFile := range matches {
				matched, err := doublestar.Match(pattern, matchFile)
				if err != nil {
					return fmt.Errorf("check file match: %w", err)
				}
				if matched {
					delete(matches, matchFile)
					logger.Debug("exclude a file", "pattern", pattern, "excluded_file", matchFile)
				}
			}
			continue
		}
		if err := f.globModuleFiles(logger, rootDir, pattern, m, file, matches, ignorePatterns); err != nil {
			return err
		}
	}
	if len(matches) == 0 {
		logger.Debug("no file matches", "pattern", m.SlashPath)
	}
	maps.Copy(matchFiles, matches)
	return nil
}

func getModuleID(rootDir, p, tag string) (string, error) {
	a, err := filepath.Rel(rootDir, p)
	if err != nil {
		return "", fmt.Errorf("get a relative path from the root directory to a module file: %w", err)
	}
	moduleID := filepath.ToSlash(a)
	if tag != "" {
		moduleID += ":" + tag
	}
	return moduleID, nil
}

func (f *FileFinder) findFilesFromLintFile(logger *slog.Logger, m *config.LintGlob, rootDir string, matchFiles map[string][]*config.LintFile, ignorePatterns []string) error { //nolint:cyclop
	if m.Excluded {
		pattern := filepath.Join(rootDir, filepath.FromSlash(m.Glob))
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
	if err := doublestar.GlobWalk(afero.NewIOFS(f.fs), filepath.Join(rootDir, filepath.FromSlash(m.Glob)), func(path string, d fs.DirEntry) error {
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
	if len(matches) == 0 {
		logger.Debug("no file matches", "pattern", m.Glob)
	}
	for file := range matches {
		relPath, err := filepath.Rel(rootDir, file)
		if err != nil {
			return fmt.Errorf("get a relative path from the root directory to a module: %w", err)
		}
		id := filepath.ToSlash(relPath)
		matchFiles[file] = append(matchFiles[file], &config.LintFile{
			ID:     id,
			Path:   file,
			Config: m.Config,
		})
	}
	return nil
}

func (f *FileFinder) findFilesFromModules(logger *slog.Logger, modules []*config.ModuleGlob, rootDir string, ignorePatterns []string) ([]*config.LintFile, error) {
	matchFiles := map[string][]*config.LintFile{}
	for _, m := range modules {
		if err := f.findFilesFromModule(logger, m, rootDir, matchFiles, ignorePatterns); err != nil {
			return nil, err
		}
	}
	arr := []*config.LintFile{}
	for _, m := range matchFiles {
		arr = append(arr, m...)
	}
	return arr, nil
}

func (f *FileFinder) findFilesFromLintFiles(logger *slog.Logger, modules []*config.LintGlob, rootDir string, ignorePatterns []string) ([]*config.LintFile, error) {
	matchFiles := map[string][]*config.LintFile{}
	for _, m := range modules {
		if err := f.findFilesFromLintFile(logger, m, rootDir, matchFiles, ignorePatterns); err != nil {
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

func (f *FileFinder) findFilesFromPath(file *config.DataFile, cfgDir string, matchFiles map[string]*domain.Path, ignoredPatterns []string) error {
	if file.Excluded {
		if err := f.excludeFiles(file.Path, cfgDir, matchFiles); err != nil {
			return err
		}
		return nil
	}
	isAbs := filepath.IsAbs(file.Path)
	line := file.Path
	if !isAbs {
		line = filepath.Join(cfgDir, file.Path)
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

func (f *FileFinder) findFilesFromPaths(files []*config.DataFile, cfgDir string, ignoredPatterns []string) ([]*domain.Path, error) {
	matchFiles := map[string]*domain.Path{}
	for _, file := range files {
		if err := f.findFilesFromPath(file, cfgDir, matchFiles, ignoredPatterns); err != nil {
			return nil, err
		}
	}
	return slices.Collect(maps.Values(matchFiles)), nil
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
