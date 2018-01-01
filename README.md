## dotfiles

```bash
# 開発するのに必須なやつたち

# xcodeいれなきゃだめっぽいけどツールだけ入れたい時はこれ
# Xcode重いしめんどくさいから
xcode-select --install

# home brew install
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install.sh)"
```

```bash
brew install ghq
git config --global ghq.root ~/Repositories
ghq get git@github.com:fumihumi/dotfiles.git
```

```bash
ssh-keygen -t rsa -b 4096 -C '${fumihumi}@${コメント}'
pbcopy < github_id_rsa.pub
ssh -T git@github.com
# うまくできないときは ssh-key関連でうまくいってないのかも
# eval `ssh-agent`
# ssh-add -K github_id_rsa

# docker install
# docker toolboxのがいっぱい入ってるけどそんなに使わないし。↓で良さそう
brew cask install docker
brew install docker-compose
```

```bash
$ cd dotfiles
$ cp .gitconfig ~/.gitconfig
# 改めてghq経由でDotfilesRepoをクローンするため && dotfiles更新する場合これをする必要がある
$ cd setup
$ brew bundle
$ sh *.sh
```

```bash
# setup for ${any} languages

# using 'asdf'
# https://github.com/asdf-vm/asdf みてよしなに対応して。

# MEMO mysql。ローカルマシンに載せるのしんどいしdockerにしても良い気がしてきてrう
# mysql install
# at-markでversion指定しない方が良さそうな気もするけど、どうだろう。
brew install mysql@5.7
# brew link mysql@5.7
# brew linkコマンドうまくいかなかったぽい？
mysql -u root
echo 'export PATH="/usr/local/opt/mysql@5.7/bin:$PATH"'
which mysql
mysql.server start
mysql -u root

```

### MissionControl の設定いじる

> 「最新の使用状況に基づいて操作スペースを自動的に並び替える」のチェックを外す
> [qiita: Mission Control のデスクトップの順番が勝手に変わる件についての対処法](https://qiita.com/ayies128/items/f036ba7d89444b3b71f0)

キーボード > ショートカット > MissionControl >

- 左の操作スペースに移動 `alt + [`
- 右の操作スペースに移動 `alt + ]`

キーボード > 入力ソース

- ライブ変換 off


```
{ fumihumi } (.ssh): head config
Host *
  UseKeychain yes
  AddKeysToAgent yes

```
