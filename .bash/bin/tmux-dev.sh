#!/bin/bash

# 🤖 Claude Code / Codex grid tmux launcher
#
# tmux セッションを起動し、Claude Code (または Codex) を rows × cols の
# グリッドで立ち上げる。デフォルトは「縦 3 行 × 横 2 列 = 6 ペイン」。
# デフォルトはデタッチ状態でセッションを作成するだけ。-a でそのままアタッチする。
#
#   例) -r 3 -c 2
#   +--------+--------+
#   | pane   | pane   |
#   +--------+--------+
#   | pane   | pane   |
#   +--------+--------+
#   | pane   | pane   |
#   +--------+--------+
#
# Usage:
#   tmux-dev.sh [-r ROWS] [-c COLS] [-s SESSION] [-d WORKDIR] [-x] [-a]
#
#   -r ROWS    : 縦の分割数 (行)            default: 3
#   -c COLS    : 横の分割数 (列)            default: 2
#   -s SESSION : tmux セッション名          default: claude (codex モード時: codex)
#   -d WORKDIR : 各ペインの作業ディレクトリ default: リポジトリルート
#   -x         : Codex を起動する (default: Claude Code)
#   -a         : 起動後そのまま tmux にアタッチする (default: detach)
#
# 環境変数:
#   CLAUDE_CMD : Claude モードで実行するコマンド default: claude
#   CODEX_CMD  : Codex モードで実行するコマンド  default: codex

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"

# --- デフォルト値 ---
ROWS=3
COLS=2
SESSION=""
WORK_DIR="$ROOT_DIR"
CLAUDE_CMD="${CLAUDE_CMD:-claude}"
CODEX_CMD="${CODEX_CMD:-codex}"
MODE="claude"
ATTACH=0

usage() {
  sed -n '3,30p' "$0" | sed 's/^# \{0,1\}//'
  exit "${1:-0}"
}

# --- 引数パース ---
while getopts ":r:c:s:d:xah" opt; do
  case "$opt" in
    r) ROWS="$OPTARG" ;;
    c) COLS="$OPTARG" ;;
    s) SESSION="$OPTARG" ;;
    d) WORK_DIR="$OPTARG" ;;
    x) MODE="codex" ;;
    a) ATTACH=1 ;;
    h) usage 0 ;;
    :) echo "❌ オプション -$OPTARG には引数が必要です" >&2; usage 1 ;;
    \?) echo "❌ 不明なオプション: -$OPTARG" >&2; usage 1 ;;
  esac
done

# --- モードに応じたコマンド / セッション名の決定 ---
if [ "$MODE" = "codex" ]; then
  RUN_CMD="$CODEX_CMD"
  SESSION="${SESSION:-codex}"
else
  RUN_CMD="$CLAUDE_CMD"
  SESSION="${SESSION:-claude}"
fi

# --- バリデーション ---
if ! [[ "$ROWS" =~ ^[1-9][0-9]*$ ]] || ! [[ "$COLS" =~ ^[1-9][0-9]*$ ]]; then
  echo "❌ ROWS / COLS は 1 以上の整数で指定してください (ROWS=$ROWS, COLS=$COLS)" >&2
  exit 1
fi

TOTAL=$((ROWS * COLS))
if [ "$TOTAL" -gt 64 ]; then
  echo "❌ ペイン数が多すぎます (${ROWS}×${COLS}=${TOTAL})。64 以下にしてください。" >&2
  exit 1
fi

# --- 依存チェック ---
if ! command -v tmux >/dev/null 2>&1; then
  echo "❌ tmux が見つかりません。先に tmux をインストールしてください (例: brew install tmux)" >&2
  exit 1
fi
if ! command -v "${RUN_CMD%% *}" >/dev/null 2>&1; then
  echo "❌ '${RUN_CMD%% *}' が見つかりません。CLI をインストールしてください。" >&2
  exit 1
fi

# 既存セッションがあれば二重起動を防止。-a 指定時はアタッチする。
if tmux has-session -t "$SESSION" 2>/dev/null; then
  if [ "$ATTACH" -eq 1 ]; then
    echo "ℹ️  既存セッション '$SESSION' が見つかりました。アタッチします。"
    exec tmux attach-session -t "$SESSION"
  fi
  echo "ℹ️  既存セッション '$SESSION' が見つかりました。アタッチするには: tmux attach -t $SESSION"
  exit 0
fi

# pane を direction 方向に n 等分し、生成された pane id を空白区切りで出力する。
#   split_even <-h|-v> <n> <start_pane_id> <workdir>
# -p は「新しく作られる側 (末尾) のペイン」に割り当てる割合。
# 末尾ペインを連鎖的に分割することで、ほぼ均等な n 分割を得る。
split_even() {
  local dir="$1" n="$2" start="$3" wd="$4"
  local ids=("$start")
  local cur="$start" k pct new
  for ((k = 1; k < n; k++)); do
    pct=$((100 * (n - k) / (n - k + 1)))
    new="$(tmux split-window "$dir" -p "$pct" -P -F '#{pane_id}' -t "$cur" -c "$wd")"
    ids+=("$new")
    cur="$new"
  done
  echo "${ids[@]}"
}

echo "🚀 tmux セッション '$SESSION' を作成中 (${ROWS} 行 × ${COLS} 列 = ${TOTAL} ペイン, mode=${MODE})"
echo "📁 作業ディレクトリ: $WORK_DIR"

# セッション作成 (起点ペイン)
tmux new-session -d -s "$SESSION" -n "$MODE" -c "$WORK_DIR"
FIRST="$(tmux display-message -p -t "$SESSION":0 '#{pane_id}')"

# まず横に COLS 分割して各列の起点ペインを得る
read -r -a COL_PANES <<<"$(split_even -h "$COLS" "$FIRST" "$WORK_DIR")"

# 各列を縦に ROWS 分割し、全ペイン id を集約
ALL_PANES=()
for col in "${COL_PANES[@]}"; do
  read -r -a row_panes <<<"$(split_even -v "$ROWS" "$col" "$WORK_DIR")"
  ALL_PANES+=("${row_panes[@]}")
done

# 各ペインでコマンドを起動
for pane in "${ALL_PANES[@]}"; do
  tmux send-keys -t "$pane" "$RUN_CMD" C-m
done

# 起点ペインを選択
tmux select-pane -t "$FIRST"

if [ "$ATTACH" -eq 1 ]; then
  echo "✅ 起動完了。アタッチします (デタッチ: Ctrl-b d)"
  exec tmux attach-session -t "$SESSION"
fi

echo "✅ 起動完了 (detach 状態)。アタッチするには: tmux attach -t $SESSION"
