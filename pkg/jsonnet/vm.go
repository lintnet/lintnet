package jsonnet

import (
	"encoding/json"
	"fmt"

	"github.com/google/go-jsonnet"
	"github.com/google/go-jsonnet/ast"
	"github.com/lintnet/go-jsonnet-native-functions/pkg/path/filepath"
	"github.com/lintnet/go-jsonnet-native-functions/pkg/regexp"
	"github.com/lintnet/go-jsonnet-native-functions/pkg/strings"
	"github.com/santhosh-tekuri/jsonschema/v5"
)

func SetNativeFunctions(vm *jsonnet.VM) {
	vm.NativeFunction(strings.Contains("strings.Contains"))
	vm.NativeFunction(strings.TrimPrefix("strings.TrimPrefix"))
	vm.NativeFunction(strings.TrimSpace("strings.TrimSpace"))
	vm.NativeFunction(regexp.MatchString("regexp.MatchString"))
	vm.NativeFunction(filepath.Base("filepath.Base"))
	vm.NativeFunction(ValidateJSONSchema("github.com/santhosh-tekuri/jsonschema/v5.ValidateJSONSchema"))
}

func MakeVM() *jsonnet.VM {
	return jsonnet.MakeVM()
}

func NewVM(param string, importer jsonnet.Importer) *jsonnet.VM {
	vm := jsonnet.MakeVM()
	vm.TLACode("param", param)
	SetNativeFunctions(vm)
	vm.Importer(importer)
	return vm
}

func ValidateJSONSchema(name string) *jsonnet.NativeFunction {
	return &jsonnet.NativeFunction{
		Name:   name,
		Params: ast.Identifiers{"schema", "v"},
		Func: func(s []any) (any, error) {
			schemaS, err := json.Marshal(s[0])
			if err != nil {
				return nil, fmt.Errorf("marshal a JSON Schema as JSON: %w", err)
			}
			sch, err := jsonschema.CompileString("schema.json", string(schemaS))
			if err != nil {
				return nil, fmt.Errorf("compile a JSON Schema: %w", err)
			}

			if err = sch.Validate(s[1]); err != nil {
				return nil, fmt.Errorf("validate data by the JSON Schema: %w", err)
			}
			return nil, nil
		},
	}
}
