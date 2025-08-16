---
sidebar_position: 300
---

# Customize Outputs

By default, lintnet outputs JSON when the lint fails.
You can custmize the JSON format with Jsonnet and Go's [text/template](https://pkg.go.dev/text/template) and [html/template](https://pkg.go.dev/html/template).

For detail, please see [the example](https://github.com/lintnet/examples/tree/main/customize-output).

## Output JSON even if lint passes

By default `lintnet lint` command outputs nothing if lint passes.
If you want to output JSON even if lint passes, please use `--output-success` option.

```sh
lintnet lint -output-success
```
