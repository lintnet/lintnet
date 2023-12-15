package lint

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/google/go-jsonnet"
	"github.com/google/go-jsonnet/ast"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
)

func (c *Controller) findJsonnet(baseDir string) ([]string, error) {
	filePaths := []string{}
	if err := filepath.WalkDir(baseDir, func(path string, d fs.DirEntry, err error) error {
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

func newVM(input []byte, filePath, fileType string) *jsonnet.VM {
	vm := jsonnet.MakeVM()
	vm.ExtCode("input", string(input))
	vm.ExtVar("file_path", filePath)
	vm.ExtVar("file_type", fileType)
	setNativeFunctions(vm)
	return vm
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
