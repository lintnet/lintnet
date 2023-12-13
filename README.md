# lintnet

Generic Linter powered by [Jsonnet](https://jsonnet.org/)

You can write lint rules with Jsonnet and lint files (JSON and YAML).

## :warning: This project is still under development

This tool doesn't work and API is unstable yet.
Please don't use this tool yet.

## Why Jsonnet?

- Powerful
  - Jsonnet has enough features to lint data
    - e.g. variables, functions, conditions, Array and Object Comprehension, Imports, Errors, External variables, Top-level arguments, Standard library
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

```sh
lintnet [<file path to be validated> ...]
```

## Lint rules

> [!WARNING]
> The specification is unstable yet.

lintnet uses Jsonnet to write lint rules.

### External Variables

- `input`: A file content to be linted
- `file_path`: A file path to be linted
- `file_type`: A file type to be linted. One of `json` and `yaml`

### Format

Coming soon.

### Example

```jsonnet
local fileType = std.extVar('filetype');
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
        for job in std.objectKeysValues(std.extVar('input').jobs)
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
