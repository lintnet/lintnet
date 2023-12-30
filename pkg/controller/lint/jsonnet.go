package lint

import (
	"encoding/json"
	"fmt"

	"github.com/lintnet/lintnet/pkg/config"
	"github.com/lintnet/lintnet/pkg/jsonnet"
	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
)

func (c *Controller) parseLintFiles(lintFiles []*config.LintFile) ([]*Node, error) {
	nodes := make([]*Node, 0, len(lintFiles))
	for _, lintFile := range lintFiles {
		node, err := c.parseLintFile(lintFile)
		if err != nil {
			return nil, logerr.WithFields(err, logrus.Fields{ //nolint:wrapcheck
				"file_path": lintFile.Path,
			})
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}

func (c *Controller) parseLintFile(lintFile *config.LintFile) (*Node, error) {
	node, err := jsonnet.ReadToNode(c.fs, lintFile.Path)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}
	return &Node{
		Node:   node,
		Key:    lintFile.ID,
		Config: lintFile.Config,
	}, nil
}

type Node struct {
	Node   jsonnet.Node
	Config map[string]any
	Key    string
}

func (c *Controller) evaluateLintFile(data *Data, lintFile *Node) *JsonnetEvaluateResult {
	tla := &TopLevelArgment{
		Data:   data,
		Config: lintFile.Config,
	}
	if tla.Config == nil {
		tla.Config = map[string]any{}
	}
	tlaB, err := json.Marshal(tla)
	if err != nil {
		return &JsonnetEvaluateResult{
			Key:   lintFile.Key,
			Error: fmt.Errorf("marshal a top level argument as JSON: %w", err).Error(),
		}
	}
	vm := jsonnet.NewVM(string(tlaB), c.importer)
	result, err := vm.Evaluate(lintFile.Node)
	if err != nil {
		return &JsonnetEvaluateResult{
			Key:   lintFile.Key,
			Error: err.Error(),
		}
	}
	return &JsonnetEvaluateResult{
		Key:    lintFile.Key,
		Result: result,
	}
}

func (c *Controller) evaluate(data *Data, lintFiles []*Node) []*JsonnetEvaluateResult {
	results := make([]*JsonnetEvaluateResult, len(lintFiles))
	for i, lintFile := range lintFiles {
		results[i] = c.evaluateLintFile(data, lintFile)
	}
	return results
}
