#!/usr/bin/env bash

set -euo pipefail

for file in $(find json-schema/jsonnet -name '*.jsonnet'); do
  base=$(basename "$file")
  jsonnet "$file" > "json-schema/${base/jsonnet/json}"
done
