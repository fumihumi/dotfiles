#!/bin/sh

set -u
GITHUB_USERNAME="fumihumi"
REPOSITORY_NAME="dotfiles"
# DOT_DIRECTORY="$(ghq root)/${GITHUB_USERNAME}/${REPOSITORY_NAME}"
DOT_DIRECTORY="${HOME}/Repositories/github.com/${GITHUB_USERNAME}/${REPOSITORY_NAME}"

cd ${DOT_DIRECTORY}

GITCONFIG_FILE=".gitconfig"
ln -snfv ${DOT_DIRECTORY}/${GITCONFIG_FILE} ${HOME}/${GITCONFIG_FILE}

for f in .??*
do
  # 無視したいファイル
  [[ ${f} = ".git" ]] && continue
  [[ ${f} = ".gitignore" ]] && continue
  ln -snfv ${DOT_DIRECTORY}/${f} ${HOME}/${f}
done

ln -snfv ~/Repositories/github.com/fumihumi/dotfiles/.config/ ~/.config

echo 'liked dotfiles complete!'
