#!/usr/bin/env bash
set -euo pipefail

USAGE="Usage: mark.sh <type|--status|--clean>"
CHECK_TYPES=(lint test build review)

head_sha=$(git rev-parse HEAD)
head_short=$(git rev-parse --short HEAD)

status() {
  for type in "${CHECK_TYPES[@]}"; do
    tag="check/$type"
    tag_sha=$(git rev-parse --verify "$tag" 2>/dev/null || true)
    if [[ -z "$tag_sha" ]]; then
      printf "%-16s ✗ (未設置)\n" "$tag"
    elif [[ "$tag_sha" == "$head_sha" ]]; then
      printf "%-16s ✓ (現在の HEAD)\n" "$tag"
    else
      tag_short=$(git rev-parse --short "$tag_sha")
      behind=$(git rev-list --count "$tag_sha".."$head_sha" 2>/dev/null || echo "?")
      printf "%-16s ✗ (%s — %s commits behind)\n" "$tag" "$tag_short" "$behind"
    fi
  done
}

clean() {
  tags=$(git tag -l 'check/*')
  if [[ -z "$tags" ]]; then
    echo "check/ タグはありません"
    return
  fi
  echo "$tags" | xargs git tag -d
}

mark() {
  local type="$1"
  local valid=false
  for t in "${CHECK_TYPES[@]}"; do
    [[ "$t" == "$type" ]] && valid=true && break
  done
  if ! $valid; then
    echo "Error: unknown type '$type'. Valid types: ${CHECK_TYPES[*]}" >&2
    exit 1
  fi
  git tag -f "check/$type" HEAD
  echo "check/$type → $head_short"
}

if [[ $# -eq 0 ]]; then
  echo "$USAGE" >&2
  exit 1
fi

case "$1" in
  --status) status ;;
  --clean)  clean ;;
  *)        mark "$1" ;;
esac
