---
sidebar_position: 200
---

# Supported data format

:::info
If you want to request new data format, please [create an issue](https://github.com/lintnet/lintnet/issues/new?assignees=&labels=enhancement%2Cnew-data-format&projects=&template=new-data-format-request.yml&title=New+Data+Format+Request%3A+%3Cformat+name%3E).
:::

lintnet can lint the following file formats.
lintnet judges file types by file extensions.
We're considering supporting additional file formats. [#37](https://github.com/lintnet/lintnet/issues/37)

format | `file_type` | file extensions | parser
--- | --- | --- | ---
CSV | csv | `.csv` | [encoding/csv](https://pkg.go.dev/encoding/csv#Reader)
HCL 2 | hcl2 | `.hcl` | [tmccombs/hcl2json](https://pkg.go.dev/github.com/tmccombs/hcl2json/convert)
JSON | json | `.json` | [encoding/json](https://pkg.go.dev/encoding/json#Decoder)
TOML | toml | `.toml` | [BurntSushi/toml](https://godocs.io/github.com/BurntSushi/toml#Decoder)
TSV | tsv | `.tsv` | [encoding/csv](https://pkg.go.dev/encoding/csv#Reader)
YAML | yaml | `.yml`, `.yaml` | [gopkg.in/yaml.v3](https://pkg.go.dev/gopkg.in/yaml.v3#Decoder)

## YAML is parsed as multiple documents

A YAML file is comprised of multiple documents separated by `---`.

e.g.

```yaml
---
# document 0
name: foo
---
# document 1
name: bar
```

So A YAML file is parsed as multiple documents even if the file includes only one document.
In lint rules, you need to iterate multiple documents or get the first document by specifying the index.

e.g.

```jsonnet
for doc in param.data.value # Iterate multiple documents
```

```jsonnet
param.data.value[0] # Get the first document
```

## Plain Text

lintnet judges file types by file extensions.
If no parser is found, lintnet parse the file as a plain text file.
The external variable `file_type` is `plain_text`.
The external variable `input` is empty, but you can still lint the file with other external variables such as `file_path` and `file_text`.

