---
sidebar_position: 1
---

# Home

[Release Notes](https://github.com/lintnet/lintnet/releases) | [Versioning Policy](https://github.com/suzuki-shunsuke/versioning-policy) | [MIT LICENSE](https://github.com/lintnet/lintnet/blob/main/LICENSE)

lintnet is a general-purpose linter that brings [Policy as Code](https://developer.hashicorp.com/sentinel/docs/concepts/policy-as-code) to your software development.
It is a command-line tool like [Conftest](https://www.conftest.dev/) but offers greater reusability and an enhanced user experience. 
It is available for Terraform, Kubernetes, GitHub Actions, and any kind of configuration files.
Its versatility allows it to cover many use cases, eliminating the need for multiple different linters.
Unlike other linters, lintnet itself does not come with built-in lint rules; instead, it runs user-defined lint rules.
This means you no longer need to develop linters from scratch.
Instead, you can focus on developing lint rules while lintnet handles the rest.
These lint rules can be reused and published as Modules.
Our goal is to create an ecosystem for lint rules where everyone can easily use and publish them, thereby promoting the Policy as Code approach in software development.
You can define lint rules using [Jsonnet](https://jsonnet.org/), a simple, powerful, and secure configuration language.
lintnet enhances Jsonnet's capabilities with [go-jsonnet](https://github.com/google/go-jsonnet)'s native functions, making it even more powerful.

## Features

- [Support various formats](supported-data-format.md)
- [Define lint rules by Jsonnet](#why-jsonnet)
- [Easy to install](install.md). lintnet is a single binary written in [Go](https://go.dev/), so you only need to install an executable file into `$PATH`. lintnet has no dependency that you need to install
- [Unit testing of lint rules](test-rule.md)
- [Share and reuse lint rules as Modules](module.md)
  - [Update Modules by Renovate](module.md#update-modules-by-renovate)
- [Lint across multiple files](guides/lint-across-files.md)
- [Customize output format](guides/customize-output.md)

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
  - Jsonnet retricts access to filesystems and networks, and external command execution

## How does lintnet achieve lint using Jsonnet?

![image](https://github.com/lintnet/lintnet/assets/13323303/d53e3739-c6ae-4d52-86f3-2caa38812251)

## Comparison

### Conftest

- 👍 High reusability
- 👍 Some people would prefer Jsonnet over [Rego](https://www.openpolicyagent.org/docs/latest/policy-language/), though this is subjective and up to you
- 👍 Declarative configuration

#### 👍 High reusability

When we used Conftest, we complaint we couldn't reuse Conftest policies well.

1. Conftest has the mechanism to push and pull policies, but we think this isn't enough. More sophisticated and standardized way is necessary
1. It's a little difficult to share Conftest policies between multiple repositories.
Of course it's possible, but there is no standard way
1. People write similar policies from scratch independently.
This isn't good. Ideally, policies should be shared and reused all over the world

lintnet has the module mechanism. you can distribute and reuse modules so easily in the standard way.
Not only lint rules but also Jsonnet functions can be shared as modules.
You can update modules continuously by Renovate.

About modules, please see [Module](module.md).

#### 👍 Some people would prefer Jsonnet over Rego

This is so subjective and up to you, but some people would feel Jsonnet is easier than Rego.

Rego is awesome, but it's different from other programing languages such as JavaScript and Python, so some people have difficulty in learning Rego.

If you complain about Rego, maybe you like Jsonnet.

### Programing languages such as Python and JavaScript

- 👍 Secure
- 👍 You only need to implement lint logic. You don't need to implement other feature such as reading and parsing files and outputs results

If you reuse third party libraries as lint rules, you need to check if they are secure.
Common programing languages such as Python and JavaScript can do anything, so attackers can execute malicious codes. It would be difficult to ensure security.
On the other hand, Jsonnet restricts access to filesystem and network, and OS command execution so it's securer than those programming languages.

## Sub projects

https://github.com/orgs/lintnet/repositories

- [lintnet](https://github.com/lintnet/lintnet): CLI
- [modules](https://github.com/lintnet-modules): Official modules
- [examples](https://github.com/lintnet/examples): Examples
- [renovate-config](https://github.com/lintnet/renovate-config): Renovate Config Preset to update modules
- [lintnet.github.io](https://github.com/lintnet/lintnet.github.io): Official web site
- [go-jsonnet-native-functions](https://github.com/lintnet/go-jsonnet-native-functions): Go package porting several Go's Standard libraries functions to go-jsonnet's Native functions
