package lint

func (c *Controller) getResults(targets []*Target) ([]*Result, error) {
	results := make([]*Result, 0, len(targets))
	for _, target := range targets {
		rs, err := c.lintTarget(target)
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

func (c *Controller) lintTarget(target *Target) ([]*Result, error) {
	lintFiles, err := c.lintFileParser.Parses(target.LintFiles)
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

	results := c.lintNonCombineFiles(target, nonCombineFiles)

	if len(combineFiles) > 0 {
		rs, err := c.lintCombineFiles(target, combineFiles)
		if err != nil {
			return nil, err
		}
		return append(results, rs...), nil
	}
	return results, nil
}

func (c *Controller) lintCombineFiles(target *Target, combineFiles []*Node) ([]*Result, error) {
	rs, err := c.lint(&DataSet{
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

func (c *Controller) lintNonCombineFiles(target *Target, nonCombineFiles []*Node) []*Result {
	results := make([]*Result, 0, len(target.DataFiles))
	for _, dataFile := range target.DataFiles {
		results = append(results, c.lintNonCombineFile(nonCombineFiles, dataFile)...)
	}
	return results
}

func (c *Controller) lintNonCombineFile(nonCombineFiles []*Node, dataFile *Path) []*Result {
	rs, err := c.lint(&DataSet{
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

func (c *Controller) getTLA(dataSet *DataSet) (*TopLevelArgment, error) {
	if dataSet.File != nil {
		return c.parseDataFile(dataSet.File)
	}
	if len(dataSet.Files) > 0 {
		combinedData := make([]*Data, len(dataSet.Files))
		for i, dataFile := range dataSet.Files {
			data, err := c.parseDataFile(dataFile)
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

func (c *Controller) lint(dataSet *DataSet, lintFiles []*Node) ([]*Result, error) {
	tla, err := c.getTLA(dataSet)
	if err != nil {
		return nil, err
	}
	return c.lintFileEvaluator.Evaluates(tla, lintFiles), nil
}
