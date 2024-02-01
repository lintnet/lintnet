# lintnet

Powerful, Secure, Shareable linter powered by [Jsonnet](https://jsonnet.org/)

<p align="center" width="100%">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/lintnet/logo/add-logo/images/lintnet.png">
    <img src="https://raw.githubusercontent.com/lintnet/logo/add-logo/images/lintnet.png" alt="logo" width="400">
  </picture>
</p>

## :warning: This project is still under development

This tool doesn't work and API is unstable yet.
Please don't use this tool yet.

## Features

- [Support various configuration file formats](https://lintnet.github.io/docs/supported-data-format)
- Powerful. You can lint configuration files flexibly by Jsonnet. And lintnet extends Jsonnet by native functions
- Secure. Jsonnet can't access filesystem and network so it's secure compared with common programming languages such as JavaScript
- Shareable. lintnet provides Module system that you can share lint rules between other projects. You can develop lint rules as both OSS and in-house libraries
- Cross Platform. lintnet works on Linux, macOS, Windows / amd64, arm64
- Easy to install. lintnet is a single binary written in [Go](https://go.dev/), so you only need to install an execurable file into `$PATH`. lintnet has no dependency that you need to install.

## Document

https://lintnet.github.io/

## LICENSE

[MIT](LICENSE)

About the license of logo, please see [here](https://github.com/lintnet/logo).
