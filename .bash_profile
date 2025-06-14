#Bashrc読み込み
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
eval "$(~/.local/bin/mise activate bash)"

# docker fzf
if [ -f '$(ghq root)/github.com/kwhrtsk/docker-fzf-completion/docker-fzf.bash' ]; then
  source $(ghq root)/github.com/kwhrtsk/docker-fzf-completion/docker-fzf.bash
fi

BREW_PREFIX=$(brew --prefix)

if [[ -f "$BREW_PREFIX/etc/profile.d/bash_completion.sh" ]]; then
  source "$BREW_PREFIX/etc/profile.d/bash_completion.sh"
fi

if [ -f ${BREW_PREFIX}/etc/profile.d/z.sh ]; then
  source $BREW_PREFIX/etc/profile.d/z.sh
fi

if [ -f ~/.bashrc ] ; then
  . ~/.bashrc
  for f in ~/.bash/*.rc; do source $f; done

  if [ -d ~/.bash/completion.d ]; then
    for c in ~/.bash/completion.d/*; do source "$c"; done
  fi

  for f in ~/.bash/works/*.rc; do source $f; done
fi

is_git_worktree() {
  if [ "$(git rev-parse --is-inside-work-tree 2>/dev/null)" = "true" ]; then
    toplevel=$(git rev-parse --show-toplevel 2>/dev/null)
    if [ -n "$toplevel" ] && [ -f "$toplevel/.git" ]; then
      echo " (worktree)"
    fi
  fi
}

git_prompt_string() {
  if type __git_ps1 &>/dev/null; then
    echo "$(__git_ps1)$(is_git_worktree)"
  fi
}

export PS1='\[\033[34m\]{ \u }\[\033[31m\]$(git_prompt_string)\[\033[33m\] [\t]\[\033[34m\] \n\[\033[32m\](\W)\[\033[00m\]: \[\033[00m\]'

if command -v rustc &> /dev/null; then
  export RUBY_CONFIGURE_OPTS="--enable-yjit"
fi
