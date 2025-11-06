# git-wcp

Git worktree 間でファイルをコピーするツール

## 概要

`git-wcp` (worktree copy) は、異なる Git worktree 間でファイルを簡単にコピーするためのコマンドです。worktree 名の TAB 補完にも対応しており、`gwq` コマンドとの連携により効率的な操作が可能です。

## インストール

```bash
# git-tools ディレクトリで
./build.sh

# または個別にビルド
cd git-wcp
go build -o git-wcp
```

## 使用方法

```bash
git wcp <source> <destination>
```

### 形式

- `<worktree>:<path>` - 指定した worktree のファイル
- `@:<path>` - 現在の worktree のファイル

### 例

```bash
# feature-branch から main へファイルをコピー
git wcp feature-branch:src/config.json main:src/config.json

# 現在の worktree から test1 へ README をコピー
git wcp @:README.md test1:docs/README.md

# sample01 から現在の worktree へ設定ファイルをコピー
git wcp sample01:config.json @:config.json

# 詳細な出力を表示
git wcp -v feature-branch:src/main.go @:src/main.go
```

## 機能

### TAB 補完

Bash と Zsh で TAB 補完を有効にするには：

#### Bash
```bash
# 一時的に有効化
source <(git-wcp completion bash)

# 永続的に有効化 (Linux)
sudo git-wcp completion bash > /etc/bash_completion.d/git-wcp

# 永続的に有効化 (macOS with Homebrew)
git-wcp completion bash > $(brew --prefix)/etc/bash_completion.d/git-wcp
```

#### Zsh
```bash
# 一時的に有効化
source <(git-wcp completion zsh)

# 永続的に有効化
git-wcp completion zsh > "${fpath[1]}/_git-wcp"
```

### gwq との連携

`gwq` がインストールされている場合、`gwq list --json` を使用して worktree 情報を取得します。これにより、より正確で高速な worktree 検出が可能になります。

### フラグ

- `-v, --verbose`: 詳細な出力（コピー元とコピー先のフルパスを表示）

## 実装の特徴

- **Go 実装**: 高速で信頼性の高い動作
- **gwq 対応**: `gwq list --json` による効率的な worktree 検出
- **エラーハンドリング**: 詳細なエラーメッセージ
- **自動ディレクトリ作成**: コピー先のディレクトリが存在しない場合は自動作成

## 技術仕様

- 言語: Go 1.21+
- 依存関係: 
  - github.com/spf13/cobra (CLI フレームワーク)
  - gwq (オプション、worktree 管理ツール)