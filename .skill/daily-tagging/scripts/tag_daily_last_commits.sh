#!/usr/bin/env bash
set -euo pipefail

remote="${1:-origin}"
branch="$(git rev-parse --abbrev-ref HEAD)"

if ! git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  echo "Not a git repository" >&2
  exit 1
fi

# Pull latest changes first (fast-forward only for safety).
git pull --ff-only "$remote" "$branch"

# git log outputs newest first, so the first commit per date is that day's latest.
# Output: "YYYY-MM-DD <hash>" with one entry per date.
date_commit_pairs="$(git log --date=format:%Y-%m-%d --pretty=format:'%ad %H' | tr -d '\r' | awk '!seen[$1]++ { print }')"

changed_tags=""
output_created=""
output_retagged=""

while IFS=' ' read -r date_part target_commit; do
  [ -z "$date_part" ] && continue
  tag_name="daily-$date_part"

  if git rev-parse -q --verify "refs/tags/$tag_name" >/dev/null; then
    current_commit="$(git rev-parse "$tag_name" | tr -d '\r')"
    if [ "$current_commit" != "$target_commit" ]; then
      git tag -f "$tag_name" "$target_commit"
      output_retagged="${output_retagged}Retagged $tag_name: $current_commit -> $target_commit
"
      changed_tags="${changed_tags}${tag_name}
"
    fi
  else
    git tag "$tag_name" "$target_commit"
    output_created="${output_created}Created $tag_name -> $target_commit
"
    changed_tags="${changed_tags}${tag_name}
"
  fi
done <<< "$date_commit_pairs"

if [ -z "$changed_tags" ]; then
  echo "No tag changes required."
  exit 0
fi

# Print results sorted
if [ -n "$output_created" ]; then
  printf '%s' "$output_created" | sort
fi
if [ -n "$output_retagged" ]; then
  printf '%s' "$output_retagged" | sort
fi

# Push all changed tags in one command
printf '%s' "$changed_tags" | tr -d '\r' | sed '/^$/d' | sed 's|.*|+refs/tags/&:refs/tags/&|' | xargs git push "$remote"
