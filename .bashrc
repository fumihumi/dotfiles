# .bashrc

# オリジナルのTERM=xtermはカラー表示できないと思われる。
if [ "$TERM" == xterm ]; then
    export TERM=xterm-color
fi

PATH="$PATH:~/.mos/bin"

# .bash/bin フォルダのスクリプトをコマンドとして利用できるようにPATHに追加
export PATH="$HOME/Repositories/github.com/fumihumi/dotfiles/.bash/bin:$PATH"

export EDITOR=nvim
export TIG_EDITOR=nvim
export GIT_EDITOR=nvim
export PROMPT_COMMAND='history -a; history -r'
export LANG=ja_JP.UTF-8
export LSCOLORS=gxfxcxdxbxegedabagacad
export GREP_COLOR='1;37;41'
export GIT_PS1_SHOWDIRTYSTATE=true

. "$HOME/.cargo/env"
