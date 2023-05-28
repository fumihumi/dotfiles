#bashrc読み込み
eval $(/opt/homebrew/bin/brew shellenv)

# The next line updates PATH for the Google Cloud SDK.
if [ -f '/Users/fumihumi/google-cloud-sdk/path.bash.inc' ]; then
  source '/Users/fumihumi/google-cloud-sdk/path.bash.inc';
fi

# The next line enables shell command completion for gcloud.
if [ -f '/Users/fumihumi/google-cloud-sdk/completion.bash.inc' ]; then
  source '/Users/fumihumi/google-cloud-sdk/completion.bash.inc';
fi

if [ -f '~/.bash/completion/exercism_completion.bash' ]; then
  source '~/.bash/completion/exercism_completion.bash';
fi

# docker fzf
if [ -f '$(ghq root)/github.com/kwhrtsk/docker-fzf-completion/docker-fzf.bash' ]; then
  source $(ghq root)/github.com/kwhrtsk/docker-fzf-completion/docker-fzf.bash
fi

BREW_PREFIX=$(brew --prefix)

if [ -f ${BREW_PREFIX}/etc/bash_completion ]; then
    source ${BREW_PREFIX}/etc/bash_completion
fi

#CHECK https://qiita.com/koyopro/items/3fce94537df2be6247a3
if [ -f ${BREW_PREFIX}/etc/bash_completion.d ]; then
    source "${BREW_PREFIX}/etc/bash_completion.d/git-prompt.sh"
    source "${BREW_PREFIX}/etc/bash_completion.d/git-completion.bash"
fi


if [ -f ${BREW_PREFIX}/etc/profile.d/z.sh ]; then
  source ${BREW_PREFIX}/etc/profile.d/z.sh
fi

if [ -f $(which asdf) ]; then
  source "${BREW_PREFIX}/opt/asdf/libexec/asdf.sh"
  source "${BREW_PREFIX}/opt/asdf/etc/bash_completion.d/asdf.bash"
fi

if [ -f '~/anaconda3/etc/profile.d/conda.sh' ]; then
  source ~/anaconda3/etc/profile.d/conda.sh
fi

if [ -f ~/.bashrc ] ; then
  . ~/.bashrc
  for f in ~/.bash/*.rc; do source $f; done
  for f in ~/.bash/works/*.rc; do source $f; done
fi

export PS1='\[\033[34m\]{ \u }\[\033[31m\]$(__git_ps1)\[\033[33m\] [\t]\[\033[34m\] \n\[\033[32m\](\W)\[\033[00m\]: \[\033[00m\]'
