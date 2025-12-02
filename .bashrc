# .bashrc

# オリジナルのTERM=xtermはカラー表示できないと思われる。
if [ "$TERM" == xterm ]; then
    export TERM=xterm-color
fi

PATH="$PATH:~/.mos/bin"

export EDITOR=nvim
export TIG_EDITOR=nvim
export GIT_EDITOR=nvim
export PROMPT_COMMAND='history -a; history -r'
export LANG=ja_JP.UTF-8
export LSCOLORS=gxfxcxdxbxegedabagacad
export GREP_COLOR='1;37;41'
export GIT_PS1_SHOWDIRTYSTATE=true

. "$HOME/.cargo/env"

### MANAGED BY RANCHER DESKTOP START (DO NOT EDIT)
export PATH="/Users/takafumi.suzuki/.rd/bin:$PATH"
### MANAGED BY RANCHER DESKTOP END (DO NOT EDIT)
alias claude="/Users/takafumi.suzuki/.claude/local/claude"

# gw shell integration
eval "$(gw shell-integration --show-script --shell=bash)"

[[ "$TERM_PROGRAM" == "kiro" ]] && . "$(kiro --locate-shell-integration-path bash)"
