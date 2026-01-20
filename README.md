## dotfiles

```bash
xcode-select --install

# home brew install
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install.sh)"
```

```bash
# setup for ${any} languages

# using 'mise'
# ref: https://mise.jdx.dev/getting-started.html
```

```bash
ssh-keygen -t rsa -b 4096 -C '${fumihumi}@${コメント}'
pbcopy < github_id_rsa.pub
ssh -T git@github.com
# troubleshooting ...
# eval `ssh-agent`
# ssh-add -K github_id_rsa
```

```bash
brew install ghq
git config --global ghq.root ~/Repositories
ghq get git@github.com:fumihumi/dotfiles.git

cd dotfiles
CURRENT_DIR=$(pwd)
ln -s $CURRENT_DIR/.gitconfig ~/.gitconfig
ln -s $CURRENT_DIR/.gitignore_global ~/.gitignore_global

$ cd setup
$ brew bundle
$ sh *.sh
```

```bash
# setup for ${any} languages

# using 'mise'
# https://github.com/jdx/mise みてよしなに対応して。
```

### MissionControl の設定いじる

> 「最新の使用状況に基づいて操作スペースを自動的に並び替える」のチェックを外す
> [ref: Mission Control のデスクトップの順番が勝手に変わる件についての対処法](https://qiita.com/ayies128/items/f036ba7d89444b3b71f0)

キーボード > ショートカット > MissionControl >

- 左の操作スペースに移動 `alt + [`
- 右の操作スペースに移動 `alt + ]`

キーボード > 入力ソース

- ライブ変換 off

## デスクトップとDock

-　> ホットコーナー

- 左上 > スクリーンセーバ
- 右上 > 通知センター
- 右上 > スリープ
- 右下 > Mission Control

### MissionControl

> 「最新の使用状況に基づいて操作スペースを自動的に並び替える」のチェックを外す
> [ref: Mission Control のデスクトップの順番が勝手に変わる件についての対処法](https://qiita.com/ayies128/items/f036ba7d89444b3b71f0)

```shell
# { fumihumi } (.ssh): head config
Host *
  UseKeychain yes
  AddKeysToAgent yes

```

## Tools

- BetterTouchTools
- Karabiner-Elements
- AltTab
- alfred

## Git

- `git cleanup-branches`: マージ済みのローカルブランチを一括削除（worktree で使用中のブランチは削除しません）
  - Dry-run: `git cleanup-branches -n`
  - 基準ブランチ（base）を指定: `git cleanup-branches -n --base <branch>`
  - 除外（正規表現）: `export GIT_IGNORE_BRANCH_LIST='^(main|master|develop)$'`
  - リポジトリごとの除外（推奨）: `git config --local cleanup.ignoreBranchRegex '^(main|master|develop)$'`
  - `--base` 省略時のデフォルト基準（リポジトリ単位で設定）
    - 現在ブランチ基準にする: `git config --local cleanup.baseStrategy current`
    - 特定ブランチ固定にする: `git config --local cleanup.baseRef develop`
