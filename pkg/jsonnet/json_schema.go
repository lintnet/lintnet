package jsonnet

import (
	"encoding/json"
	"errors"

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
				return "compile a JSON Schema: " + err.Error(), nil
			}

			if err := sch.Validate(s[1]); err != nil {
				ve := &jsonschema.ValidationError{}
				if errors.Is(err, ve) {
					var a any
					b, err := json.Marshal(ve.DetailedOutput())
					if err != nil {
						return nil, err
					}
					if err := json.Unmarshal(b, &a); err != nil {
						return nil, err
					}
					return a, nil
				}
				return err.Error(), nil
			}
			return nil, nil
		},
	}
}
