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

func (c *Controller) Lint(_ context.Context, args ...string) error {
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
		return fmt.Errorf("walks the file tree of the unarchived package: %w", err)
	}
	jsonnetAsts := make([]ast.Node, len(filePaths))
	for i, filePath := range filePaths {
		b, err := afero.ReadFile(c.fs, filePath)
		if err != nil {
			return fmt.Errorf("read a jsonnet file: %w", err)
		}
		ja, err := jsonnet.SnippetToAST(filePath, string(b))
		if err != nil {
			return fmt.Errorf("parse a jsonnet file: %w", err)
		}
		jsonnetAsts[i] = ja
	}

	for _, arg := range args {
		var newDecoder func(io.Reader) decoder
		if strings.HasSuffix(arg, ".json") {
			newDecoder = func(r io.Reader) decoder {
				return json.NewDecoder(r)
			}
		} else if strings.HasSuffix(arg, ".yaml") {
			newDecoder = func(r io.Reader) decoder {
				return yaml.NewDecoder(r)
			}
		} else {
			return logerr.WithFields(errors.New("lintnet supports linting only JSON or YAML"), logrus.Fields{ //nolint:wrapcheck
				"file_path": arg,
			})
		}
		f, err := c.fs.Open(arg)
		if err != nil {
			return fmt.Errorf("open a yaml file: %w", err)
		}
		defer f.Close()
		var input interface{}
		if err := newDecoder(f).Decode(&input); err != nil {
			return fmt.Errorf("decode a file: %w", err)
		}
		inputB, err := json.Marshal(input)
		if err != nil {
			return fmt.Errorf("marshal input as JSON: %w", err)
		}

		vm := jsonnet.MakeVM()
		vm.ExtCode("input", string(inputB))
		for _, ja := range jsonnetAsts {
			result, err := vm.Evaluate(ja)
			if err != nil {
				return fmt.Errorf("evaluate Jsonnet: %w", err)
			}
			fmt.Println(result) //nolint:forbidigo
		}
	}

	return nil
}

type decoder interface {
	Decode(interface{}) error
}

func decodeJSON(f io.Reader, dest interface{}) error {
	return json.NewDecoder(f).Decode(dest)
}

func decodeYAML(f io.Reader, dest interface{}) error {
	return yaml.NewDecoder(f).Decode(dest)
}
