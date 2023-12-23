package lint

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/lintnet/lintnet/pkg/config"
	"github.com/lintnet/lintnet/pkg/jsonnet"
	"github.com/lintnet/lintnet/pkg/module"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
	"golang.org/x/exp/maps"
)

type LintFile struct { //nolint:revive
	Path       string
	ModulePath string
	Imports    map[string]string
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

func (c *Controller) findFilesFromPaths(files string) ([]string, error) {
	lines := strings.Split(files, "\n")
	matchFiles := map[string]struct{}{}
	for _, line := range lines {
		if strings.HasPrefix(line, "#") {
			// ignore comments
			continue
		}
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
		matches, err := doublestar.Glob(afero.NewIOFS(c.fs), line)
		if err != nil {
			return nil, fmt.Errorf("search files: %w", err)
		}
		for _, file := range matches {
			matchFiles[file] = struct{}{}
		}
	}
	return maps.Keys(matchFiles), nil
}

func (c *Controller) convertStringsToTargets(logE *logrus.Entry, ruleBaseDir string, dataFiles []string) ([]*Target, error) {
	lintFiles, err := jsonnet.FindFiles(logE, c.fs, ruleBaseDir)
	if err != nil {
		return nil, err
	}
	arr := make([]*LintFile, len(lintFiles))
	for i, lintFile := range lintFiles {
		arr[i] = &LintFile{
			Path: lintFile,
		}
	}
	return []*Target{
		{
			LintFiles: arr,
			DataFiles: dataFiles,
		},
	}, nil
}

func (c *Controller) findFiles(logE *logrus.Entry, cfg *config.Config, modulesList [][]*module.Module, ruleBaseDir string, dataFiles []string, rootDir string) ([]*Target, error) {
	if ruleBaseDir != "" {
		return c.convertStringsToTargets(logE, ruleBaseDir, dataFiles)
	}
	if len(cfg.Targets) == 0 {
		return c.convertStringsToTargets(logE, "lintnet", dataFiles)
	}

	targets := make([]*Target, len(cfg.Targets))
	for i, target := range cfg.Targets {
		modules := modulesList[i]
		t, err := c.findTarget(target, modules, rootDir)
		if err != nil {
			return nil, err
		}
		targets[i] = t
	}
	return targets, nil
}

func (c *Controller) readJsonnets(filePaths []*LintFile) (map[string]jsonnet.Node, error) {
	jsonnetAsts := make(map[string]jsonnet.Node, len(filePaths))
	for _, filePath := range filePaths {
		ja, err := jsonnet.Read(c.fs, filePath.Path)
		if err != nil {
			return nil, logerr.WithFields(err, logrus.Fields{ //nolint:wrapcheck
				"file_path": filePath,
			})
		}
		if filePath.ModulePath != "" {
			jsonnetAsts[filePath.ModulePath] = ja
			continue
		}
		jsonnetAsts[filePath.Path] = ja
	}
	return jsonnetAsts, nil
}

func (c *Controller) evaluate(tla string, jsonnetAsts map[string]jsonnet.Node) map[string]*JsonnetEvaluateResult {
	vm := jsonnet.NewVM(tla, c.importer)

	results := make(map[string]*JsonnetEvaluateResult, len(jsonnetAsts))
	for k, ja := range jsonnetAsts {
		result, err := vm.Evaluate(ja)
		if err != nil {
			results[k] = &JsonnetEvaluateResult{
				Error: err.Error(),
			}
			continue
		}
		results[k] = &JsonnetEvaluateResult{
			Result: result,
		}
	}
	return results
}
