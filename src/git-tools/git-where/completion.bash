#!/usr/bin/env bash

# Bash completion for git-where
_git_where() {
    local cur="${COMP_WORDS[COMP_CWORD]}"
    local prev="${COMP_WORDS[COMP_CWORD-1]}"
    
    # Check if we're completing an option
    if [[ "$cur" == -* ]]; then
        # Offer options
        local options="--json --fzf --help -h"
        COMPREPLY=($(compgen -W "$options" -- "$cur"))
        return
    fi
    
    # Check if previous word was an option that doesn't take arguments
    case "$prev" in
        --json|--fzf|--help|-h)
            # After these options, we can still complete worktree names
            ;;
        git-where|where)
            # First argument after command (handle both 'git-where' and 'git where')
            ;;
        *)
            # Check if we're in a git command context
            if [[ "${COMP_WORDS[0]}" == "git" && "${COMP_WORDS[1]}" == "where" ]]; then
                # We're in 'git where' context, continue with completion
                :
            else
                # Already have a worktree argument, no more completion
                return
            fi
            ;;
    esac
    
    # Get list of worktrees using gwq
    if command -v gwq >/dev/null 2>&1; then
        local worktrees=$(gwq list --json 2>/dev/null | jq -r '.[].branch' 2>/dev/null)
        COMPREPLY=($(compgen -W "$worktrees" -- "$cur"))
    fi
}

# Register completion for both git-where and git where
complete -F _git_where git-where

# Also register for git subcommand
_git_where_wrapper() {
    # When called as 'git where', adjust COMP_WORDS
    _git_where
}

# Register git where completion if git completion is loaded
if declare -F __git_complete >/dev/null 2>&1; then
    __git_complete git-where _git_where
fi
