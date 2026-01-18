#!/usr/bin/env bash
set -euo pipefail

plugin_name="${1:-}"
plugin_src="${2:-}"

if [[ -z "${plugin_name}" || -z "${plugin_src}" ]]; then
  echo "Usage: scripts/build-plugin.sh <name> <source-file>"
  echo "Example: scripts/build-plugin.sh hello plugins/examples/hello.go"
  exit 1
fi

goos="$(go env GOOS)"
goarch="$(go env GOARCH)"
out_dir="dist/plugins"
out_file="${out_dir}/${plugin_name}_${goos}_${goarch}.so"

mkdir -p "${out_dir}"
go build -buildmode=plugin -o "${out_file}" "${plugin_src}"

echo "Built ${out_file}"
echo "Upload this file to a GitHub Release and add it to plugins/manifest.json."
