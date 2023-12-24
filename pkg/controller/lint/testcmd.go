package lint

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/lintnet/lintnet/pkg/config"
	"github.com/lintnet/lintnet/pkg/jsonnet"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

func (c *Controller) Test(_ context.Context, logE *logrus.Entry, param *ParamLint) error { //nolint:funlen,cyclop
	cfg := &config.Config{}
	if err := c.findAndReadConfig(param.ConfigFilePath, cfg); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}

	targets, err := c.findFiles(logE, cfg, nil, param.RuleBaseDir, param.FilePaths, param.RootDir)
	if err != nil {
		return err
	}

	pairs := c.filterTargetsWithTest(logE, targets)
	failedResults := make([]*FailedResult, 0, len(pairs))
	for _, pair := range pairs {
		testData := []*TestData{}
		if err := jsonnet.Read(c.fs, pair.TestFilePath, "{}", c.importer, &testData); err != nil {
			failedResults = append(failedResults, &FailedResult{
				LintFilePath: pair.LintFilePath,
				TestFilePath: pair.TestFilePath,
				Error:        fmt.Errorf("read a test file: %w", err).Error(),
			})
			continue
		}
		for _, td := range testData {
			var result interface{}
			tlaB, err := json.Marshal(td.Param)
			if err != nil {
				failedResults = append(failedResults, &FailedResult{
					LintFilePath: pair.LintFilePath,
					TestFilePath: pair.TestFilePath,
					Param:        td.Param,
					Error:        fmt.Errorf("marshal param as JSON: %w", err).Error(),
				})
				continue
			}
			if err := jsonnet.Read(c.fs, pair.LintFilePath, string(tlaB), c.importer, &result); err != nil {
				failedResults = append(failedResults, &FailedResult{
					LintFilePath: pair.LintFilePath,
					TestFilePath: pair.TestFilePath,
					Param:        td.Param,
					Error:        fmt.Errorf("read a lint file: %w", err).Error(),
				})
				continue
			}
			if diff := cmp.Diff(td.Result, result); diff != "" {
				failedResults = append(failedResults, &FailedResult{
					LintFilePath: pair.LintFilePath,
					TestFilePath: pair.TestFilePath,
					Wanted:       td.Result,
					Param:        td.Param,
					Got:          result,
					Diff:         diff,
				})
			}
		}
	}
	if len(failedResults) == 0 {
		return nil
	}
	for _, failedResult := range failedResults {
		fmt.Fprintf(c.stdout, `name: %s
lint file: %s
test file: %s
param: %v
error: %s
wanted: %v
got: %v
diff: %s
`, failedResult.Name, failedResult.LintFilePath, failedResult.TestFilePath, failedResult.Param, failedResult.Error, failedResult.Wanted, failedResult.Got, failedResult.Diff)
	}

	return nil
}

type TestData struct {
	Name   string      `json:"name,omitempty"`
	Param  interface{} `json:"param,omitempty"`
	Result interface{} `json:"result,omitempty"`
}

type TestPair struct {
	LintFilePath string
	TestFilePath string
}

type FailedResult struct {
	Name         string      `json:"name,omitempty"`
	LintFilePath string      `json:"lint_file_path,omitempty"`
	TestFilePath string      `json:"test_file_path,omitempty"`
	Param        interface{} `json:"param,omitempty"`
	Wanted       interface{} `json:"wanted,omitempty"`
	Got          interface{} `json:"got,omitempty"`
	Diff         string      `json:"diff,omitempty"`
	Error        string      `json:"error,omitempty"`
}

func (c *Controller) filterTargetsWithTest(logE *logrus.Entry, targets []*Target) []*TestPair {
	pairs := []*TestPair{}
	for _, target := range targets {
		for _, lintFile := range target.LintFiles {
			if lintFile.Path == "" {
				continue
			}
			baseName := filepath.Base(lintFile.Path)
			ext := filepath.Ext(baseName)
			testFileName := fmt.Sprintf("%s_test%s", strings.TrimSuffix(baseName, filepath.Ext(baseName)), ext)
			testFilePath := filepath.Join(filepath.Dir(lintFile.Path), testFileName)
			f, err := afero.Exists(c.fs, testFilePath)
			if err != nil {
				logE.WithError(err).Warn("check if a test file exists")
				continue
			}
			if !f {
				continue
			}
			pairs = append(pairs, &TestPair{
				LintFilePath: lintFile.Path,
				TestFilePath: testFilePath,
			})
		}
	}
	return pairs
}
