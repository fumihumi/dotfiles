#!/usr/bin/env bash
set -euo pipefail

# git-cp: Copy files between git worktrees
# Usage: git cp <source> <destination>
# Format: <worktree>:<path> or @:<path> (@ means current worktree)

print_usage() {
    cat << EOF
Usage: git cp <source> <destination>

Copy files between git worktrees.

Format:
  <worktree>:<path>  - File in specified worktree
  @:<path>           - File in current worktree

Examples:
  git cp worktreeA:src/file.txt worktreeB:dest/file.txt
  git cp @:README.md test1:docs/README.md
  git cp sample01:config.json @:config.json
EOF
}

error() {
    echo "Error: $1" >&2
    exit 1
}

# Parse worktree:path format
parse_location() {
    local location="$1"

    if [[ ! "$location" =~ ^([^:]+):(.+)$ ]]; then
        error "Invalid format: $location. Expected format: <worktree>:<path>"
    fi

    local worktree="${BASH_REMATCH[1]}"
    local path="${BASH_REMATCH[2]}"

    echo "$worktree"
    echo "$path"
}

# Get worktree root directory by name
get_worktree_root() {
    local worktree_name="$1"

    # Handle @ as current worktree
    if [[ "$worktree_name" == "@" ]]; then
        git rev-parse --show-toplevel 2>/dev/null || error "Not in a git repository"
        return
    fi

    # Get worktree info in porcelain format
    local worktree_info
    worktree_info=$(git worktree list --porcelain 2>/dev/null) || error "Failed to list worktrees"

    # Parse worktree info to find matching worktree
    local current_worktree=""
    local current_branch=""
    local found=false

    while IFS= read -r line; do
        if [[ "$line" =~ ^worktree[[:space:]](.+)$ ]]; then
            current_worktree="${BASH_REMATCH[1]}"
        elif [[ "$line" =~ ^branch[[:space:]]refs/heads/(.+)$ ]]; then
            current_branch="${BASH_REMATCH[1]}"
            if [[ "$current_branch" == "$worktree_name" ]]; then
                echo "$current_worktree"
                found=true
                break
            fi
        fi
    done <<< "$worktree_info"

    if [[ "$found" != "true" ]]; then
        error "Worktree '$worktree_name' not found"
    fi
}

# Main function
main() {
    if [[ $# -ne 2 ]]; then
        print_usage
        exit 1
    fi

    local source="$1"
    local destination="$2"

    # Parse source
    local source_parts
    IFS=$'\n' read -d '' -ra source_parts < <(parse_location "$source" && printf '\0')
    local source_worktree="${source_parts[0]}"
    local source_path="${source_parts[1]}"

    # Parse destination
    local dest_parts
    IFS=$'\n' read -d '' -ra dest_parts < <(parse_location "$destination" && printf '\0')
    local dest_worktree="${dest_parts[0]}"
    local dest_path="${dest_parts[1]}"

    # Get worktree roots
    local source_root
    source_root=$(get_worktree_root "$source_worktree")
    local dest_root
    dest_root=$(get_worktree_root "$dest_worktree")

    # Build full paths
    local source_full="$source_root/$source_path"
    local dest_full="$dest_root/$dest_path"

    # Check if source exists
    if [[ ! -e "$source_full" ]]; then
        error "Source file not found: $source_full"
    fi

    # Create destination directory if needed
    local dest_dir
    dest_dir=$(dirname "$dest_full")
    if [[ ! -d "$dest_dir" ]]; then
        mkdir -p "$dest_dir" || error "Failed to create directory: $dest_dir"
    fi

    # Copy file
    cp "$source_full" "$dest_full" || error "Failed to copy file"

    echo "Copied: $source -> $destination"
    echo "  From: $source_full"
    echo "  To:   $dest_full"
}

# Run main function
main "$@"
