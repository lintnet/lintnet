package lint

func isFailed(results map[string]*FileResult) bool {
	// data file -> result
	for _, result := range results {
		if result.isFailed() {
			return true
		}
	}
	return false
}

func (r *JsonnetResult) isFailed() bool {
	if r.Failed {
		return true
	}
	if r.Message != "" {
		return true
	}
	return false
}
