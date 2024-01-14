package lint

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/google/go-cmp/cmp"
	"github.com/lintnet/lintnet/pkg/config"
	"github.com/lintnet/lintnet/pkg/domain"
	"github.com/lintnet/lintnet/pkg/jsonnet"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

//go:embed test_diff.txt
var testResultTemplateByte []byte

type TestData struct {
	Name     string                  `json:"name,omitempty"`
	DataFile string                  `json:"data_file,omitempty"`
	Param    *domain.TopLevelArgment `json:"param,omitempty"`
	Result   []any                   `json:"result,omitempty"`
}

type TestPair struct {
	LintFilePath string
	TestFilePath string
}

type FailedResult struct {
	Name         string `json:"name,omitempty"`
	LintFilePath string `json:"lint_file_path,omitempty"`
	TestFilePath string `json:"test_file_path,omitempty"`
	Param        any    `json:"param,omitempty"`
	Wanted       any    `json:"wanted,omitempty"`
	Got          any    `json:"got,omitempty"`
	Diff         string `json:"diff,omitempty"`
	Error        string `json:"error,omitempty"`
}

type TestResult struct {
	Name     string `json:"name,omitempty"`
	Links    any    `json:"links,omitempty"`
	Message  string `json:"message,omitempty"`
	Level    string `json:"level,omitempty"`
	Location any    `json:"location,omitempty"`
	Custom   any    `json:"custom,omitempty"`
	Excluded bool   `json:"excluded,omitempty"`
}

func (tr *TestResult) Any() any {
	m := map[string]any{}
	if tr.Name != "" {
		m["name"] = tr.Name
	}
	if tr.Links != nil {
		m["links"] = tr.Links
	}
	if tr.Message != "" {
		m["message"] = tr.Message
	}
	if tr.Level != "" {
		m["level"] = tr.Level
	}
	if tr.Location != nil {
		m["location"] = tr.Location
	}
	if tr.Custom != nil {
		m["custom"] = tr.Custom
	}
	return m
}

func (c *Controller) Test(_ context.Context, logE *logrus.Entry, param *ParamLint) error {
	rawCfg := &config.RawConfig{}
	if err := c.configReader.Read(param.ConfigFilePath, rawCfg); err != nil {
		return err
	}
	cfg, err := rawCfg.Parse()
	if err != nil {
		return fmt.Errorf("parse a configuration file: %w", err)
	}

	testResultTemplate, err := template.New("_").Parse(string(testResultTemplateByte))
	if err != nil {
		return fmt.Errorf("parse the template of test result: %w", err)
	}

	cfgDir := filepath.Dir(rawCfg.FilePath)

	targets, err := c.fileFinder.Find(logE, cfg, param.RootDir, cfgDir)
	if err != nil {
		return fmt.Errorf("find files: %w", err)
	}

	pairs := c.filterTargetsWithTest(logE, targets)
	failedResults := make([]*FailedResult, 0, len(pairs))
	for _, pair := range pairs {
		if results := c.tests(pair); len(results) > 0 {
			failedResults = append(failedResults, results...)
		}
	}
	if len(failedResults) == 0 {
		return nil
	}
	if err := testResultTemplate.Execute(c.stdout, failedResults); err != nil {
		return fmt.Errorf("render the result: %w", err)
	}
	return nil
}

func (c *Controller) test(pair *TestPair, td *TestData) *FailedResult { //nolint:cyclop
	if td.DataFile != "" {
		p := &domain.Path{
			Raw: td.DataFile,
			Abs: filepath.Join(filepath.Dir(pair.TestFilePath), td.DataFile),
		}
		data, err := c.dataFileParser.Parse(p)
		if err != nil {
			return &FailedResult{
				Error: fmt.Errorf("read a data file: %w", err).Error(),
			}
		}
		if td.Param != nil && td.Param.Data != nil && td.Param.Data.FilePath != "" {
			data.Data.FilePath = td.Param.Data.FilePath
		}
		if td.Param != nil {
			data.Config = td.Param.Config
		}
		td.Param = data
	}
	if td.Param.Config == nil {
		td.Param.Config = map[string]any{}
	}
	tlaB, err := json.Marshal(td.Param)
	if err != nil {
		return &FailedResult{
			Error: fmt.Errorf("marshal param as JSON: %w", err).Error(),
		}
	}
	var results []*TestResult
	if err := jsonnet.Read(c.fs, pair.LintFilePath, string(tlaB), c.importer, &results); err != nil {
		return &FailedResult{
			Error: fmt.Errorf("read a lint file: %w", err).Error(),
		}
	}
	rs := make([]any, 0, len(results))
	for _, result := range results {
		if result.Excluded {
			continue
		}
		rs = append(rs, result.Any())
	}
	if len(rs) == 0 && len(td.Result) == 0 {
		return nil
	}
	if diff := cmp.Diff(td.Result, rs); diff != "" {
		return &FailedResult{
			Wanted: td.Result,
			Got:    rs,
			Diff:   diff,
		}
	}
	return nil
}

func (c *Controller) tests(pair *TestPair) []*FailedResult {
	testData := []*TestData{}
	if err := jsonnet.Read(c.fs, pair.TestFilePath, "{}", c.importer, &testData); err != nil {
		return []*FailedResult{
			{
				LintFilePath: pair.LintFilePath,
				TestFilePath: pair.TestFilePath,
				Error:        fmt.Errorf("read a test file: %w", err).Error(),
			},
		}
	}
	results := make([]*FailedResult, 0, len(testData))
	for _, td := range testData {
		if result := c.test(pair, td); result != nil {
			result.Name = td.Name
			result.LintFilePath = pair.LintFilePath
			result.TestFilePath = pair.TestFilePath
			result.Param = td.Param
			results = append(results, result)
		}
	}
	return results
}

func (c *Controller) filterTargetsWithTest(logE *logrus.Entry, targets []*domain.Target) []*TestPair {
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
