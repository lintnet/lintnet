package jsonnet

import (
	"encoding/json"
	"errors"

	"github.com/google/go-jsonnet"
	"github.com/google/go-jsonnet/ast"
	"github.com/lintnet/go-jsonnet-native-functions/util"
	"github.com/santhosh-tekuri/jsonschema/v6"
)

func ValidateJSONSchema(name string) *jsonnet.NativeFunction {
	return &jsonnet.NativeFunction{
		Name:   name,
		Params: ast.Identifiers{"schema", "v"},
		Func: func(s []any) (any, error) {
			schema, err := json.Marshal(s[0])
			if err != nil {
				return util.NewError("marshal a JSON Schema as JSON: " + err.Error()), nil //nolint:nilerr
			}
			sch, err := jsonschema.CompileString("", string(schema))
			if err != nil {
				return util.NewError("compile a JSON Schema: " + err.Error()), nil //nolint:nilerr
			}

			if err := sch.Validate(s[1]); err != nil {
				return handleJSONSchemaError(err)
			}
			return nil, nil
		},
	}
}

func handleJSONSchemaError(err error) (any, error) {
	ve := &jsonschema.ValidationError{}
	if errors.As(err, &ve) {
		var a any
		b, err := json.Marshal(ve.DetailedOutput())
		if err != nil {
			return util.NewError("marshal a DetailedOutput as JSON: " + err.Error()), nil //nolint:nilerr
		}
		if err := json.Unmarshal(b, &a); err != nil {
			return util.NewError("unmarshal a DetailedOutput as JSON: " + err.Error()), nil //nolint:nilerr
		}
		return a, nil
	}
	return util.NewError(err.Error()), nil
}
