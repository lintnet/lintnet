package lint

func isFailed(results []*FlatError, errLevel ErrorLevel) (bool, error) {
	for _, result := range results {
		e := result.Level
		if e == "" {
			e = "error"
		}
		feErrLevel, err := newErrorLevel(e)
		if err != nil {
			return false, err
		}
		if feErrLevel >= errLevel {
			return true, nil
		}
	}
	return false, nil
}

func (r *JsonnetResult) isFailed() bool {
	return r.Failed
}
