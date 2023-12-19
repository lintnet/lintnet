package lint

import (
	"errors"
	"fmt"
	"io/fs"
	"regexp"
	"strings"

	"github.com/google/go-jsonnet"
	"github.com/google/go-jsonnet/ast"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/lintnet/pkg/config"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
)

type LintFile struct { //nolint:revive
	Path    string
	Imports map[string]string
}

func (c *Controller) findTarget(logE *logrus.Entry, target *config.Target) (*Target, error) {
	lintFiles, err := c.findFilesFromPaths(logE, target.LintFiles.SearchType, target.LintFiles.Paths)
	if err != nil {
		return nil, err
	}
	dataFiles, err := c.findFilesFromPaths(logE, target.DataFiles.SearchType, target.DataFiles.Paths)
	if err != nil {
		return nil, err
	}
	a := make([]*LintFile, len(lintFiles))
	for i, b := range lintFiles {
		a[i] = &LintFile{
			Path: b,
		}
	}
	return &Target{
		LintFiles: a,
		DataFiles: dataFiles,
	}, nil
}

func (c *Controller) findFilesbyGlob(paths []*config.Path) ([]string, error) {
	filePaths := make([]string, 0, len(paths))
	for _, p := range paths {
		matches, err := fs.Glob(afero.NewIOFS(c.fs), p.Path)
		if err != nil {
			return nil, fmt.Errorf("search files by glob: %w", err)
		}
		filePaths = append(filePaths, matches...)
	}
	return filePaths, nil
}

func (c *Controller) findFilesByRegexp(logE *logrus.Entry, paths []*config.Path) ([]string, error) {
	patterns := make([]*regexp.Regexp, len(paths))
	for i, p := range paths {
		p, err := regexp.Compile(p.Path)
		if err != nil {
			return nil, fmt.Errorf("compile a regular expression to search files: %w", err)
		}
		patterns[i] = p
	}
	filePaths := make([]string, 0, len(paths))
	if err := fs.WalkDir(afero.NewIOFS(c.fs), "", func(p string, dirEntry fs.DirEntry, e error) error {
		if e != nil {
			logE.WithError(e).Warn("error occurred while searching files")
			return nil
		}
		for _, pattern := range patterns {
			if pattern.MatchString(p) {
				filePaths = append(filePaths, p)
			}
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("search files: %w", err)
	}
	return filePaths, nil
}

func (c *Controller) findFilesByEqual(paths []*config.Path) ([]string, error) {
	filePaths := make([]string, len(paths))
	for i, p := range paths {
		filePaths[i] = p.Path
	}
	return filePaths, nil
}

type Target struct {
	LintFiles []*LintFile
	DataFiles []string
}

func (c *Controller) findFilesFromPaths(logE *logrus.Entry, searchType string, paths []*config.Path) ([]string, error) {
	switch searchType {
	case "equal":
		return c.findFilesByEqual(paths)
	case "glob":
		return c.findFilesbyGlob(paths)
	case "regexp":
		return c.findFilesByRegexp(logE, paths)
	default:
		return nil, logerr.WithFields(errors.New("search_type is invalid"), logrus.Fields{ //nolint:wrapcheck
			"search_type": searchType,
		})
	}
}

func (c *Controller) convertStringsToTargets(logE *logrus.Entry, ruleBaseDir string, dataFiles []string) ([]*Target, error) {
	lintFiles, err := c.findJsonnetFromBaseDir(logE, ruleBaseDir)
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

func (c *Controller) findFiles(logE *logrus.Entry, cfg *config.Config, ruleBaseDir string, dataFiles []string) ([]*Target, error) {
	if ruleBaseDir != "" {
		return c.convertStringsToTargets(logE, ruleBaseDir, dataFiles)
	}
	if len(cfg.Targets) == 0 {
		return c.convertStringsToTargets(logE, "lintnet", dataFiles)
	}

	targets := make([]*Target, len(cfg.Targets))
	for i, target := range cfg.Targets {
		t, err := c.findTarget(logE, target)
		if err != nil {
			return nil, err
		}
		targets[i] = t
	}
	return targets, nil
}

func (c *Controller) findJsonnetFromBaseDir(logE *logrus.Entry, baseDir string) ([]string, error) {
	filePaths := []string{}
	if err := fs.WalkDir(afero.NewIOFS(c.fs), baseDir, func(p string, dirEntry fs.DirEntry, e error) error {
		if e != nil {
			logE.WithError(e).Warn("error occurred while searching files")
			return nil
		}
		if dirEntry.Type().IsDir() {
			return nil
		}
		if !strings.HasSuffix(p, ".jsonnet") {
			return nil
		}
		filePaths = append(filePaths, p)
		return nil
	}); err != nil {
		return nil, fmt.Errorf("walks the file tree of the unarchived package: %w", err)
	}
	return filePaths, nil
}

func (c *Controller) readJsonnets(filePaths []*LintFile, modules map[string]string) (map[string]ast.Node, error) {
	jsonnetAsts := make(map[string]ast.Node, len(filePaths))
	for _, filePath := range filePaths {
		ja, err := c.readJsonnet(filePath.Path)
		if err != nil {
			return nil, logerr.WithFields(err, logrus.Fields{ //nolint:wrapcheck
				"file_path": filePath,
			})
		}
		jsonnetAsts[filePath.Path] = ja
	}
	return jsonnetAsts, nil
}

func newVM(data *Data, libs map[string]string) *jsonnet.VM {
	vm := jsonnet.MakeVM()
	vm.ExtCode("input", string(data.JSON))
	vm.ExtVar("file_path", data.FilePath)
	vm.ExtVar("file_type", data.FileType)
	vm.ExtVar("file_text", data.Text)
	setNativeFunctions(vm)

	if len(libs) != 0 {
		m := make(map[string]jsonnet.Contents, len(libs))
		for k, v := range libs {
			m[k] = jsonnet.MakeContents(v)
		}
		mi := &jsonnet.MemoryImporter{
			Data: m,
		}
		vm.Importer(mi)
	}

	return vm
}

func (c *Controller) readJsonnet(filePath string) (ast.Node, error) {
	b, err := afero.ReadFile(c.fs, filePath)
	if err != nil {
		return nil, fmt.Errorf("read a jsonnet file: %w", err)
	}
	ja, err := jsonnet.SnippetToAST(filePath, string(b))
	if err != nil {
		return nil, fmt.Errorf("parse a jsonnet file: %w", err)
	}
	return ja, nil
}

func (c *Controller) evaluate(data *Data, jsonnetAsts map[string]ast.Node, libs map[string]string) map[string]*JsonnetEvaluateResult {
	vm := newVM(data, libs)

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
