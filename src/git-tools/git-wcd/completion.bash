# Bash completion for git-wcd
_git_wcd() {
    local cur="${COMP_WORDS[COMP_CWORD]}"
    
    # Complete worktree names
    local worktrees=$(git worktree list --porcelain 2>/dev/null | grep ^branch | cut -d' ' -f2 | sed 's|refs/heads/||')
    COMPREPLY=($(compgen -W "$worktrees" -- "$cur"))
}

complete -F _git_wcd git-wcd