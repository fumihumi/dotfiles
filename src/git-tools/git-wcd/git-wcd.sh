#!/usr/bin/env bash
set -euo pipefail

# git-wcd: Get path to a git worktree
# Usage: git wcd <worktree>
# Format: <worktree> (e.g., worktreeA)

print_usage() {
    cat << EOF
Usage: git wcd [worktree]
  Without arguments: Interactive selection with fzf
  With worktree name: Direct path output
EOF
}

error() {
    echo "Error: $1" >&2
    exit 1
}

# gwq command is "ghq" like git worktree commandline
main() {
  local worktree_path=""

  if [[ $# -gt 0 ]]; then
    # Direct worktree specification
    local worktree="$1"
    worktree_path=$(gwq get "$worktree" --json | jq -r '.path')
  else
    # Interactive selection with fzf
    worktree_path=$(gwq list --json | jq -c '.[] | { branch: .branch, path: .path}' | fzf --prompt="Select worktree: " | jq -rc '.path')
  fi

  if [[ -z "$worktree_path" || "$worktree_path" == "null" ]]; then
    error "No worktree selected or found"
  fi

  # Output only the path
  echo "$worktree_path"
}

main "$@"
