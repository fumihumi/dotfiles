# git-diff-summary

Git のブランチ間の差分を、ファイルごとの最終更新者と日付を含めて表示するツールです。

## 機能

- ブランチ間の差分ファイル一覧を最終更新者情報と共に表示
- 削除されたファイルの明示的な表示
- ファイル変更統計（追加/削除行数）の表示
- JSON形式での出力サポート
- 並行処理による高速な情報取得
- Bash補完サポート

## インストール

```bash
go build -o git-diff-summary .
```

## 使用方法

### Git サブコマンドとして使用

```bash
# 基本的な使用方法
git diff-summary origin/main origin/develop

# 特定のパスでフィルタリング
git diff-summary origin/main origin/develop db/migrate/

# 整形された出力
git diff-summary --format origin/main origin/develop

# 統計情報付きで表示
git diff-summary --stat origin/main origin/develop

# JSON形式で出力
git diff-summary --json origin/main origin/develop
```

### 直接実行

```bash
git-diff-summary origin/main origin/develop
```

## オプション

```text
USAGE:
    git diff-summary [OPTIONS] <SRC_BRANCH> <DST_BRANCH> [file_path]

ARGUMENTS:
    <SRC_BRANCH>    比較元ブランチ (例: origin/main)
    <DST_BRANCH>    比較先ブランチ (例: origin/develop)
    [file_path]     フィルタリング用のパス (例: db/migrate, app/models)

OPTIONS:
    -f, --format        整列された列で出力
    --stat              ファイル変更統計を表示 (A/M/D/R/C)
    --json              JSON形式で出力
    -w, --workers int   並行処理のワーカー数 (デフォルト: 10)
    -h, --help          ヘルプを表示
```

## 出力形式

### デフォルト形式

```text
src/main.go: John Doe (2023/12/01)
src/utils.go: Jane Smith (2023/11/30)
README.md: [DELETED] Alice Johnson (2023/11/29)
```

### 整形出力（--format）

```text
src/main.go                          John Doe      (2023/12/01)
src/utils.go                         Jane Smith    (2023/11/30)
README.md                            [DELETED]     Alice Johnson (2023/11/29)
```

### 統計情報付き（--stat）

```text
M  src/main.go     |  25 +++++++++++++++++++++++++
A  src/utils.go    |  15 +++++++++++++++
D  README.md       |  10 ----------
 3 files changed, 40 insertions(+), 10 deletions(-)
```

### JSON形式（--json）

```json
[
  {
    "file": "src/main.go",
    "author": "John Doe",
    "date": "2023/12/01",
    "status": "M",
    "insertions": 25,
    "deletions": 0
  },
  {
    "file": "README.md",
    "author": "Alice Johnson",
    "date": "2023/11/29",
    "deleted": true,
    "status": "D",
    "deletions": 10
  }
]
```

## Bash補完

補完スクリプトの生成：

```bash
# 補完スクリプトを生成
git-diff-summary completion > ~/.bash_completion.d/git-diff-summary

# または ~/.bashrc に追加
eval "$(git-diff-summary completion)"
```

補完機能：

- ブランチ名の自動補完
- オプションの補完
- ファイルパスの補完

## 使用例

### プルリクエストレビュー時

```bash
# マージ対象の変更ファイルと最終更新者を確認
git diff-summary origin/main feature/new-feature

# 特定のディレクトリのみチェック
git diff-summary origin/main feature/new-feature src/api/
```

### リリース準備時

```bash
# リリースブランチとmainの差分を確認
git diff-summary origin/main origin/release/v1.2.0 --stat

# JSON形式で記録を残す
git diff-summary origin/main origin/release/v1.2.0 --json > release-changes.json
```

### 並行処理の調整

```bash
# 大規模なリポジトリで並行処理数を増やす
git diff-summary -w 20 origin/main origin/develop

# 低スペック環境で並行処理数を減らす
git diff-summary -w 5 origin/main origin/develop
```

## テスト

```bash
# 全テストを実行
go test -v

# 短時間テスト（統合テストをスキップ）
go test -short -v

# カバレッジ付きテスト
go test -cover

# ベンチマークテスト
go test -bench=.

# 特定のテストのみ実行
go test -run TestParseStatOutput -v
```

## 必要な環境

- Go 1.21以上
- Git 2.0以上
