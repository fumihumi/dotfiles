# git-tools

Git サブコマンドコレクション - 日常的な Git 操作を効率化するツール群

## 概要

このディレクトリには、Git の操作を拡張・効率化するためのカスタムサブコマンドが含まれています。規約ベースの自動ビルドシステムにより、新しいツールの追加も簡単です。

## ツール一覧

### git-diff-summary

ブランチ間の差分を最終更新者情報と共に表示

```bash
git diff-summary origin/main origin/develop
```

### git-wcp

Git worktree 間でファイルをコピー

```bash
git wcp feature-branch:src/config.json @:src/config.json
```

### git-wcd

Git worktree のパスを取得・移動

```bash
cd $(git wcd feature-branch)
```

## インストール

```bash
# すべてのツールをビルド
./build.sh

# ビルド状態を確認
./build.sh check

# Bash補完をインストール
./build.sh completion install
```

## ビルドシステム

### 規約ベースの自動検出

- `go.mod` がある → Go プロジェクトとしてビルド
- `<dirname>.sh` がある → Bash スクリプトとしてインストール
- `completion.bash` がある → 補完スクリプトをロード

### コマンド

```bash
./build.sh              # 全ツールをビルド（デフォルト）
./build.sh check        # ビルド状態を確認
./build.sh completion show     # 補完設定を表示
./build.sh completion install  # 補完を ~/.bashrc に追加
./build.sh clean        # ビルド済みバイナリを削除
./build.sh help         # ヘルプを表示
```

## 新しいツールの追加

### Go製ツールの場合

```bash
# 1. ディレクトリを作成
mkdir git-newtool

# 2. Go モジュールを初期化
cd git-newtool
go mod init git-newtool

# 3. main.go を作成
# 4. build.sh を実行すると自動的にビルドされる
```

### Bashスクリプトの場合

```bash
# 1. ディレクトリを作成
mkdir git-newtool

# 2. スクリプトを作成（ディレクトリ名と同じ名前）
touch git-newtool/git-newtool.sh
chmod +x git-newtool/git-newtool.sh

# 3. 補完スクリプトを追加（オプション）
touch git-newtool/completion.bash
```

## ディレクトリ構造

```
git-tools/
├── build.sh              # ビルドスクリプト
├── build.go              # ビルドシステム本体
├── git-diff-summary/     # Go製ツール例
│   ├── go.mod
│   ├── main.go
│   └── README.md
├── git-wcp/              # Go製ツール（worktree間コピー）
│   ├── go.mod
│   ├── main.go
│   └── main_test.go
└── git-wcd/
    ├── git-wcd.sh
    └── completion.bash
```

## 必要な環境

- Go 1.21以上（Go製ツールをビルドする場合）
- Bash 4.0以上
- Git 2.0以上
