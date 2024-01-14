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

type LintFile struct { //nolint:revive
	Path       string
	ModulePath string
	Config     map[string]any
}

type Target struct {
	ID        string
	LintFiles []*config.LintFile
	DataFiles Paths
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
			if dataFile.Abs == filePath {
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

func (c *Controller) filterTargetsByDataRootDir(logE *logrus.Entry, param *ParamLint, targets []*Target) error {
	for _, target := range targets {
		if err := c.filterTargetByDataRootDir(logE, param, target); err != nil {
			return err
		}
	}
	return nil
}

func (c *Controller) filterTargetByDataRootDir(logE *logrus.Entry, param *ParamLint, target *Target) error {
	arr := make([]*Path, 0, len(target.DataFiles))
	for _, dataFile := range target.DataFiles {
		if c.filterFileByDataRootDir(logE, param, dataFile.Abs) {
			arr = append(arr, dataFile)
		} else {
			logE.WithField("data_file", dataFile).Warn("this data file is ignored because this is out of the data root directory")
		}
	}
	target.DataFiles = arr
	return nil
}

func (c *Controller) filterFileByDataRootDir(logE *logrus.Entry, param *ParamLint, dataFile string) bool {
	p := dataFile
	if a, err := filepath.Rel(param.DataRootDir, p); err != nil {
		logE.WithError(err).Warn("get a relative path")
	} else if !strings.HasPrefix(a, "..") {
		return true
	}
	for _, c := range param.FilePaths {
		b, err := filepath.Rel(c, dataFile)
		if err != nil {
			logE.WithError(err).Warn("get a relative path")
			continue
		}
		if b == "." {
			return true
		}
	}
	return false
}

func (c *Controller) findFiles(logE *logrus.Entry, cfg *config.Config, rootDir, cfgDir string) ([]*Target, error) {
	finder := &FileFinder{
		fs: c.fs,
	}
	return finder.Find(logE, cfg, rootDir, cfgDir)
}

type FileFinder struct {
	fs afero.Fs
}

func (f *FileFinder) Find(logE *logrus.Entry, cfg *config.Config, rootDir, cfgDir string) ([]*Target, error) {
	if len(cfg.Targets) == 0 {
		return nil, nil
	}

	targets := make([]*Target, len(cfg.Targets))
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

func (f *FileFinder) findTarget(logE *logrus.Entry, target *config.Target, rootDir, cfgDir string, ignorePatterns []string) (*Target, error) {
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
	return &Target{
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

func (f *FileFinder) excludeFiles(pattern, cfgDir string, matchFiles map[string]*Path) error {
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

func (f *FileFinder) findFilesFromPath(line, cfgDir string, matchFiles map[string]*Path, ignoredPatterns []string) error {
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
		p := &Path{
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

func (f *FileFinder) findFilesFromPaths(files []string, cfgDir string, ignoredPatterns []string) ([]*Path, error) {
	matchFiles := map[string]*Path{}
	for _, line := range files {
		if err := f.findFilesFromPath(line, cfgDir, matchFiles, ignoredPatterns); err != nil {
			return nil, err
		}
	}
	return maps.Values(matchFiles), nil
}
