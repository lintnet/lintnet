package jsonnet

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/go-jsonnet"
	"github.com/google/go-jsonnet/ast"
	"github.com/santhosh-tekuri/jsonschema/v5"
)

func ValidateJSONSchema(name string) *jsonnet.NativeFunction {
	return &jsonnet.NativeFunction{
		Name:   name,
		Params: ast.Identifiers{"schema", "v"},
		Func: func(s []any) (any, error) {
			schemaS, ok := s[0].(string)
			if !ok {
				return "the first argument must be a string", nil
			}
			sch, err := jsonschema.Compile(schemaS)
			if err != nil {
				return "compile a JSON Schema: " + err.Error(), nil //nolint:nilerr
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
			return nil, fmt.Errorf("marshal a DetailedOutput as JSON: %w", err)
		}
		if err := json.Unmarshal(b, &a); err != nil {
			return nil, fmt.Errorf("unmarshal DetailedOutput as JSON: %w", err)
		}
		return a, nil
	}
	return err.Error(), nil
}
