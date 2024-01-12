package lint

import (
	"encoding/json"
	"fmt"
	"strings"

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
		Node:    node,
		Key:     lintFile.ID,
		Config:  lintFile.Config,
		Combine: strings.HasSuffix(lintFile.Path, "_combine.jsonnet"),
	}, nil
}

type Node struct {
	Node    jsonnet.Node
	Config  map[string]any
	Key     string
	Combine bool
}

func (c *Controller) evaluateLintFile(tla *TopLevelArgment, lintFile jsonnet.Node) (string, error) {
	if tla.Config == nil {
		tla.Config = map[string]any{}
	}
	tlaB, err := json.Marshal(tla)
	if err != nil {
		return "", fmt.Errorf("marshal a top level argument as JSON: %w", err)
	}
	vm := jsonnet.NewVM(string(tlaB), c.importer)
	result, err := vm.Evaluate(lintFile)
	if err != nil {
		return "", fmt.Errorf("evaluate a lint file as Jsonnet: %w", err)
	}
	return result, nil
}

func (c *Controller) evaluate(tla *TopLevelArgment, lintFiles []*Node) []*Result {
	results := make([]*Result, len(lintFiles))
	for i, lintFile := range lintFiles {
		tla := &TopLevelArgment{
			Data:         tla.Data,
			CombinedData: tla.CombinedData,
			Config:       lintFile.Config,
		}
		s, err := c.evaluateLintFile(tla, lintFile.Node)
		if err != nil {
			results[i] = &Result{
				LintFile: lintFile.Key,
				Error:    err.Error(),
			}
			continue
		}
		rs, a, err := c.parseResult([]byte(s))
		results[i] = &Result{
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
