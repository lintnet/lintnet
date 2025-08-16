---
sidebar_position: 400
---

# Configuration

## Environment variables

- `LINTNET_CONFIG`: Configuration file path
- [LINTNET_ERROR_LEVEL](guides/error-level.md): `debug|info|warn|error`
- [LINTNET_SHOWN_ERROR_LEVEL](guides/error-level.md): `debug|info|warn|error`
- `LINTNET_OUTPUT_SUCCESS`: `true|false`
- `LINTNET_LOG_LEVEL`: `trace|debug|info|warn|error|fatal|panic`
- `LINTNET_LOG_COLOR`: `auto|always|never`
- `LINTNET_GITHUB_TOKEN`: GitHub Access Token for getting Modules
- `LINTNET_ROOT_DIR`: Root directory where modules are installed
- `GITHUB_TOKEN`: GitHub Access Token for getting Modules

## Configuration file path

lintnet reads a configuration file `^\.?lintnet\.jsonnet$` on the current directory.
You can also specify the configuration file path by the command line option `--config (-c)` and the environment variable `LINTNET_CONFIG`.

```sh
lintnet -c foo.yaml lint
```

## Scaffold a configuration file

You can scaffold the configuration file by `lintnet init` command.

```sh
lintnet init
```

## Configuration file format

[JSON Schema](https://github.com/lintnet/lintnet/blob/main/json-schema/lintnet.json)

The configuration is a function returning an object.

```jsonnet
function(param) {
  // targets is a list of lint configuration.
  targets: [
    //
  ],
  // ignored_dirs is a list of ignored directories.
  // ignored_dirs is optional.
  // The default value is [".git", "node_modules"].
  ignored_dirs: [
    ".git",
    "node_modules",
  ],
  // outputs is a list of output configuration.
  // outputs is optional.
  // outputs is used to customize output format.
  outputs: [
    // ...
  ],
}
```

### .targets

The element of `targets` is a pair of lint files and data files.
`data_files` is required.
Either `lint_files` or `modules` is required.
Both `lint_files` and `modules` can also be used.

```jsonnet
{
  id: 'target id', // optional
  base_data_path: '', // optional
  // data_files is a list of glob patterns.
  data_files: [
    'examples/**/*.csv', // relative path from the configuration file
    // ...
  ],
  // lint_files is a list of local lint files.
  lint_files: [
    'main.jsonnet', // relative path from the configuration file
  ],
  // modules is a list of modules.
  modules: [
    'github_archive/github.com/suzuki-shunsuke/example-lintnet-modules/newline.jsonnet@32ca3be646ec5b5861aab72fed30cd71f6eba9bf:v0.1.2',
  ],
}
```

### .targets[].lint_files

An element of `lint_files` is either a string or an object.
If you want to pass configuration to lint rules, you need to use an object. Otherwise, a string is enough.

```jsonnet
[
  'main.jsonnet',
  {
    path: 'examples/lint/filename.jsonnet', // required
    config: {}, // optional. lint rule configuration
  },
],
```

### .targets[].modules

An element of `modules` is either a string or an object.

```jsonnet
[
  // github_archive/github.com/<repo owner>/<repo name>/<path>@<full commit hash>[:<version>]
  // version is optional.
  'github_archive/github.com/lintnet-modules/ghalint/workflow/action_ref_should_be_full_length_commit_sha/main.jsonnet@00571db321e413d45be457f39e48cd4237399bb7:v0.3.0',
  {
    // path is required
    path: 'github_archive/github.com/lintnet-modules/ghalint/workflow/action_ref_should_be_full_length_commit_sha/main.jsonnet@00571db321e413d45be457f39e48cd4237399bb7:v0.3.0',
    config: {},
  },
  {
    path: 'github_archive/github.com/lintnet-modules/ghalint@00571db321e413d45be457f39e48cd4237399bb7:v0.3.0',
    // You can specify file paths in a module with the attribute files.
    // This style is useful to specify multiple file path patterns in a module and set config parameter by lint rule
    files: [
      'workflow/**/main.jsonnet',
      '!workflow/action_ref_should_be_full_length_commit_sha/main.jsonnet',
      {
        path: 'workflow/action_ref_should_be_full_length_commit_sha/main.jsonnet',
        config: {},
      },
    ],
  },
],
```

### .targets[].base_data_path

`base_data_path` is useful to lint files by service directory in a Monorepo.

```jsonnet
{
  // data files which are on the same directory as tfaction.yaml.
  base_data_path: '**/tfaction.yaml',
  data_files: [
    // relative path from base_data_path
    // Glob is also available
    '*.tf',
  ],
  // ...
},
```

In case of [linting across multiple files](/docs/guides/lint-across-files/), `base_data_path` is useful to separate files.
In the above case, if `**/tfaction.yaml` matches `foo/tfaction.yaml` and `bar/tfaction.yaml`, `foo/*.tf` and `bar/*.tf` are linted separately.

### .outputs

Please see [Customize Output](/docs/guides/customize-output/).

## File paths in configuration files

The file path separator must be a slash `/`.
If file paths are relative paths, the base must be the configuration file.

## Top level argument

Now the top level argument `param` is empty. This argument is reserved for future enhancement.

## Glob

`data_files` and `lint_files` are lists of patterns matching with data and lint files.
Each string is parsed with [bmatcuk/doublestar](https://github.com/bmatcuk/doublestar).

### Exclude by `!`

If each string starts with `!`, files matching with the pattern are excluded.

e.g. foo/example.jsonnet is excluded

```
**/*.jsonnet
!foo/example.jsonnet
```

e.g. foo/example.jsonnet isn't excluded because the later pattern `foo/*.jsonnet` takes precedence

```
**/*.jsonnet
!foo/example.jsonnet
foo/*.jsonnet
```

### Excluded directories and files

lintnet doesn't check `.gitignore`.
lintnet ignores the following directories.

- `.git`
- `node_modules`

And in `lint_files`, files `*_test.jsonnet` are ignored.

## See also

- [Command line arguments](guides/usage.md)
- [Parameterize lint rules](guides/parameterize-rule.md)
- [Customize Ouptputs](guides/customize-output.md)
- [Lint across files](guides/lint-across-files.md)
