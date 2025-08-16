package lint

import (
	"fmt"

	"github.com/lintnet/lintnet/pkg/config"
	"github.com/lintnet/lintnet/pkg/domain"
	"github.com/lintnet/lintnet/pkg/filefind"
	"github.com/lintnet/lintnet/pkg/jsonnet"
	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
)

type Linter struct {
	dataFileParser    DataFileParser
	lintFileParser    LintFileParser
	lintFileEvaluator LintFileEvaluator
}

func NewLinter(dataFileParser DataFileParser, lintFileParser LintFileParser, lintFileEvaluator LintFileEvaluator) *Linter {
	return &Linter{
		dataFileParser:    dataFileParser,
		lintFileParser:    lintFileParser,
		lintFileEvaluator: lintFileEvaluator,
	}
}

type DataFileParser interface {
	Parse(filePath *domain.Path) (*domain.TopLevelArgument, error)
}

type LintFileParser interface { //nolint:revive
	Parse(lintFile *config.LintFile) (*domain.Node, error)
	Parses(lintFiles []*config.LintFile) ([]*domain.Node, error)
}

type LintFileEvaluator interface { //nolint:revive
	Evaluate(tla *domain.TopLevelArgument, lintFile jsonnet.Node) (string, error)
	Evaluates(tla *domain.TopLevelArgument, lintFiles []*domain.Node) []*domain.Result
}

func (l *Linter) Lint(targets []*filefind.Target) ([]*domain.Result, error) {
	results := make([]*domain.Result, 0, len(targets))
	for _, target := range targets {
		rs, err := l.lintTarget(target)
		if err != nil {
			return nil, err
		}
		for _, r := range rs {
			r.TargetID = target.ID
		}
		results = append(results, rs...)
	}
	return results, nil
}

func (l *Linter) lintTarget(target *filefind.Target) ([]*domain.Result, error) {
	lintFiles, err := l.lintFileParser.Parses(target.LintFiles)
	if err != nil {
		return nil, fmt.Errorf("parse lint files: %w", err)
	}

	combineFiles := []*domain.Node{}
	nonCombineFiles := []*domain.Node{}
	for _, lintFile := range lintFiles {
		if lintFile.Combine {
			combineFiles = append(combineFiles, lintFile)
			continue
		}
		nonCombineFiles = append(nonCombineFiles, lintFile)
	}

	results := l.lintNonCombineFiles(target, nonCombineFiles)

	if len(combineFiles) > 0 {
		rs, err := l.lintCombineFiles(target, combineFiles)
		if err != nil {
			return nil, err
		}
		return append(results, rs...), nil
	}
	return results, nil
}

func (l *Linter) lintCombineFiles(target *filefind.Target, combineFiles []*domain.Node) ([]*domain.Result, error) {
	rs, err := l.lint(&domain.DataSet{
		Files: target.DataFiles,
	}, combineFiles)
	if err != nil {
		return nil, err
	}
	for _, r := range rs {
		arr := make([]string, len(target.DataFiles))
		for i, dataFile := range target.DataFiles {
			arr[i] = dataFile.Raw
		}
		r.DataFiles = arr
	}
	return rs, nil
}

func (l *Linter) lintNonCombineFiles(target *filefind.Target, nonCombineFiles []*domain.Node) []*domain.Result {
	results := make([]*domain.Result, 0, len(target.DataFiles))
	for _, dataFile := range target.DataFiles {
		results = append(results, l.lintNonCombineFile(nonCombineFiles, dataFile)...)
	}
	return results
}

func (l *Linter) lintNonCombineFile(nonCombineFiles []*domain.Node, dataFile *domain.Path) []*domain.Result {
	rs, err := l.lint(&domain.DataSet{
		File: dataFile,
	}, nonCombineFiles)
	if err != nil {
		return []*domain.Result{
			{
				DataFile: dataFile.Raw,
				Error:    err.Error(),
			},
		}
	}
	for _, r := range rs {
		r.DataFile = dataFile.Raw
	}
	return rs
}

func (l *Linter) getTLA(dataSet *domain.DataSet) (*domain.TopLevelArgument, error) {
	if dataSet.File != nil {
		tla, err := l.dataFileParser.Parse(dataSet.File)
		if err != nil {
			return nil, fmt.Errorf("parse a data file: %w", err)
		}
		return tla, nil
	}
	if len(dataSet.Files) > 0 {
		combinedData := make([]*domain.Data, len(dataSet.Files))
		for i, dataFile := range dataSet.Files {
			data, err := l.dataFileParser.Parse(dataFile)
			if err != nil {
				return nil, fmt.Errorf("parse a data file: %w", logerr.WithFields(err, logrus.Fields{
					"data_file": dataFile.Raw,
				}))
			}
			combinedData[i] = data.Data
		}
		return &domain.TopLevelArgument{
			CombinedData: combinedData,
		}, nil
	}
	return &domain.TopLevelArgument{}, nil
}

func (l *Linter) lint(dataSet *domain.DataSet, lintFiles []*domain.Node) ([]*domain.Result, error) {
	tla, err := l.getTLA(dataSet)
	if err != nil {
		return nil, err
	}
	return l.lintFileEvaluator.Evaluates(tla, lintFiles), nil
}
