# lintnet

General purpose linter powered by [Jsonnet](https://jsonnet.org/).

<p align="center" width="100%">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/lintnet/logo/main/images/lintnet.png">
    <img src="https://raw.githubusercontent.com/lintnet/logo/main/images/lintnet.png" alt="logo" width="400">
  </picture>
</p>

## Features

- [Support various configuration file formats](https://lintnet.github.io/docs/supported-data-format)
- Powerful. You can lint files flexibly by Jsonnet. And lintnet extends Jsonnet by native functions
- Secure. Jsonnet restricts access to filesystem and network so it's secure compared with common programming languages such as Python
- Shareable. lintnet supports sharing lint rules as Modules. You can utilize third party lint rules, reuse your lint rules in multiple projects, and distribute lint rules as OSS and in-house libraries
- Easy to install. lintnet is a single binary written in [Go](https://go.dev/), so you only need to install an execurable file into `$PATH`. lintnet has no dependency that you need to install

## Document

https://lintnet.github.io/

## LICENSE

[MIT](LICENSE)

About the license of logo, please see [here](https://github.com/lintnet/logo).
