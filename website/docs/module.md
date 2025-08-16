---
sidebar_position: 500
---

# Module

You can share and reuse lint rules and Jsonnet codes such as functions.
We call this mechanism `Module`.
You can host modules on GitHub repositories.

There are two types of modules.

1. Lint rule module
1. Imported module

## 1. Lint rule module

Lint rule module is same with the normal lint rules.

You can use Lint rule modules by specifying them in configuration files.

```jsonnet
function(param) {
  targets: [
    {
      modules: [ // Lint rule modules
        'github_archive/github.com/lintnet-modules/nllint/main.jsonnet@8cfc4eae68ec93f9b92d9048ce51b0d9646c976c:v0.1.0',
      ],
      data_files: [
        '**/*',
      ],
    },
  ],
}
```

## 2. Imported module

You can share variables and functions as Imported modules.
Imported modules are imported by Jsonnet's `import` statement.

```jsonnet
local hello = import 'github_archive/github.com/lintnet/modules/modules/hello/hello.jsonnet@60a46a4fa4c0e7b1b95f57c479e756afa2f376e9:v0.1.0';
```

You can utilize third party Jsonnet libraries unrelated to lintnet too.

e.g.

- https://github.com/jsonnet-libs/xtd ([example](https://github.com/lintnet/examples/tree/main/jsonnet-library/xtd))

## Module path format

```
${type}/${host}/${repository_owner}/${repository_name}/${file_path}@${full_commit_hash}[:${tag}]
```

Now only `github_archive` is valid as `type`, and only `github.com` is valid as `host`.

e.g.

```
github_archive/github.com/lintnet/modules/modules/hello/hello.jsonnet@60a46a4fa4c0e7b1b95f57c479e756afa2f376e9:v0.1.0'
```

## Update modules by Renovate

You can update modules by [Renovate](https://docs.renovatebot.com/) using our [Renovate Preset](https://docs.renovatebot.com/config-presets/).

https://github.com/lintnet/renovate-config

## Where to install modules

Modules are installed on the following directory.

```
${Application Data Directory}/lintnet/modules
```

`${Application Data Directory}` is `XDG_DATA_HOME` in https://github.com/adrg/xdg .

| environment | Application Data Directory                  |
| ----------- | ------------------------------------------- |
| Unix        | `~/.local/share`                            |
| macOS       | `~/Library/Application Support`             |
| Windows     | `LocalAppData`, `%LOCALAPPDATA%` (Fallback) |

Or you can change the directory by the environment variable `LINTNET_ROOT_DIR`.

You can get the install path by `lintnet info -module-root-dir`

```sh
lintnet info -module-root-dir
```

## :bulb: Cache modules in CI

You can cache modules in CI such as GitHub Actions.

e.g.

```yaml
- run: echo "module_root_dir=$(lintnet info -module-root-dir)" >> "$GITHUB_OUTPUT"
  id: lintnet

- uses: actions/cache@v3
  with:
    path: |
      ${{steps.lintnet.outputs.module_root_dir}}
    key: ${{ hashFiles('lintnet.jsonnet') }}

- run: lintnet lint
  env:
    GITHUB_TOKEN: ${{github.token}}
```

## GitHub Access Tokens

lintnet uses GitHub API to download Modules.
To avoid API rate limiting, we recommend setting a GitHub Access Token to the environment variables `LINTNET_GITHUB_TOKEN` or `GITHUB_TOKEN`.
To use modules hosted on private repositories, GitHub Access Tokens with `contents:read` permission are necessary.

## Official Modules

https://github.com/lintnet-modules

We ported some linters such as [ghalint](https://github.com/suzuki-shunsuke/ghalint) and [nllint](https://github.com/suzuki-shunsuke/nllint) to lintnet and shared them as official modules.

- [ghalint](https://github.com/lintnet-modules/ghalint)
- [nllint](https://github.com/lintnet-modules/nllint)
- [github-actions](https://github.com/lintnet-modules/github-actions)
- [k8s](https://github.com/lintnet-modules/k8s)
- [Terraform](https://github.com/lintnet-modules/terraform)
  - [Terraform AWS Provider](https://github.com/lintnet-modules/terraform-aws)
  - [Terraform Google Provider](https://github.com/lintnet-modules/terraform-google)
- etc

## Find Modules

We recommend adding the topic [lintnet-module](https://github.com/topics/lintnet-module) to Module repositories so that everyone can find modules.
So please check the topic [lintnet-module](https://github.com/topics/lintnet-module).

## Develop Modules

1. Create a GitHub Repository
1. Write lint rules
1. (Optional) Add the topic [lintnet-module](https://github.com/topics/lintnet-module) to the repository so that everyone can find your modules
1. (Optional) Write tests
1. (Optional) Set up CI running `lintnet test`
1. (Optional) Write document
1. (Optional) Create GitHub Releases

[The official modules](https://github.com/lintnet-modules) would be a good reference.

### Write document

Of course the format is free, but we recommend writing the following information.

- Description
- Example
- Why is the rule necessary?
- How to fix
- config's schema
