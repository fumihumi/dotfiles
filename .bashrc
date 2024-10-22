# .bashrc

# オリジナルのTERM=xtermはカラー表示できないと思われる。
if [ "$TERM" == xterm ]; then
    export TERM=xterm-color
fi

PATH="$PATH:/Users/fumihumi/.mos/bin"

export EDITOR="vim"
export PROMPT_COMMAND='history -a; history -r'
export LANG=ja_JP.UTF-8
export LSCOLORS=gxfxcxdxbxegedabagacad
export GREP_COLOR='1;37;41'
export GIT_PS1_SHOWDIRTYSTATE=true

. "$HOME/.cargo/env"
