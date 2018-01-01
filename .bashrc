# .bashrc

#export PS1='\[\033[34m\]{ \u }\[\033[31m\]$(__git_ps1)\[\033[33m\] [\t]\[\033[34m\] \n\[\033[32m\]( \W )\[\033[00m\]:\[\033[00m\]'

# { fumihumi } (master) [15:02:06]
# ( ~ ):

# export PS1='\[\033[34m\]{ \u }\[\033[31m\]$(__git_ps1)\[\033[33m\] [\t]\[\033[34m\] \[\033[32m\]( \W )\[\033[00m\]:\[\033[00m\]'
# { fumihumi } (dev *) [01:22:19] ( dotfiles ):

export PS1='\[\033[34m\]{ \u }\[\033[31m\]$(__git_ps1)\[\033[32m\] (\W)\[\033[32m\]: \[\033[00m\]'
# { fumihumi } (dev *) (dotfiles):

# オリジナルのTERM=xtermはカラー表示できないと思われる。
if [ "$TERM" == xterm ]; then
    export TERM=xterm-color
fi

GIT_PS1_SHOWDIRTYSTATE=true

PATH="$PATH:/Users/fumihumi/.mos/bin"

export EDITOR="vim"
export PROMPT_COMMAND='history -a; history -r'
export LANG=ja_JP.UTF-8
export LSCOLORS=gxfxcxdxbxegedabagacad
export GREP_COLOR='1;37;41'
