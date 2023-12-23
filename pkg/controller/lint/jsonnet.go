package lint

import (
	"github.com/lintnet/lintnet/pkg/jsonnet"
	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
)

func (c *Controller) readJsonnets(filePaths []*LintFile) (map[string]jsonnet.Node, error) {
	jsonnetAsts := make(map[string]jsonnet.Node, len(filePaths))
	for _, filePath := range filePaths {
		ja, err := jsonnet.Read(c.fs, filePath.Path)
		if err != nil {
			return nil, logerr.WithFields(err, logrus.Fields{ //nolint:wrapcheck
				"file_path": filePath,
			})
		}
		if filePath.ModulePath != "" {
			jsonnetAsts[filePath.ModulePath] = ja
			continue
		}
		jsonnetAsts[filePath.Path] = ja
	}
	return jsonnetAsts, nil
}

func (c *Controller) evaluate(tla string, jsonnetAsts map[string]jsonnet.Node) map[string]*JsonnetEvaluateResult {
	vm := jsonnet.NewVM(tla, c.importer)

	results := make(map[string]*JsonnetEvaluateResult, len(jsonnetAsts))
	for k, ja := range jsonnetAsts {
		result, err := vm.Evaluate(ja)
		if err != nil {
			results[k] = &JsonnetEvaluateResult{
				Error: err.Error(),
			}
			continue
		}
		results[k] = &JsonnetEvaluateResult{
			Result: result,
		}
	}
	return results
}
