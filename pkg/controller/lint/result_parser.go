package lint

func checkFailed(results map[string]*FileResult) bool {
	for _, result := range results {
		if result.Error != "" {
			return true
		}
		for _, r := range result.Results {
			if r.Error != "" {
				return true
			}
			for _, rule := range r.Output.Rules {
				if len(rule.Errors) != 0 {
					return true
				}
			}
		}
	}
	return false
}
