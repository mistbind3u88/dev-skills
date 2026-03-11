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

declare -A latest_by_date=()
while IFS= read -r row; do
  date_part="${row%% *}"
  hash_part="${row##* }"
  if [[ -z "${latest_by_date[$date_part]+x}" ]]; then
    latest_by_date["$date_part"]="$hash_part"
  fi
done < <(git log --date=format:%Y-%m-%d --pretty=format:'%ad %H')

created_tags=()
retagged_tags=()
changed_tags=()

for date in "${!latest_by_date[@]}"; do
  tag_name="daily-$date"
  target_commit="${latest_by_date[$date]}"

  if git rev-parse -q --verify "refs/tags/$tag_name" >/dev/null; then
    current_commit="$(git rev-parse "$tag_name")"
    if [[ "$current_commit" != "$target_commit" ]]; then
      git tag -f "$tag_name" "$target_commit"
      retagged_tags+=("$tag_name:$current_commit:$target_commit")
      changed_tags+=("$tag_name")
    fi
  else
    git tag "$tag_name" "$target_commit"
    created_tags+=("$tag_name:$target_commit")
    changed_tags+=("$tag_name")
  fi
done

if [[ ${#changed_tags[@]} -eq 0 ]]; then
  echo "No tag changes required."
  exit 0
fi

if [[ ${#created_tags[@]} -gt 0 ]]; then
  IFS=$'\n' created_sorted=($(printf '%s\n' "${created_tags[@]}" | sort))
  unset IFS
  for entry in "${created_sorted[@]}"; do
    tag_name="${entry%%:*}"
    commit_hash="${entry##*:}"
    echo "Created $tag_name -> $commit_hash"
  done
fi

if [[ ${#retagged_tags[@]} -gt 0 ]]; then
  IFS=$'\n' retagged_sorted=($(printf '%s\n' "${retagged_tags[@]}" | sort))
  unset IFS
  for entry in "${retagged_sorted[@]}"; do
    tag_name="${entry%%:*}"
    rest="${entry#*:}"
    old_commit="${rest%%:*}"
    new_commit="${rest##*:}"
    echo "Retagged $tag_name: $old_commit -> $new_commit"
  done
fi

IFS=$'\n' changed_sorted=($(printf '%s\n' "${changed_tags[@]}" | sort -u))
unset IFS

push_specs=()
for tag_name in "${changed_sorted[@]}"; do
  push_specs+=("+refs/tags/$tag_name:refs/tags/$tag_name")
done

git push "$remote" "${push_specs[@]}"