#!/usr/bin/env bash
set -euo pipefail

# Create a backup branch for the current HEAD.
# Naming: <branch>-<9-digit-SHA>
# If the branch name already ends with a 9-digit hex suffix, strip it first
# to avoid stacking (e.g. feature-abc123456-def789012).

branch="$(git rev-parse --abbrev-ref HEAD)"
short_sha="$(git rev-parse --short=9 HEAD)"

# Strip existing backup suffix (hyphen + exactly 9 hex chars at end)
base="$(printf '%s' "$branch" | sed 's/-[0-9a-f]\{9\}$//')"

backup_name="${base}-${short_sha}"

git branch -f "$backup_name"
echo "$backup_name"
