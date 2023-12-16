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
	if r.Error != "" {
		return true
	}
	for _, l := range r.Locations {
		if l.S != "" || l.Raw != nil {
			return true
		}
	}
	for _, e := range r.Errors {
		if e.Error != "" {
			return true
		}
	}
	for _, sub := range r.SubRules {
		if sub.isFailed() {
			return true
		}
	}
	return false
}
