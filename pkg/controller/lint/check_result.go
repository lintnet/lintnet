package lint

func isFailed(results map[string]*FileResult) bool {
	// data file -> result
	for _, result := range results {
		if result.Error != "" {
			return true
		}
		// lint file -> result
		for _, r := range result.Results {
			if r.Error != "" {
				return true
			}
			if isJsonnetResultFailed(r.RawResult) {
				return true
			}
		}
	}
	return false
}

func isJsonnetResultFailed(result *JsonnetResult) bool {
	if result.Failed {
		return true
	}
	if result.Error != "" {
		return true
	}
	for _, l := range result.Locations {
		if l.S != "" || l.Raw != nil {
			return true
		}
	}
	for _, e := range result.Errors {
		if e.Error != "" {
			return true
		}
	}
	for _, r := range result.SubRules {
		if isJsonnetResultFailed(r) {
			return true
		}
	}
	return false
}
