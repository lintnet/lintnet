---
sidebar_position: 100
---

# Native functions

:::info
If you want new functions, please [create issues](https://github.com/lintnet/lintnet/issues/new?assignees=&labels=enhancement%2Cnew-native-function&projects=&template=new-native-function-request.yml).
:::

:::caution
We don't write the document of each native functions because we have ported too many functions from Go standard library to maintain the document of them.
Please see GoDoc and [API design](https://github.com/lintnet/go-jsonnet-native-functions?tab=readme-ov-file#api-design).
:::

- [filepath.Base](https://pkg.go.dev/path/filepath#Base)
- [jsonschema.Validate](#jsonschemavalidate)
- [path.Base](https://pkg.go.dev/path#Base)
- [path.Clean](https://pkg.go.dev/path#Clean)
- [path.Dir](https://pkg.go.dev/path#Dir)
- [path.Ext](https://pkg.go.dev/path#Ext)
- [path.IsAbs](https://pkg.go.dev/path#IsAbs)
- [path.Match](https://pkg.go.dev/path#Match)
- [path.Split](https://pkg.go.dev/path#Split)
- [regexp.MatchString](https://pkg.go.dev/regexp#MatchString)
- [strings.Contains](https://pkg.go.dev/strings#Contains)
- [strings.ContainsAny](https://pkg.go.dev/strings#ContainsAny)
- [strings.Count](https://pkg.go.dev/strings#Count)
- [strings.Cut](https://pkg.go.dev/strings#Cut)
- [strings.CutPrefix](https://pkg.go.dev/strings#CutPrefix)
- [strings.CutSuffix](https://pkg.go.dev/strings#CutSuffix)
- [strings.EqualFold](https://pkg.go.dev/strings#EqualFold)
- [strings.Fields](https://pkg.go.dev/strings#Fields)
- [strings.LastIndex](https://pkg.go.dev/strings#LastIndex)
- [strings.LastIndexAny](https://pkg.go.dev/strings#LastIndexAny)
- [strings.Repeat](https://pkg.go.dev/strings#Repeat)
- [strings.Replace](https://pkg.go.dev/strings#Replace)
- [strings.TrimPrefix](https://pkg.go.dev/strings#TrimPrefix)
- [strings.TrimSpace](https://pkg.go.dev/strings#TrimSpace)
- [url.Parse](#urlparse)

## API design

Please see [API design](https://github.com/lintnet/go-jsonnet-native-functions?tab=readme-ov-file#api-design).

## jsonschema.Validate

```go
func(schema, v any) error
```

https://pkg.go.dev/github.com/santhosh-tekuri/jsonschema/v5#Schema.Validate

e.g.

```jsonnet
local schema = import 'main_config_schema.json'; // Import JSON Schema
local validateJSONSchema = std.native('jsonschema.Validate');
local vr = validateJSONSchema(schema, param.config); // Validate param.config with JSON Schema main_config_shema.json
```

Validate validates `v` with JSON Schema `schema` and returns the result.
`schema` is a object representing a JSON Schema. You can define it in Jsonnet or read a JSON Schema with `import`.
Validate returns a error message (string) if something is wrong, or returns [a detailed error object](https://pkg.go.dev/github.com/santhosh-tekuri/jsonschema/v5#Detailed) if `v` violates JSON Schema.
If there is no violation, Validate returns `null`.

## url.Parse

https://pkg.go.dev/net/url#Parse

This function converts [*url.URL](https://pkg.go.dev/net/url#URL) to an object and returns it.

e.g.

```json
[
  {
    "Scheme":      "http",
    "Opaque":      "",
    "Host":        "example.com",
    "Path":        "/foo/bar",
    "RawPath":     "",
    "OmitHost":    false,
    "ForceQuery":  false,
    "RawQuery":    "lang=en&tag=go",
    "Fragment":    "top",
    "RawFragment": "",
    "Query": {
      "lang": ["en"],
      "tag": ["go"],
    },
  },
  null
]
```
