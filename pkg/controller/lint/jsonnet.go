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

func (c *Controller) findFilesbyGlob(target *config.Target) ([]*LintFile, error) {
	filePaths := make([]*LintFile, 0, len(target.LintFiles.Paths))
	for _, p := range target.LintFiles.Paths {
		matches, err := fs.Glob(afero.NewIOFS(c.fs), p.Path)
		if err != nil {
			return nil, fmt.Errorf("search files by glob: %w", err)
		}
		for _, match := range matches {
			filePaths = append(filePaths, &LintFile{
				Path: match,
			})
		}
	}
	return filePaths, nil
}

func (c *Controller) findFilesByRegexp(logE *logrus.Entry, target *config.Target) ([]*LintFile, error) {
	patterns := make([]*regexp.Regexp, len(target.LintFiles.Paths))
	for i, p := range target.LintFiles.Paths {
		p, err := regexp.Compile(p.Path)
		if err != nil {
			return nil, fmt.Errorf("compile a regular expression to search files: %w", err)
		}
		patterns[i] = p
	}
	filePaths := make([]*LintFile, 0, len(target.LintFiles.Paths))
	if err := fs.WalkDir(afero.NewIOFS(c.fs), "", func(p string, dirEntry fs.DirEntry, e error) error {
		if e != nil {
			logE.WithError(e).Warn("error occurred while searching files")
			return nil
		}
		for _, pattern := range patterns {
			if pattern.MatchString(p) {
				filePaths = append(filePaths, &LintFile{
					Path: p,
				})
			}
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("search files: %w", err)
	}
	return filePaths, nil
}

func (c *Controller) findFilesFromTarget(logE *logrus.Entry, target *config.Target) ([]*LintFile, error) {
	switch target.LintFiles.SearchType {
	case "equal":
		filePaths := make([]*LintFile, len(target.LintFiles.Paths))
		for i, p := range target.LintFiles.Paths {
			filePaths[i] = &LintFile{
				Path: p.Path,
			}
		}
		return filePaths, nil
	case "glob":
		return c.findFilesbyGlob(target)
	case "regexp":
		return c.findFilesByRegexp(logE, target)
	default:
		return nil, logerr.WithFields(errors.New("search_type is invalid"), logrus.Fields{ //nolint:wrapcheck
			"search_type": target.LintFiles.SearchType,
		})
	}
}

func (c *Controller) findJsonnet(logE *logrus.Entry, cfg *config.Config, ruleBaseDir string) ([]*LintFile, error) {
	if ruleBaseDir != "" {
		return c.findJsonnetFromBaseDir(logE, ruleBaseDir)
	}
	if len(cfg.Targets) == 0 {
		return c.findJsonnetFromBaseDir(logE, "lintnet")
	}
	filePaths := make([]*LintFile, 0, len(cfg.Targets))
	for _, target := range cfg.Targets {
		fps, err := c.findFilesFromTarget(logE, target)
		if err != nil {
			return nil, err
		}
		filePaths = append(filePaths, fps...)
	}
	return filePaths, nil
}

func (c *Controller) findJsonnetFromBaseDir(logE *logrus.Entry, baseDir string) ([]*LintFile, error) {
	filePaths := []*LintFile{}
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
		filePaths = append(filePaths, &LintFile{
			Path: p,
		})
		return nil
	}); err != nil {
		return nil, fmt.Errorf("walks the file tree of the unarchived package: %w", err)
	}
	return filePaths, nil
}

func (c *Controller) readJsonnets(filePaths []*LintFile) (map[string]ast.Node, error) {
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

func newVM(data *Data) *jsonnet.VM {
	vm := jsonnet.MakeVM()
	vm.ExtCode("input", string(data.JSON))
	vm.ExtVar("file_path", data.FilePath)
	vm.ExtVar("file_type", data.FileType)
	vm.ExtVar("file_text", data.Text)
	setNativeFunctions(vm)
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

func (c *Controller) evaluate(data *Data, jsonnetAsts map[string]ast.Node) map[string]*JsonnetEvaluateResult {
	vm := newVM(data)

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
