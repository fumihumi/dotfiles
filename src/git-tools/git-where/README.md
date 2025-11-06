# git-where

Git worktree の場所を確認・表示するためのツール

## 概要

`git-where` は Git worktree のパスを簡単に確認できるツールです。デフォルトでは全 worktree の一覧を表示し、様々なオプションで出力形式を制御できます。

## 必要な依存関係

- `gwq` - Git worktree を管理するツール
- `fzf` - インタラクティブな選択UI（--fzf オプション使用時）
- `jq` - JSON パーサー

## インストール

```bash
# 1. このリポジトリをクローン
git clone https://github.com/fumihumi/dotfiles.git

# 2. build.sh を実行してシンボリックリンクを作成
cd dotfiles/src/git-tools
./build.sh

# 3. bash completion を有効化 (.bashrc に追加)
source ~/.bash/completion.d/git-where.rc
```

## 使い方

### 基本的な使い方

```bash
# 全 worktree を一覧表示
git where

# 特定の worktree 情報を表示
git where <worktree-name>

# 部分一致で検索（一意に特定できる場合）
git where feat  # feature/auth にマッチ
```

### オプション

- `--json` - JSON 形式で出力
- `--fzf` - fzf を使用したインタラクティブ選択
- `-h, --help` - ヘルプを表示

### 使用例

```bash
# 全 worktree の一覧表示
$ git where
main                              /Users/username/projects/repo
feature/auth                      /Users/username/projects/repo-feature-auth
feature/new-ui                    /Users/username/projects/repo-feature-new-ui

# 特定 worktree の詳細表示
$ git where feature/auth
Worktree: feature/auth
Path: /Users/username/projects/repo-feature-auth

# JSON 形式で全一覧を取得
$ git where --json
[
  {
    "branch": "main",
    "path": "/Users/username/projects/repo"
  },
  {
    "branch": "feature/auth",
    "path": "/Users/username/projects/repo-feature-auth"
  }
]

# fzf でインタラクティブに選択
$ git where --fzf
# fzf UI が開き、選択すると：
main  /Users/username/projects/repo

# 部分一致での検索
$ git where auth
Worktree: feature/auth
Path: /Users/username/projects/repo-feature-auth

# JSON 出力と組み合わせ
$ git where --json feature/auth
{
  "branch": "feature/auth",
  "path": "/Users/username/projects/repo-feature-auth"
}
```

## Tab 補完

bash completion により、以下の補完が利用可能です：

- オプション補完: `git where --<TAB>` で利用可能なオプションを表示
- worktree 名補完: `git where <TAB>` で既存の worktree 名を補完
- 部分一致補完: `git where feat<TAB>` で feature/* にマッチする worktree を補完

## エラーハンドリング

- 存在しない worktree を指定した場合、エラーメッセージが表示されます
- 部分一致で複数の候補がある場合、候補一覧を表示してエラーとなります
- fzf で何も選択せずに終了した場合、エラーメッセージが表示されます

## 関連ツール

- `git-wcd` - worktree のパスを取得（cd コマンドと組み合わせて使用）
- `git-wcp` - 現在の worktree から別の worktree へファイルをコピー
- `gwq` - Git worktree を管理する基盤ツール