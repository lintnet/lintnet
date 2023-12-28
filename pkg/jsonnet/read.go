package jsonnet

import (
	"encoding/json"
	"fmt"

	"github.com/google/go-jsonnet"
	"github.com/google/go-jsonnet/ast"
	"github.com/spf13/afero"
)

func Read(fs afero.Fs, filePath, tla string, importer jsonnet.Importer, dest any) error {
	vm := NewVM(tla, importer)
	node, err := ReadToNode(fs, filePath)
	if err != nil {
		return fmt.Errorf("parse a file as Jsonnet: %w", err)
	}
	result, err := vm.Evaluate(node)
	if err != nil {
		return fmt.Errorf("evaluate a file as Jsonnet: %w", err)
	}
	if err := json.Unmarshal([]byte(result), dest); err != nil {
		return fmt.Errorf("unmarshal as JSON: %w", err)
	}
	return nil
}

func ReadToNode(fs afero.Fs, filePath string) (ast.Node, error) {
	b, err := afero.ReadFile(fs, filePath)
	if err != nil {
		return nil, fmt.Errorf("read a jsonnet file: %w", err)
	}
	ja, err := jsonnet.SnippetToAST(filePath, string(b))
	if err != nil {
		return nil, fmt.Errorf("parse a jsonnet file: %w", err)
	}
	return ja, nil
}
