---
sidebar_position: 100
---

# Parameterize lint rules

Each lint file can take config parameters.

e.g.

```jsonnet
lint_files: [
  {
    path: 'examples/lint/filename.jsonnet',
    config: {
      excludes: ['foo'],
    },
  },
],
modules: [
  {
    path: 'github_archive/github.com/suzuki-shunsuke/example-lintnet-modules/ghalint.jsonnet@32ca3be646ec5b5861aab72fed30cd71f6eba9bf:v0.1.2',
    config: {
      excludes: ['foo'],
    },
  },
],
```

Each lint file can refer to parameters by `param.config`.

e.g.

```jsonnet
local excludes = std.get(param.config, 'excludes', [])
```
