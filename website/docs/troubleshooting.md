---
sidebar_position: 600
---

# Troubleshooting

This page is a guide for when you face any issues.

## Output Debug log

The debug log would be helpful.

```sh
lintnet -log-level debug lint
```

## Check Jsonnet behaviour

See [Check Jsonnet behaviour](learn-jsonnet.md#check-jsonnet-behaviour).

## Output variables with a custom field

```jsonnet
function(param)
  if std.objectHas(param.data.value, 'description') then [] else [{
    name: 'description is required',
    custom: {
      param: param, // Output param for debug
    },
  }]
```

## Output variables with std.trace

https://jsonnet.org/ref/stdlib.html#trace

e.g.

```jsonnet
function(param)
  std.trace(std.toString(param), // Output param to stderr
    if std.objectHas(param.data.value, 'description') then [] else [{
      name: 'description is required',
    }]
  )
```
