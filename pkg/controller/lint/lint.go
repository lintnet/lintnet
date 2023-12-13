package lint

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/google/go-jsonnet"
	"github.com/google/go-jsonnet/ast"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
	"gopkg.in/yaml.v3"
)

type (
	FileResult struct {
		Results map[string]*Result `json:"results,omitempty"`
		Error   error              `json:"error,omitempty"`
	}
	Result struct {
		Result interface{} `json:"result,omitempty"`
		Error  error       `json:"error,omitempty"`
	}

	NewDecoder func(io.Reader) decoder
	decoder    interface {
		Decode(dest interface{}) error
	}
)

func (c *Controller) Lint(_ context.Context, args ...string) error {
	filePaths, err := c.findJsonnet()
	if err != nil {
		return err
	}
	jsonnetAsts, err := c.readJsonnets(filePaths)
	if err != nil {
		return err
	}

	results := make(map[string]*FileResult, len(args))
	for _, arg := range args {
		rs, err := c.lint(arg, jsonnetAsts)
		if err != nil {
			return logerr.WithFields(err, logrus.Fields{ //nolint:wrapcheck
				"file_path": arg,
			})
		}
		results[arg] = &FileResult{
			Results: rs,
			Error:   err,
		}
	}
	if err := json.NewEncoder(c.stdout).Encode(results); err != nil {
		return fmt.Errorf("encode the result as JSON: %w", err)
	}

	return nil
}

func (c *Controller) findJsonnet() ([]string, error) {
	filePaths := []string{}
	if err := filepath.WalkDir("lintnet", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.Type().IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".jsonnet") {
			return nil
		}
		filePaths = append(filePaths, path)
		return nil
	}); err != nil {
		return nil, fmt.Errorf("walks the file tree of the unarchived package: %w", err)
	}
	return filePaths, nil
}

func getNewDecoder(fileName string) (NewDecoder, error) {
	switch {
	case strings.HasSuffix(fileName, ".json"):
		return func(r io.Reader) decoder {
			return json.NewDecoder(r)
		}, nil
	case strings.HasSuffix(fileName, ".yaml"):
		return func(r io.Reader) decoder {
			return yaml.NewDecoder(r)
		}, nil
	default:
		return nil, errors.New("lintnet supports linting only JSON or YAML")
	}
}

func (c *Controller) readJsonnets(filePaths []string) (map[string]ast.Node, error) {
	jsonnetAsts := make(map[string]ast.Node, len(filePaths))
	for _, filePath := range filePaths {
		ja, err := c.readJsonnet(filePath)
		if err != nil {
			return nil, logerr.WithFields(err, logrus.Fields{ //nolint:wrapcheck
				"file_path": filePath,
			})
		}
		jsonnetAsts[filePath] = ja
	}
	return jsonnetAsts, nil
}

func (c *Controller) readJsonnet(filePath string) (ast.Node, error) {
	b, err := afero.ReadFile(c.fs, filePath)
	if err != nil {
		return nil, fmt.Errorf("read a jsonnet file: %w", err)
	}
	ja, err := jsonnet.SnippetToAST(filePath, string(b))
	if err != nil {
		return nil, fmt.Errorf("parse a jsonnet file: %w", err)
	}
	return ja, nil
}

func (c *Controller) lint(arg string, jsonnetAsts map[string]ast.Node) (map[string]*Result, error) {
	input, err := c.parse(arg)
	if err != nil {
		return nil, err
	}

	vm := jsonnet.MakeVM()
	vm.ExtCode("input", string(input))
	results := make(map[string]*Result, len(jsonnetAsts))
	for k, ja := range jsonnetAsts {
		result, err := vm.Evaluate(ja)
		if err != nil {
			results[k] = &Result{
				Result: result,
				Error:  err,
			}
			continue
		}
		var a interface{}
		if err := json.Unmarshal([]byte(result), &a); err != nil {
			results[k] = &Result{
				Result: result,
				Error:  err,
			}
			continue
		}
		results[k] = &Result{
			Result: a,
			Error:  err,
		}
	}
	return results, nil
}

func (c *Controller) parse(arg string) ([]byte, error) {
	newDecoder, err := getNewDecoder(arg)
	if err != nil {
		return nil, logerr.WithFields(err, logrus.Fields{ //nolint:wrapcheck
			"file_path": arg,
		})
	}
	f, err := c.fs.Open(arg)
	if err != nil {
		return nil, fmt.Errorf("open a yaml file: %w", err)
	}
	defer f.Close()
	var input interface{}
	if err := newDecoder(f).Decode(&input); err != nil {
		return nil, fmt.Errorf("decode a file: %w", err)
	}
	inputB, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("marshal input as JSON: %w", err)
	}
	return inputB, nil
}
