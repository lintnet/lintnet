package testcmd

import (
	_ "embed"
	"io"

	"github.com/lintnet/lintnet/pkg/config"
	"github.com/lintnet/lintnet/pkg/config/reader"
	"github.com/lintnet/lintnet/pkg/domain"
	"github.com/lintnet/lintnet/pkg/encoding"
	"github.com/lintnet/lintnet/pkg/filefind"
	"github.com/lintnet/lintnet/pkg/jsonnet"
	"github.com/lintnet/lintnet/pkg/lint"
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

type TestController struct {
	fs             afero.Fs
	stdout         io.Writer
	importer       *jsonnet.ModuleImporter
	param          *ParamController
	dataFileParser lint.DataFileParser
	fileFinder     FileFinder
	configReader   *reader.Reader
}

type FileFinder interface {
	Find(logE *logrus.Entry, cfg *config.Config, rootDir, cfgDir string) ([]*domain.Target, error)
}

type ParamController struct {
	Version string
}

func NewTestController(param *ParamController, fs afero.Fs, stdout io.Writer, importer *jsonnet.ModuleImporter) *TestController {
	dp := encoding.NewDataFileParser(fs)
	return &TestController{
		param:          param,
		fs:             fs,
		stdout:         stdout,
		importer:       importer,
		dataFileParser: dp,
		fileFinder:     filefind.NewFileFinder(fs),
		configReader:   reader.New(fs, importer),
	}
}
