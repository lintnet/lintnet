package testcmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/google/go-cmp/cmp"
	"github.com/lintnet/lintnet/pkg/config"
	"github.com/lintnet/lintnet/pkg/domain"
	"github.com/lintnet/lintnet/pkg/filefilter"
	"github.com/lintnet/lintnet/pkg/jsonnet"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
)

type ParamTest struct {
	RootDir        string
	ConfigFilePath string
	TargetID       string
	PWD            string
	FilePaths      []string
}

func (p *ParamTest) FilterParam() *filefilter.Param {
	return &filefilter.Param{
		TargetID: p.TargetID,
		PWD:      p.PWD,
	}
}

func (c *Controller) Test(_ context.Context, logE *logrus.Entry, param *ParamTest) error {
	pairs, err := c.listPairs(logE, param)
	if err != nil {
		return err
	}

	testResultTemplate, err := template.New("_").Parse(string(testResultTemplateByte))
	if err != nil {
		return fmt.Errorf("parse the template of test result: %w", err)
	}

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
	return errors.New("test failed")
}

func (c *Controller) listPairs(logE *logrus.Entry, param *ParamTest) ([]*TestPair, error) {
	if len(param.FilePaths) != 0 {
		return c.listPairsWithFilePaths(param.FilePaths)
	}

	rawCfg := &config.RawConfig{}
	if err := c.configReader.Read(param.ConfigFilePath, rawCfg); err != nil {
		if param.ConfigFilePath == "" && errors.Is(err, fs.ErrNotExist) {
			return c.listPairsWithFilePaths([]string{"."})
		}
		return nil, fmt.Errorf("read a configuration file: %w", err)
	}

	if param.TargetID != "" {
		target, err := rawCfg.GetTarget(param.TargetID)
		if err != nil {
			return nil, fmt.Errorf("get a target from configuration file by target id: %w", err)
		}
		rawCfg.Targets = []*config.RawTarget{target}
	}

	cfg, err := rawCfg.Parse()
	if err != nil {
		return nil, fmt.Errorf("parse a configuration file: %w", err)
	}

	cfgDir := filepath.Dir(rawCfg.FilePath)

	lintFiles, err := c.fileFinder.FindLintFiles(logE, cfg, cfgDir)
	if err != nil {
		return nil, fmt.Errorf("find files: %w", err)
	}

	return c.filterLintFilesWithTest(logE, lintFiles), nil
}

func getTestFilePath(lintFilePath string) string {
	return lintFilePath[:len(lintFilePath)-len(".jsonnet")] + "_test.jsonnet"
}

func getLintFilePath(testFilePath string) string {
	return testFilePath[:len(testFilePath)-len("_test.jsonnet")] + ".jsonnet"
}

func (c *Controller) listPairsWithFilePath(filePath string) ([]*TestPair, error) { //nolint:cyclop
	switch {
	case strings.HasSuffix(filePath, "_test.jsonnet"):
		lintFile := getLintFilePath(filePath)
		return []*TestPair{
			{
				LintFilePath: lintFile,
				TestFilePath: filePath,
			},
		}, nil
	case strings.HasSuffix(filePath, ".jsonnet"):
		tp := getTestFilePath(filePath)
		if f, err := afero.Exists(c.fs, tp); err != nil {
			return nil, fmt.Errorf("check if a file exists: %w", err)
		} else if !f {
			return nil, nil
		}
		return []*TestPair{
			{
				LintFilePath: filePath,
				TestFilePath: tp,
			},
		}, nil
	default:
		if b, err := afero.IsDir(c.fs, filePath); err != nil {
			return nil, fmt.Errorf("check if a path is a directory: %w", err)
		} else if !b {
			return nil, nil
		}
		pairs := []*TestPair{}
		if err := doublestar.GlobWalk(afero.NewIOFS(c.fs), filePath+"/**/*_test.jsonnet", func(testFile string, _ fs.DirEntry) error {
			lintFile := getLintFilePath(testFile)
			a, err := afero.Exists(c.fs, lintFile)
			if err != nil {
				return fmt.Errorf("check if a lint file exists: %w", err)
			}
			if !a {
				return nil
			}
			pairs = append(pairs, &TestPair{
				LintFilePath: lintFile,
				TestFilePath: testFile,
			})
			return nil
		}, doublestar.WithNoFollow()); err != nil {
			return nil, fmt.Errorf("search files: %w", err)
		}
		return pairs, nil
	}
}

func (c *Controller) listPairsWithFilePaths(filePaths []string) ([]*TestPair, error) {
	pairs := make([]*TestPair, 0, len(filePaths))
	for _, p := range filePaths {
		ps, err := c.listPairsWithFilePath(p)
		if err != nil {
			return nil, err
		}
		pairs = append(pairs, ps...)
	}
	return pairs, nil
}

func (c *Controller) test(pair *TestPair, td *TestData) *FailedResult { //nolint:cyclop
	if td.DataFile != "" {
		if err := c.readDatafile(pair, td); err != nil {
			return &FailedResult{
				Error: err.Error(),
			}
		}
	}

	if len(td.DataFiles) != 0 {
		if err := c.readDatafiles(pair, td); err != nil {
			return &FailedResult{
				Error: err.Error(),
			}
		}
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

func (c *Controller) readDatafile(pair *TestPair, td *TestData) error {
	p := &domain.Path{
		Raw: td.DataFile,
		Abs: filepath.Join(filepath.Dir(pair.TestFilePath), td.DataFile),
	}
	data, err := c.dataFileParser.Parse(p)
	if err != nil {
		return fmt.Errorf("read a data file: %w", err)
	}
	if td.Param != nil && td.Param.Data != nil && td.Param.Data.FilePath != "" {
		data.Data.FilePath = td.Param.Data.FilePath
	}
	if td.FakeDataFile != "" {
		data.Data.FilePath = td.FakeDataFile
	}
	if td.Param != nil {
		data.Config = td.Param.Config
	}
	td.Param = data
	return nil
}

func (c *Controller) readDatafiles(pair *TestPair, td *TestData) error {
	combinedData := make([]*domain.Data, len(td.DataFiles))
	for i, dataFile := range td.DataFiles {
		p := &domain.Path{
			Raw: dataFile.Path,
			Abs: filepath.Join(filepath.Dir(pair.TestFilePath), dataFile.Path),
		}
		data, err := c.dataFileParser.Parse(p)
		if err != nil {
			return fmt.Errorf("read a data file: %w", logerr.WithFields(err, logrus.Fields{
				"data_file": dataFile.Path,
			}))
		}
		if dataFile.FakePath != "" {
			data.Data.FilePath = dataFile.FakePath
		}
		if td.Param != nil {
			data.Config = td.Param.Config
		}
		combinedData[i] = data.Data
	}
	if td.Param == nil {
		td.Param = &domain.TopLevelArgument{}
	}
	td.Param.CombinedData = combinedData
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

func (c *Controller) filterLintFilesWithTest(logE *logrus.Entry, lintFiles []*config.LintFile) []*TestPair {
	pairs := []*TestPair{}
	for _, lintFile := range lintFiles {
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
	return pairs
}
