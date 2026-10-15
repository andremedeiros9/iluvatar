#!/usr/bin/env bash
set -euo pipefail

# Determines the next release version and creates/pushes a release branch
# from origin/main. Intended to be called from .github/workflows/release-branch.yml.
#
# Usage:
#   release-branch.sh determine-version   # writes new_version to $GITHUB_OUTPUT
#   release-branch.sh create-branch <version>
#   release-branch.sh set-default-branch <version>

command="${1:-}"

case "$command" in
  determine-version)
    # Fetch all remote branches so we can see existing release/v* branches
    git fetch origin

    # List existing release branches matching release/vX.Y.Z
    existing=$(git branch -r | grep 'origin/release/v' | sed 's/.*origin\/release\/v//' | sort -V || true)

    # Start from 0.0.0 if none exist
    major=0
    minor=0
    patch=0

    if [ -n "$existing" ]; then
      # Take the highest existing version
      latest=$(echo "$existing" | tail -n 1)
      major=$(echo "$latest" | cut -d. -f1)
      minor=$(echo "$latest" | cut -d. -f2)
      patch=$(echo "$latest" | cut -d. -f3)
    fi

    # Increment patch
    patch=$((patch + 1))

    # Carry over when patch exceeds 9
    if [ "$patch" -gt 9 ]; then
      patch=0
      minor=$((minor + 1))
    fi

    # Carry over when minor exceeds 9 (only major can exceed 9)
    if [ "$minor" -gt 9 ]; then
      minor=0
      major=$((major + 1))
    fi

    new_version="v${major}.${minor}.${patch}"
    echo "new_version=${new_version}" >> "$GITHUB_OUTPUT"
    echo "Next release version: ${new_version}"
    ;;

  create-branch)
    version="${2:?version required}"
    branch_name="release/${version}"
    git checkout -b "${branch_name}" origin/main
    git push -u origin "${branch_name}"
    echo "Created and pushed branch: ${branch_name}"
    ;;

  set-default-branch)
    version="${2:?version required}"
    branch_name="release/${version}"
    gh api -X PATCH "repos/${GITHUB_REPOSITORY}" \
      -f default_branch="${branch_name}" >/dev/null
    echo "Default branch set to: ${branch_name}"
    ;;

  *)
    echo "Unknown command: $command" >&2
    echo "Usage: $0 {determine-version|create-branch <version>}" >&2
    exit 1
    ;;
esac
