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
  # .claude はランタイムデータ(projects, plugins, cache 等)が大量に書き込まれる
  # ディレクトリなので、まるごとリンクせず管理対象の設定ファイルだけを個別にリンクする
  [[ ${f} = ".claude" ]] && continue
  ln -snfv ${DOT_DIRECTORY}/${f} ${HOME}/${f}
done

ln -snfv ~/Repositories/github.com/fumihumi/dotfiles/.config/ ~

# Claude Code の設定ファイルのみを ~/.claude 配下にリンクする
mkdir -p ${HOME}/.claude
for f in CLAUDE.md settings.json statusline-command.sh
do
  ln -snfv ${DOT_DIRECTORY}/.claude/${f} ${HOME}/.claude/${f}
done

echo 'liked dotfiles complete!'
