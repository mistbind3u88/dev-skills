#!/usr/bin/env bash
set -euo pipefail

USAGE="Usage: mark.sh <type|--status|--clean>"
CHECK_TYPES=(lint test build review)

head_sha=$(git rev-parse HEAD)
head_short=$(git rev-parse --short HEAD)
branch=$(git rev-parse --abbrev-ref HEAD)

tag_prefix="mark/$branch"

status() {
  for type in "${CHECK_TYPES[@]}"; do
    tag="$tag_prefix/$type"
    tag_sha=$(git rev-parse --verify "$tag" 2>/dev/null || true)
    if [[ -z "$tag_sha" ]]; then
      printf "%-16s ✗ (未設置)\n" "$type"
    elif [[ "$tag_sha" == "$head_sha" ]]; then
      printf "%-16s ✓ (現在の HEAD)\n" "$type"
    else
      tag_short=$(git rev-parse --short "$tag_sha")
      behind=$(git rev-list --count "$tag_sha".."$head_sha" 2>/dev/null || echo "?")
      printf "%-16s ✗ (%s — %s commits behind)\n" "$type" "$tag_short" "$behind"
    fi
  done
}

clean() {
  tags=$(git tag -l "$tag_prefix/*")
  if [[ -z "$tags" ]]; then
    echo "$tag_prefix/ タグはありません"
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
  git tag -f "$tag_prefix/$type" HEAD
  echo "$tag_prefix/$type → $head_short"
}

if [[ $# -eq 0 ]]; then
  for t in "${CHECK_TYPES[@]}"; do
    mark "$t"
  done
  exit 0
fi

case "$1" in
  --status) status ;;
  --clean)  clean ;;
  *)
    for arg in "$@"; do
      mark "$arg"
    done
    ;;
esac
