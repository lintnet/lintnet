package lint

import (
	"encoding/json"
	"fmt"

	"github.com/lintnet/lintnet/pkg/domain"
	"github.com/lintnet/lintnet/pkg/jsonnet"
)

type LintFileEvaluator struct { //nolint:revive
	importer *jsonnet.Importer
}

func (le *LintFileEvaluator) Evaluate(tla *domain.TopLevelArgment, lintFile jsonnet.Node) (string, error) {
	if tla.Config == nil {
		tla.Config = map[string]any{}
	}
	tlaB, err := json.Marshal(tla)
	if err != nil {
		return "", fmt.Errorf("marshal a top level argument as JSON: %w", err)
	}
	vm := jsonnet.NewVM(string(tlaB), le.importer)
	result, err := vm.Evaluate(lintFile)
	if err != nil {
		return "", fmt.Errorf("evaluate a lint file as Jsonnet: %w", err)
	}
	return result, nil
}

func (le *LintFileEvaluator) Evaluates(tla *domain.TopLevelArgment, lintFiles []*domain.Node) []*domain.Result {
	results := make([]*domain.Result, len(lintFiles))
	for i, lintFile := range lintFiles {
		tla := &domain.TopLevelArgment{
			Data:         tla.Data,
			CombinedData: tla.CombinedData,
			Config:       lintFile.Config,
		}
		s, err := le.Evaluate(tla, lintFile.Node)
		if err != nil {
			results[i] = &domain.Result{
				LintFile: lintFile.Key,
				Error:    err.Error(),
			}
			continue
		}
		rs, a, err := parseResult([]byte(s))
		results[i] = &domain.Result{
			LintFile:  lintFile.Key,
			RawResult: rs,
			RawOutput: s,
			Interface: a,
		}
		if err != nil {
			results[i].Error = err.Error()
		}
	}
	return results
}

func parseResult(result []byte) ([]*domain.JsonnetResult, any, error) {
	var rs any
	if err := json.Unmarshal(result, &rs); err != nil {
		return nil, nil, fmt.Errorf("unmarshal the result as JSON: %w", err)
	}

	out := []*domain.JsonnetResult{}
	if err := json.Unmarshal(result, &out); err != nil {
		return nil, rs, fmt.Errorf("unmarshal the result as JSON: %w", err)
	}
	return out, rs, nil
}
