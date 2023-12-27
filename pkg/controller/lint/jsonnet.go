package lint

import (
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
				Node: ja,
				Key:  filePath.ModulePath,
			})
			continue
		}
		jsonnetAsts = append(jsonnetAsts, &Node{
			Node: ja,
			Key:  filePath.Path,
		})
	}
	return jsonnetAsts, nil
}

type Node struct {
	Node jsonnet.Node
	Key  string
}

func (c *Controller) evaluate(tla string, jsonnetAsts []*Node) []*JsonnetEvaluateResult {
	vm := jsonnet.NewVM(tla, c.importer)

	results := make([]*JsonnetEvaluateResult, 0, len(jsonnetAsts))
	for _, ja := range jsonnetAsts {
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
