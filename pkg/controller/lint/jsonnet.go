package lint

import (
	"encoding/json"
	"fmt"

	"github.com/lintnet/lintnet/pkg/jsonnet"
	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
)

func (c *Controller) readJsonnets(filePaths []*LintFile) ([]*Node, error) {
	jsonnetAsts := make([]*Node, 0, len(filePaths))
	for _, filePath := range filePaths {
		ja, err := jsonnet.ReadToNode(c.fs, filePath.Path)
		if err != nil {
			return nil, logerr.WithFields(err, logrus.Fields{ //nolint:wrapcheck
				"file_path": filePath,
			})
		}
		if filePath.ModulePath != "" {
			jsonnetAsts = append(jsonnetAsts, &Node{
				Node:   ja,
				Key:    filePath.ModulePath,
				Custom: filePath.Param,
			})
			continue
		}
		jsonnetAsts = append(jsonnetAsts, &Node{
			Node:   ja,
			Key:    filePath.Path,
			Custom: filePath.Param,
		})
	}
	return jsonnetAsts, nil
}

type Node struct {
	Node   jsonnet.Node
	Custom interface{}
	Key    string
}

func (c *Controller) evaluate(tla *TopLevelArgment, jsonnetAsts []*Node) []*JsonnetEvaluateResult {
	results := make([]*JsonnetEvaluateResult, 0, len(jsonnetAsts))
	for _, ja := range jsonnetAsts {
		tla := &TopLevelArgment{
			Data:   tla.Data,
			Custom: ja.Custom,
		}
		tlaB, err := json.Marshal(tla)
		if err != nil {
			results = append(results, &JsonnetEvaluateResult{
				Key:   ja.Key,
				Error: fmt.Errorf("marshal a top level argument as JSON: %w", err).Error(),
			})
			continue
		}
		vm := jsonnet.NewVM(string(tlaB), c.importer)
		result, err := vm.Evaluate(ja.Node)
		if err != nil {
			results = append(results, &JsonnetEvaluateResult{
				Key:   ja.Key,
				Error: err.Error(),
			})
			continue
		}
		results = append(results, &JsonnetEvaluateResult{
			Key:    ja.Key,
			Result: result,
		})
	}
	return results
}
