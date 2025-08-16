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

1. [GitHub CLI](https://cli.github.com/)
1. [slsa-verifier](https://github.com/slsa-framework/slsa-verifier)
1. [Cosign](https://github.com/sigstore/cosign)

--

1. GitHub CLI

lintnet >= v0.4.8

You can install GitHub CLI by aqua.

```sh
aqua g -i cli/cli
```

```sh
gh release download -R lintnet/lintnet v0.4.8 -p lintnet_darwin_arm64.tar.gz
gh attestation verify lintnet_darwin_arm64.tar.gz \
  -R lintnet/lintnet \
  --signer-workflow suzuki-shunsuke/go-release-workflow/.github/workflows/release.yaml
```

2. slsa-verifier

You can install slsa-verifier by aqua.

```sh
aqua g -i slsa-framework/slsa-verifier
```

```sh
gh release download -R lintnet/lintnet v0.4.8 -p lintnet_darwin_arm64.tar.gz  -p multiple.intoto.jsonl
slsa-verifier verify-artifact lintnet_darwin_arm64.tar.gz \
  --provenance-path multiple.intoto.jsonl \
  --source-uri github.com/lintnet/lintnet \
  --source-tag v0.4.8
```

3. Cosign

You can install Cosign by aqua.

```sh
aqua g -i sigstore/cosign
```

```sh
gh release download -R lintnet/lintnet v0.4.8
cosign verify-blob \
  --signature lintnet_0.4.8_checksums.txt.sig \
  --certificate lintnet_0.4.8_checksums.txt.pem \
  --certificate-identity-regexp 'https://github\.com/suzuki-shunsuke/go-release-workflow/\.github/workflows/release\.yaml@.*' \
  --certificate-oidc-issuer "https://token.actions.githubusercontent.com" \
  lintnet_0.4.8_checksums.txt
```

After verifying the checksum, verify the artifact.

```sh
cat lintnet_0.4.8_checksums.txt | sha256sum -c --ignore-missing
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
