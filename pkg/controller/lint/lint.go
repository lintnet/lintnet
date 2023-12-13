package lint

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/go-jsonnet"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

func (c *Controller) Lint(_ context.Context) error {
	jsonnetFileName := "test.jsonnet"
	b, err := afero.ReadFile(c.fs, jsonnetFileName)
	if err != nil {
		return fmt.Errorf("read a jsonnet file: %w", err)
	}

	ja, err := jsonnet.SnippetToAST(jsonnetFileName, string(b))
	if err != nil {
		return fmt.Errorf("parse a jsonnet file: %w", err)
	}

	f, err := c.fs.Open("test.yaml")
	if err != nil {
		return fmt.Errorf("open a yaml file: %w", err)
	}
	defer f.Close()
	var input interface{}
	if err := yaml.NewDecoder(f).Decode(&input); err != nil {
		return fmt.Errorf("decode a yaml file: %w", err)
	}
	inputB, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("marshal input as JSON: %w", err)
	}

	vm := jsonnet.MakeVM()
	vm.ExtCode("input", string(inputB))
	result, err := vm.Evaluate(ja)
	if err != nil {
		return fmt.Errorf("evaluate Jsonnet: %w", err)
	}
	fmt.Println(result) //nolint:forbidigo
	return nil
}
