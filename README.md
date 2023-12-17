# lintnet

Configuration file linter powered by [Jsonnet](https://jsonnet.org/)

You can write lint rules in Jsonnet and lint configuration files ([Supported formats](#supported-file-format)).

## :warning: This project is still under development

This tool doesn't work and API is unstable yet.
Please don't use this tool yet.

## Features

- Lint any configuration files ([Supported formats](#supported-file-format))
- Powerful. You can lint configuration files flexibly by Jsonnet. And lintnet extends Jsonnet by native functions
- Secure. Jsonnet can't access filesystem and network so it's secure compared with common programming languages such as JavaScript
- Cross Platform. lintnet works on Linux, macOS, and Windows. And it works on both amd64 and arm64
- Easy to install. lintnet is a single binary written in [Go](https://go.dev/), so you only need to install an execurable file into `$PATH`. lintnet has no dependency that you need to install.

## Why Jsonnet?

- Powerful
  - Jsonnet has enough features to lint data
    - e.g. Variables, Functions, Conditions, Array and Object Comprehension, Imports, Errors, External variables, Top-level arguments, Standard library
- Simple
  - The learning cost is not so high
- Popular
  - You can search information and ask help to others when you have some troubles
  - You can utilize the knowledge for not only this tool but also other projects
- Secure
  - Jsonnet can't access file systems and networks and can't execute external commands

## Install

lintnet is a single binary written in [Go](https://go.dev/). So you only need to install an execurable file into `$PATH`.

1. [Homebrew](https://brew.sh/)

```sh
brew install suzuki-shunsuke/lintnet/lintnet
```

2. [Scoop](https://scoop.sh/)

```sh
scoop bucket add suzuki-shunsuke https://github.com/suzuki-shunsuke/scoop-bucket
scoop install lintnet
```

3. [aqua](https://aquaproj.github.io/)

```sh
aqua g -i suzuki-shunsuke/lintnet
```

4. Download a prebuilt binary from [GitHub Releases](https://github.com/suzuki-shunsuke/lintnet/releases) and install it into `$PATH`

## How to use

1. Write lint rules with Jsonnet
1. Run the command `lintnet lint`

```sh
lintnet lint [<file path to be validated> ...]
```

## Supported file format

lintnet can lint the following file formats.
lintnet judges file types by file extensions.
We're considering supporting additional file formats. [#37](https://github.com/suzuki-shunsuke/lintnet/issues/37)

format | file extensions | parser
--- | --- | ---
CSV | `.csv` | [encoding/csv](https://pkg.go.dev/encoding/csv#Reader)
JSON | `.json` | [encoding/json](https://pkg.go.dev/encoding/json#Decoder)
TOML | `.toml` | [github.com/BurntSushi/toml](https://godocs.io/github.com/BurntSushi/toml#Decoder)
TSV | `.tsv` | [encoding/csv](https://pkg.go.dev/encoding/csv#Reader)
YAML | `.yml`, `.yaml` | [gopkg.in/yaml.v3](https://pkg.go.dev/gopkg.in/yaml.v3#Decoder)

### Plain Text

lintnet judges file types by file extensions.
If no parser is found, lintnet parse the file as a plain text file.
The external variable `file_type` is `plain_text`.
The external variable `input` is empty, but you can still lint the file with other external variables such as `file_path` and `file_text`.

## Lint rules

> [!WARNING]
> The specification is unstable yet.

lintnet uses Jsonnet to write lint rules.

### Location of lint files

`lintnet` looks for lint files `*.jsonnet` recursively from the base directory `lintnet`.
You can change the base directory with the command line option `--rule-base-dir (-d)`.

e.g. Change the base directory to `policy`

```sh
lintnet lint -d policy foo.yaml bar.yaml
```

### External Variables

The following [External Variables](https://jsonnet.org/ref/language.html#external-variables-extvars) are passed to lint files.

- `input`: A file content to be linted
- `file_path`: A file path to be linted
- `file_type`: A file type to be linted. One of `json` and `yaml`
- `file_text`: A file content to be linted

### Native functions

lintnet supports all [native functions](https://pkg.go.dev/github.com/google/go-jsonnet#NativeFunction) supported by [suzuki-shunsuke/go-jsonnet-native-functions](https://github.com/suzuki-shunsuke/go-jsonnet-native-functions), which ports Go stanard libraries to Jsonnet.
The following native functions are available.

- strings.Contains
- strings.TrimPrefix
- strings.TrimSpace
- regexp.MatchString
- filepath.Base

You can executed these functions by `std.native("{native function name}")`.

e.g.

```jsonnet
local contained = std.native("strings.Contains")("hello", "ll"); // true
```

### Format of Jsonnet

JSONPath | type | description
--- | --- | ---
`.name` | string | Rule name
`.description` | string | Group description
`.message` | string | Error message
`.failed` | bool | If this is true, this means the file violates the rule
`.level` | string | Error level
`.location` | `string|any` | Location where errors occur
`.metadata` | `any` | Metadata that users can set freely

### Error level

lintnet supports the following error levels.

severity | error level
--- | ---
1 | debug
2 | info
3 | warn
4 | error

You can specify the error level by command line option `--error-level (-e)` or the environment variable `LINTNET_ERROR_LEVEL`.

e.g.

```sh
lintnet lint -e error
```

Error level is similar to log level.
Errors whose error level is lower than the error level of the command are ignored.
The default error level is `info`.

## Configuration file

lintnet reads a configuration file `^\.?lintnet\.ya?ml$` on the current directory.
You can also specify the configuration file path by the command line option `--config (-c)` and the environment variable `LINTNET_CONFIG`.

```sh
lintnet -c foo.yaml lint
```

e.g.

```yaml
error_level: info
modules:
  - id: toml
    source: github_content>suzuki-shunsuke/lintnet-example//toml.jsonnet#v0.1.0
  - id: suzuki-shunsuke/lintnet-example-2
    source: github_content>suzuki-shunsuke/lintnet-example-2#v0.1.0
  - id: suzuki-shunsuke/lintnet-example-3
    source: github_archive>suzuki-shunsuke/lintnet-example-3#v0.1.0
  - id: yaml
    source: http>https://example.com/v0.1.0/yaml.jsonnet
data_transformation:
  - id: foo
    type: jsonnet
    file: transform.jsonnet
outputs:
  - dest: stdout # stdout, stderr, file, github issue
    format: json
    template: foo.jsonnet
  - dest: file
    format: markdown
    path: foo.md
    template: foo.tmpl
targets:
  - lint_files:
      search_type: equal
      paths:
        - path: lintnet/csv.jsonnet
      # imports:
      #   - module: suzuki-shunsuke/lintnet-example-2
      #     path: util.libsonnet
      #     import: utils.libsonnet # optional
    data_files:
      search_type: equal
      paths:
        - path: examples/hello.csv
  # - lint_files:
  #     module: toml
  #   data_files:
  #     search_type: glob
  #     paths:
  #       - path: "*.toml"
  # - lint_files:
  #     module: suzuki-shunsuke/lintnet-example-2
  #     search_type: equal
  #     paths:
  #       - path: actions.jsonnet
  #   data_files:
  #     search_type: glob
  #     paths:
  #       - path: .github/workflows/*.yml
  #       - path: .github/workflows/*.yaml
  #       - path: .github/workflows/fooyaml
  #         exclude: true
  # - lint_files:
  #     module: suzuki-shunsuke/lintnet-example-3
  #     paths:
  #       - path: lintnet/json.jsonnet
  #   data_files:
  #     search_type: glob
  #     paths:
  #       - path: "*/*.json"
  # - lint_files:
  #     module: yaml
  #   data_files:
  #     search_type: glob
  #     paths:
  #       - path: "*/*.yaml"
```

### Example

Please see [lintnet](lintnet).

## LICENSE

[MIT](LICENSE)
