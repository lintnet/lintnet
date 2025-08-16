---
sidebar_position: 700
---

# Error Level

Each lint error has it's `error level`.
`lintnet lint` command has it's `error level` and `shown error level`.

The following error levels are supported.

severity | error level
--- | ---
1 | debug
2 | info
3 | warn
4 | error

The default error level of lint command and each lint error is `error`.
And the default `shown error level` of lint command is `info`.

lint command's `error level` is greater than `shown error level`.
For example, if you set `error level` to `debug`, `shown error level` becomes `debug` too.

You can specify the error level of lint command by command line option `--error-level (-e)` or the environment variable `LINTNET_ERROR_LEVEL` or configuration file `lintnet.jsonnet`.

e.g.

```sh
lintnet lint -e error
```

```jsonnet
// lintnet.jsonnet
error_level: 'warn',
shown_error_level: 'debug',
```

You can also specify the shown error level of lint command by command line option `--shown-error-level` or the environment variable `LINTNET_SHOWN_ERROR_LEVEL`.

The error level of each lintnet error is specified with `level` field.

e.g.

```jsonnet
function(param)
  if std.objectHas(param.data.value, 'description') then [] else [{
    name: 'description is required',
    level: 'warn', // Error level is 'warn'
  }]
```

If all errors' error level is lower than the error level of lint command, the command succeeds.
If the error level of a lint error is lower than the shown error level of lint command, the error is excluded from the output.
