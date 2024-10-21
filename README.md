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

cd setup
brew bundle
sh *.sh
```

### MissionControl の設定いじる

> 「最新の使用状況に基づいて操作スペースを自動的に並び替える」のチェックを外す
> [ref: Mission Control のデスクトップの順番が勝手に変わる件についての対処法](https://qiita.com/ayies128/items/f036ba7d89444b3b71f0)

キーボード > ショートカット > MissionControl >

- 左の操作スペースに移動 `alt + [`
- 右の操作スペースに移動 `alt + ]`

キーボード > 入力ソース

- ライブ変換 off

```shell
# { fumihumi } (.ssh): head config
Host *
  UseKeychain yes
  AddKeysToAgent yes
```
