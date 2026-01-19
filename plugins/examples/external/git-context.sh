#!/usr/bin/env bash
set -euo pipefail

root="$(git rev-parse --show-toplevel 2>/dev/null || true)"
if [[ -z "$root" ]]; then
  exit 0
fi

branch="$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "detached")"

dirty="clean"
if ! git diff --quiet --ignore-submodules -- 2>/dev/null || ! git diff --cached --quiet --ignore-submodules -- 2>/dev/null; then
  dirty="dirty"
fi

ahead=0
behind=0
if git rev-parse --abbrev-ref --symbolic-full-name @{u} >/dev/null 2>&1; then
  counts="$(git rev-list --left-right --count @{upstream}...HEAD 2>/dev/null || echo "0\t0")"
  behind="$(printf "%s" "$counts" | awk '{print $1}')"
  ahead="$(printf "%s" "$counts" | awk '{print $2}')"
fi

value="Branch: ${branch} · ${dirty} · ↑${ahead} ↓${behind}\nRoot: ${root}"

json_escape() {
  printf '%s' "$1" | sed -e 's/\\/\\\\/g' -e 's/\"/\\\"/g'
}

printf '{\"label\":\"Git\",\"icon\":\"\",\"value\":\"%s\"}\n' "$(json_escape "$value")"
