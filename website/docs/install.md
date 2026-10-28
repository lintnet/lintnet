---
sidebar_position: 100
---

# Install

lintnet is a single binary written in [Go](https://go.dev/).
So you only need to install an executable file into `$PATH`.

1. [Homebrew](https://brew.sh/)

```sh
brew install lintnet/lintnet/lintnet
```

2. [Scoop](https://scoop.sh/)

```sh
scoop bucket add lintnet https://github.com/lintnet/scoop-bucket
scoop install lintnet
```

3. [aqua](https://aquaproj.github.io/)

```sh
aqua g -i lintnet/lintnet
```

4. Download a prebuilt binary from [GitHub Releases](https://github.com/lintnet/lintnet/releases) and install it into `$PATH`

<details>
<summary>Verify downloaded assets from GitHub Releases</summary>

You can verify downloaded assets using some tools.

1. [GitHub CLI](https://cli.github.com/) `lintnet >= v0.4.8`
1. [Cosign](https://github.com/sigstore/cosign)
1. [Packslip](https://packslip.dev/) `lintnet >= v1.0.1`
1. [slsa-verifier](https://github.com/slsa-framework/slsa-verifier) `lintnet <= v1.0.0`

--

1. GitHub CLI

lintnet >= v0.4.8

You can install GitHub CLI by aqua.

```sh
aqua g -i cli/cli
```

```sh
version=v1.0.1
asset=lintnet_darwin_arm64.tar.gz
gh release download -R lintnet/lintnet "$version" -p "$asset"
gh attestation verify "$asset" \
  -R lintnet/lintnet \
  --signer-workflow suzuki-shunsuke/go-release-workflow/.github/workflows/release.yaml
```

2. Cosign

You can install Cosign by aqua.

```sh
aqua g -i sigstore/cosign
```

From lintnet v1.0.1, the bundle file is `lintnet_<version>_checksums.txt.sigstore.json`.
Cosign v2.4.2 or later is required.

```sh
version=v1.0.1
checksum_file="lintnet_${version#v}_checksums.txt"
asset=lintnet_darwin_arm64.tar.gz
gh release download "$version" \
  -R lintnet/lintnet \
  -p "$asset" \
  -p "$checksum_file" \
  -p "${checksum_file}.sigstore.json"
cosign verify-blob \
  "$checksum_file" \
  --bundle "${checksum_file}.sigstore.json" \
  --certificate-identity-regexp 'https://github\.com/suzuki-shunsuke/go-release-workflow/\.github/workflows/release\.yaml@.*' \
  --certificate-oidc-issuer "https://token.actions.githubusercontent.com"
cat "$checksum_file" | sha256sum -c --ignore-missing -
```

For lintnet v1.0.0 or older, `*.pem` and `*.sig` are available instead.

```sh
version=v1.0.0
checksum_file="lintnet_${version#v}_checksums.txt"
asset=lintnet_darwin_arm64.tar.gz
gh release download "$version" \
  -R lintnet/lintnet \
  -p "$asset" \
  -p "$checksum_file" \
  -p "${checksum_file}.pem" \
  -p "${checksum_file}.sig"
cosign verify-blob \
  --signature "${checksum_file}.sig" \
  --certificate "${checksum_file}.pem" \
  --certificate-identity-regexp 'https://github\.com/suzuki-shunsuke/go-release-workflow/\.github/workflows/release\.yaml@.*' \
  --certificate-oidc-issuer "https://token.actions.githubusercontent.com" \
  "$checksum_file"
cat "$checksum_file" | sha256sum -c --ignore-missing -
```

3. Packslip

lintnet >= v1.0.1

You can install packslip by aqua.

```sh
aqua g -i jdx/packslip
```

```sh
version=v1.0.1
asset=lintnet_darwin_arm64.tar.gz
gh release download -R lintnet/lintnet "$version" -p "$asset" -p packslip.sigstore.json
packslip verify packslip.sigstore.json \
    --pin ps1_oy7uitmykcmnmavndmaorcvtta \
    --identity-prefix https://github.com/suzuki-shunsuke/go-release-workflow/.github/workflows/release.yaml@ \
    --artifact "$asset"
```

- `--pin` requires the release to be signed from the repository `lintnet/lintnet`, identified by its repository ID rather than its name
- `--identity-prefix` requires the release to be signed by the release workflow of [suzuki-shunsuke/go-release-workflow](https://github.com/suzuki-shunsuke/go-release-workflow)

On success, packslip prints a line starting with `ok:`.
Please make sure it shows the project and version you meant, e.g. `github.com/lintnet/lintnet 1.0.1`.

4. slsa-verifier

:::caution
As of lintnet v1.0.1, `multiple.intoto.jsonl` hasn't been published anymore.
:::

You can install slsa-verifier by aqua.

```sh
aqua g -i slsa-framework/slsa-verifier
```

```sh
version=v1.0.0
asset=lintnet_darwin_arm64.tar.gz
gh release download -R lintnet/lintnet "$version" -p "$asset" -p multiple.intoto.jsonl
slsa-verifier verify-artifact "$asset" \
  --provenance-path multiple.intoto.jsonl \
  --source-uri github.com/lintnet/lintnet \
  --source-tag "$version"
```

</details>

5. Go

```sh
go install github.com/lintnet/lintnet/cmd/lintnet@latest
```

## Shell completion

lintnet >= v0.4.7

lintnet supports shell completion for bash, zsh, and fish.

bash

```sh
source <(lintnet completion bash)
```

zsh

```sh
source <(lintnet completion zsh)
```

fish

```sh
lintnet completion fish > ~/.config/fish/completions/lintnet.fish
```
