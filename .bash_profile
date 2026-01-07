eval $(/opt/homebrew/bin/brew shellenv)
eval "$(mise activate bash)"

if [ -f ~/.bashrc ] ; then
  . ~/.bashrc
fi


GHQ_ROOT=$(ghq root)

BREW_PREFIX=$(brew --prefix)


if [ -d "$BREW_PREFIX/etc/profile.d" ]; then
  for f in "$BREW_PREFIX"/etc/profile.d/*.sh; do
    if [ -f "$f" ]; then
      source "$f"
    fi
  done
fi

if [ -d "$BREW_PREFIX/etc/bash_completion.d" ]; then
  for f in "$BREW_PREFIX"/etc/bash_completion.d/*.sh; do
    if [ -f "$f" ]; then
      source "$f"
    fi
  done
fi

# NOTE: after bash_completion.sh
for f in ~/.bash/*.rc; do source $f; done

if [ -d ~/.bash/completion.d ]; then
  for c in ~/.bash/completion.d/*; do source "$c"; done
fi

for f in ~/.bash/works/*.rc; do source $f; done

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

### MANAGED BY RANCHER DESKTOP START (DO NOT EDIT)
export PATH="/Users/takafumi.suzuki/.rd/bin:$PATH"
### MANAGED BY RANCHER DESKTOP END (DO NOT EDIT)
export PATH="/opt/homebrew/opt/curl/bin:$PATH"
export PATH="/opt/homebrew/opt/libpq/bin:$PATH"

# .bash/bin フォルダのスクリプトをコマンドとして利用できるようにPATHに追加
export GIT_BIN_PATH="$GHQ_ROOT/github.com/fumihumi/dotfiles/.bash/bin"
export GIT_WORK_BIN_PATH="$GHQ_ROOT/github.com/fumihumi/dotfiles/.bash/works/bin"

export PATH="$GIT_BIN_PATH:$PATH"
export PATH="$GIT_WORK_BIN_PATH:$PATH"
