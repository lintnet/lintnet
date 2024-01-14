package lint

type Linter struct {
	dataFileParser    *DataFileParser
	lintFileParser    *LintFileParser
	lintFileEvaluator *LintFileEvaluator
}

func (l *Linter) Lint(targets []*Target) ([]*Result, error) {
	results := make([]*Result, 0, len(targets))
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

func (l *Linter) lintTarget(target *Target) ([]*Result, error) {
	lintFiles, err := l.lintFileParser.Parses(target.LintFiles)
	if err != nil {
		return nil, err
	}

	combineFiles := []*Node{}
	nonCombineFiles := []*Node{}
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

func (l *Linter) lintCombineFiles(target *Target, combineFiles []*Node) ([]*Result, error) {
	rs, err := l.lint(&DataSet{
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

func (l *Linter) lintNonCombineFiles(target *Target, nonCombineFiles []*Node) []*Result {
	results := make([]*Result, 0, len(target.DataFiles))
	for _, dataFile := range target.DataFiles {
		results = append(results, l.lintNonCombineFile(nonCombineFiles, dataFile)...)
	}
	return results
}

func (l *Linter) lintNonCombineFile(nonCombineFiles []*Node, dataFile *Path) []*Result {
	rs, err := l.lint(&DataSet{
		File: dataFile,
	}, nonCombineFiles)
	if err != nil {
		return []*Result{
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

func (l *Linter) getTLA(dataSet *DataSet) (*TopLevelArgment, error) {
	if dataSet.File != nil {
		return l.dataFileParser.Parse(dataSet.File)
	}
	if len(dataSet.Files) > 0 {
		combinedData := make([]*Data, len(dataSet.Files))
		for i, dataFile := range dataSet.Files {
			data, err := l.dataFileParser.Parse(dataFile)
			if err != nil {
				return nil, err
			}
			combinedData[i] = data.Data
		}
		return &TopLevelArgment{
			CombinedData: combinedData,
		}, nil
	}
	return &TopLevelArgment{}, nil
}

func (l *Linter) lint(dataSet *DataSet, lintFiles []*Node) ([]*Result, error) {
	tla, err := l.getTLA(dataSet)
	if err != nil {
		return nil, err
	}
	return l.lintFileEvaluator.Evaluates(tla, lintFiles), nil
}
