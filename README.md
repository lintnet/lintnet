# lintnet

Generic Linter powered by [Jsonnet](https://jsonnet.org/)

You can write lint rules with Jsonnet and lint files (JSON and YAML).

## :warning: This project is still under development

This tool doesn't work and API is unstable yet.
Please don't use this tool yet.

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

Coming soon.

## How to use

1. Write lint rules with Jsonnet
1. Run the command `lintnet lint`

```sh
lintnet lint [<file path to be validated> ...]
```

## Lint rules

> [!WARNING]
> The specification is unstable yet.

lintnet uses Jsonnet to write lint rules.

### Location of lint files

`lintnet` looks for lint files `*.jsonnet` recursively from the base directory `lintnet`.
You can change the base directory with the command line option `--rule-baes-dir (-d)`.

e.g. Change the base directory to `policy`

```sh
lintnet lint -d policy foo.yaml bar.yaml
```

### External Variables

The following [External Variables](https://jsonnet.org/ref/language.html#external-variables-extvars) are passed to lint files.

- `input`: A file content to be linted
- `file_path`: A file path to be linted
- `file_type`: A file type to be linted. One of `json` and `yaml`

### Format of Jsonnet

JSONPath | type | description
--- | --- | ---
`.name` | string | Group name
`.description` | string | Group description
`.rules[].name` | string | Rule name
`.rules[].description` | string | Rule description
`.rules[].errors[].message` | string | Error message

If `.rules[].errors` isn't empty, lintnet treats as the lint fails.

### Example

```jsonnet
local fileType = std.extVar('file_type');
local input = std.extVar('input');

{
  group_name: "GitHub Actions",
  description: |||
    Lint rules regarding GitHub Actions
  |||,
  rules: [
    {
      name: "uses_must_not_be_main",
      description: |||
        actions reference must not be main.
      |||,
      errors: [
        {
          message: "",
          job_name: job.key,
          uses: step.uses,
        }
        for job in std.objectKeysValues(input.jobs)
        if std.objectHas(job.value, 'steps')
        for step in job.value.steps
        if fileType == 'yaml' && std.objectHas(step, 'uses') && std.endsWith(step.uses, '@main')
      ]
    }
  ],
}
```

## LICENSE

[MIT](LICENSE)
