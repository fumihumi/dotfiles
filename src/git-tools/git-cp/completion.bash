# Bash completion for git-cp
_git_cp() {
    local cur="${COMP_WORDS[COMP_CWORD]}"
    
    # Complete worktree names followed by colon
    if [[ "$cur" == *:* ]]; then
        # Already has worktree, complete file paths
        local prefix="${cur%%:*}:"
        local path="${cur#*:}"
        COMPREPLY=($(compgen -f -- "$path" | sed "s|^|$prefix|"))
    else
        # Complete worktree names
        local worktrees=$(git worktree list --porcelain | grep ^branch | cut -d' ' -f2 | sed 's|refs/heads/||')
        worktrees="@ $worktrees"
        COMPREPLY=($(compgen -W "$worktrees" -- "$cur"))
        # Add colon suffix for completed worktrees
        if [[ ${#COMPREPLY[@]} -eq 1 ]]; then
            COMPREPLY[0]="${COMPREPLY[0]}:"
        fi
    fi
}

complete -F _git_cp git-cp