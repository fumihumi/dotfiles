#!/usr/bin/env bash
set -euo pipefail

# git-where: Display or find git worktree paths
# Usage: git where [options] [worktree]
# Options:
#   --json    Output in JSON format
#   --fzf     Use fzf for interactive selection
# Without arguments: List all worktrees
# With worktree name: Display specific worktree

print_usage() {
    cat << EOF
Usage: git where [options] [worktree]

Options:
  --json    Output in JSON format
  --fzf     Use fzf for interactive selection
  -h, --help    Show this help message

Examples:
  git where                  # List all worktrees
  git where --fzf            # Interactive selection with fzf
  git where feature/auth     # Show specific worktree info
  git where --json           # List all worktrees in JSON
  git where --json main      # Show specific worktree in JSON
EOF
}

error() {
    echo "Error: $1" >&2
    exit 1
}

# Parse command line arguments
parse_args() {
    local args=()
    JSON_OUTPUT=false
    FZF_MODE=false

    while [[ $# -gt 0 ]]; do
        case "$1" in
            --json)
                JSON_OUTPUT=true
                shift
                ;;
            --fzf)
                FZF_MODE=true
                shift
                ;;
            -h|--help)
                print_usage
                exit 0
                ;;
            --)
                shift
                args+=("$@")
                break
                ;;
            -*)
                error "Unknown option: $1"
                ;;
            *)
                args+=("$1")
                shift
                ;;
        esac
    done

    ARGS=("${args[@]}")
}

# Display single worktree info
show_worktree() {
    local worktree="$1"
    local worktree_info=$(gwq get "$worktree" --json 2>/dev/null)

    if [[ -z "$worktree_info" || "$worktree_info" == "null" ]]; then
        # Try partial match
        local all_worktrees=$(gwq list --json)
        local matches=$(echo "$all_worktrees" | jq -r --arg wt "$worktree" '.[] | select(.branch | contains($wt))')
        local match_count=$(echo "$matches" | jq -s 'length')

        if [[ $match_count -eq 0 ]]; then
            error "Worktree '$worktree' not found"
        elif [[ $match_count -eq 1 ]]; then
            # Single match found
            worktree_info="$matches"
        else
            # Multiple matches - show in table format
            if [[ "$JSON_OUTPUT" == "true" ]]; then
                # For JSON output, return array of matches
                echo "$matches" | jq -s '.'
            else
                echo "Multiple worktrees found matching '$worktree':"
                echo "$matches" | jq -r '"\(.branch)\t\(.path)"' | column -t -s $'\t' | sed 's/^/  /'
            fi
            exit 1
        fi
    fi

    if [[ "$JSON_OUTPUT" == "true" ]]; then
        echo "$worktree_info"
    else
        local path=$(echo "$worktree_info" | jq -r '.path')
        local branch=$(echo "$worktree_info" | jq -r '.branch')

        # Format as table to match list_all output
        printf "%-30s %s\n" "$branch" "$path"
    fi
}

# List all worktrees
list_all() {
    if [[ "$JSON_OUTPUT" == "true" ]]; then
        gwq list --json
    else
        local worktrees=$(gwq list --json | jq -r '.[] | "\(.branch)\t\(.path)"')
        if [[ -z "$worktrees" ]]; then
            echo "No worktrees found"
        else
            echo "$worktrees" | column -t -s $'\t'
        fi
    fi
}

# Interactive selection with fzf
select_with_fzf() {
    # If specific worktree is given, filter the list
    local filter_arg=""
    if [[ ${#ARGS[@]} -gt 0 ]]; then
        filter_arg="${ARGS[0]}"
    fi

    local worktrees=$(gwq list --json)

    # Apply filter if provided
    if [[ -n "$filter_arg" ]]; then
        worktrees=$(echo "$worktrees" | jq --arg filter "$filter_arg" '[.[] | select(.branch | contains($filter))]')
    fi

    local selected=$(echo "$worktrees" | \
        jq -r '.[] | "\(.branch)\t\(.path)"' | \
        column -t -s $'\t' | \
        fzf --prompt="Select worktree: " --preview-window=hidden)

    if [[ -z "$selected" ]]; then
        error "No worktree selected"
    fi

    # Extract branch name from the selected line
    local branch=$(echo "$selected" | awk '{print $1}')

    if [[ "$JSON_OUTPUT" == "true" ]]; then
        # Get the full JSON info for the selected branch
        echo "$worktrees" | jq --arg branch "$branch" '.[] | select(.branch == $branch)'
    else
        echo "$selected"
    fi
}

main() {
    parse_args "$@"

    if [[ "$FZF_MODE" == "true" ]]; then
        # FZF mode
        select_with_fzf
    elif [[ ${#ARGS[@]} -eq 0 ]]; then
        # No arguments - list all
        list_all
    else
        # Specific worktree
        show_worktree "${ARGS[0]}"
    fi
}

main "$@"
