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
	return r.Failed
}
