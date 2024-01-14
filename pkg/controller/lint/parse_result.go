package lint

type (
	// return of vm.Evaluate()
	JsonnetEvaluateResult struct {
		// Key    string
		Result string
		Error  string
	}
)
