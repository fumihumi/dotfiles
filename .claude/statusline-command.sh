#!/usr/bin/env bash
# Claude Code statusLine command
# Converted from PS1 in ~/.bash_profile

input=$(cat)
cwd=$(echo "$input" | jq -r '.cwd')

# ---------------------------------------------------------------------------
# PR number cache (TTL=30s, keyed by repo-root + branch)
# Avoids calling `gh` on every render. Cache files live under /tmp.
# ---------------------------------------------------------------------------
pr_number=""
if command -v gh >/dev/null 2>&1 && git -C "$cwd" rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  _repo_root=$(git -C "$cwd" rev-parse --show-toplevel 2>/dev/null)
  _branch_now=$(git -C "$cwd" symbolic-ref --short HEAD 2>/dev/null)
  if [ -n "$_repo_root" ] && [ -n "$_branch_now" ]; then
    # Build a safe cache filename from repo-root + branch
    _cache_key=$(printf '%s|%s' "$_repo_root" "$_branch_now" | tr -cs 'A-Za-z0-9_-' '_')
    _cache_file="${TMPDIR:-/tmp}/claude_statusline_pr_${_cache_key}"
    _now=$(date +%s)
    _cache_valid=0
    if [ -f "$_cache_file" ]; then
      _cache_mtime=$(stat -f '%m' "$_cache_file" 2>/dev/null || stat -c '%Y' "$_cache_file" 2>/dev/null || echo 0)
      _age=$(( _now - _cache_mtime ))
      [ "$_age" -lt 30 ] && _cache_valid=1
    fi
    if [ "$_cache_valid" -eq 1 ]; then
      pr_number=$(cat "$_cache_file" 2>/dev/null)
    else
      # Fetch from gh with a 5-second timeout.
      # Run gh from the repo root so it picks up the correct remote automatically.
      # -C to cd into the repo root avoids issues when cwd is a subdirectory.
      pr_number=$(cd "$_repo_root" 2>/dev/null && timeout 5 gh pr view --json number -q '.number' 2>/dev/null || true)
      printf '%s' "$pr_number" > "$_cache_file" 2>/dev/null || true
    fi
  fi
fi

# ---------------------------------------------------------------------------
# Git prompt string (replicates git_prompt_string + is_git_worktree from .bash_profile)
# ---------------------------------------------------------------------------
git_info=""
if git -C "$cwd" rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  branch=$(git -C "$cwd" symbolic-ref --short HEAD 2>/dev/null || git -C "$cwd" describe --tags --exact-match 2>/dev/null || git -C "$cwd" rev-parse --short HEAD 2>/dev/null)
  if [ -n "$branch" ]; then
    short_hash=$(git -C "$cwd" rev-parse --short HEAD 2>/dev/null)
    # dirty state indicator
    dirty=""
    if ! git -C "$cwd" diff --quiet 2>/dev/null || ! git -C "$cwd" diff --cached --quiet 2>/dev/null; then
      dirty="*"
    fi
    # worktree indicator
    worktree=""
    toplevel=$(git -C "$cwd" rev-parse --show-toplevel 2>/dev/null)
    if [ -n "$toplevel" ] && [ -f "$toplevel/.git" ]; then
      worktree=" (worktree)"
    fi
    # PR badge (appended after hash, e.g. "main*:fbd15d0 #71")
    pr_badge=""
    if [ -n "$pr_number" ]; then
      pr_badge=" #${pr_number}"
    fi
    git_info=" (${branch}${dirty}:${short_hash}${pr_badge}${worktree})"
  fi
fi

user=$(whoami)
dir=$(basename "$cwd")
time_str=$(date +%H:%M:%S)

# Model name (display name) from JSON input
model_name=$(echo "$input" | jq -r '.model.display_name // empty')

# Context usage percentage (pre-calculated field; null before first API call)
used_pct=$(echo "$input" | jq -r '.context_window.used_percentage // empty')

# ---------------------------------------------------------------------------
# Session metrics from the JSONL transcript.
#
# IMPORTANT: /clear starts a NEW transcript file, so every value computed here
# resets on /clear. The .cost.* fields in the JSON input (total_cost_usd,
# total_duration_ms, total_lines_*) do NOT — the Claude Code *process*
# accumulates them across /clear, so reading them shows stale, pre-clear data.
# We therefore derive all of these from transcript_path instead:
#
#   tokens   : input + output + cache_creation + cache_read usage
#   cost     : recomputed from usage × per-model pricing (USD per 1M tokens)
#   lines    : +/- counted from edit/write structuredPatch hunks
#   runtime  : last transcript timestamp − first
#
# Pricing (USD/MTok): input / output / cache-write (5m, 1.25×) / cache-read (0.1×)
#   Opus   5 / 25 / 6.25 / 0.5     Sonnet 3 / 15 / 3.75 / 0.3
#   Haiku  1 /  5 / 1.25 / 0.1     Fable 10 / 50 / 12.5 / 1.0
# ---------------------------------------------------------------------------
cost_part=""
tokens_part=""
runtime_part=""
lines_part=""
transcript_path=$(echo "$input" | jq -r '.transcript_path // empty')
if [ -n "$transcript_path" ] && [ -f "$transcript_path" ]; then
  metrics=$(jq -rs '
    def rate(m):
      if   (m|test("opus"))         then {i:5,  o:25, cw:6.25, cr:0.5}
      elif (m|test("sonnet"))       then {i:3,  o:15, cw:3.75, cr:0.3}
      elif (m|test("haiku"))        then {i:1,  o:5,  cw:1.25, cr:0.1}
      elif (m|test("fable|mythos")) then {i:10, o:50, cw:12.5, cr:1.0}
      else {i:5, o:25, cw:6.25, cr:0.5} end;
    def ts: .timestamp // empty | sub("\\.[0-9]+Z$";"Z") | fromdateiso8601;
    ( [ .[]
        | select(.type=="assistant" and .message.usage != null)
        | .message as $m | $m.usage as $u | rate($m.model // "") as $r
        | { tok:  (($u.input_tokens//0) + ($u.output_tokens//0)
                   + ($u.cache_creation_input_tokens//0) + ($u.cache_read_input_tokens//0)),
            cost: ((($u.input_tokens//0)*$r.i + ($u.output_tokens//0)*$r.o
                    + ($u.cache_creation_input_tokens//0)*$r.cw
                    + ($u.cache_read_input_tokens//0)*$r.cr) / 1000000) }
      ] ) as $usage
    | ( [ .[]
          | select((.toolUseResult|type)=="object")
          | .toolUseResult.structuredPatch // []
          | .[].lines[] ] ) as $lines
    | ( [ .[] | ts ] ) as $times
    | [ ($usage | map(.tok)  | add // 0),
        ($usage | map(.cost) | add // 0),
        ([ $lines[] | select(startswith("+")) ] | length),
        ([ $lines[] | select(startswith("-")) ] | length),
        (if ($times|length) > 0 then (($times|max) - ($times|min)) else 0 end)
      ] | @tsv
  ' "$transcript_path" 2>/dev/null)

  IFS=$'\t' read -r total_tokens total_cost lines_added lines_removed duration_s <<< "$metrics"

  # Cumulative tokens (dim gray)
  tok_int=$(printf '%.0f' "${total_tokens:-0}" 2>/dev/null || echo 0)
  if [ "$tok_int" -gt 0 ] 2>/dev/null; then
    if   [ "$tok_int" -ge 1000000 ]; then tok_fmt=$(awk "BEGIN{printf \"%.1fM\", $tok_int/1000000}")
    elif [ "$tok_int" -ge 1000 ];    then tok_fmt=$(awk "BEGIN{printf \"%.1fk\", $tok_int/1000}")
    else tok_fmt="$tok_int"; fi
    tokens_part="\033[90m${tok_fmt} tok\033[00m"
  fi

  # Session cost (cyan)
  if [ -n "$total_cost" ]; then
    cost_fmt=$(printf '%.2f' "$total_cost" 2>/dev/null)
    [ -n "$cost_fmt" ] && cost_part="\033[36m\$${cost_fmt}\033[00m"
  fi

  # Runtime (dim gray) — wall-clock span of this (post-clear) transcript
  dur_int=$(printf '%.0f' "${duration_s:-0}" 2>/dev/null || echo 0)
  if [ "$dur_int" -gt 0 ] 2>/dev/null; then
    if [ "$dur_int" -ge 60 ]; then
      runtime_fmt="$(( dur_int / 60 ))m $(( dur_int % 60 ))s"
    else
      runtime_fmt="${dur_int}s"
    fi
    runtime_part="\033[90m${runtime_fmt}\033[00m"
  fi

  # Changed lines (green+/red-) — only shown once edits have happened this session
  if [ "${lines_added:-0}" -gt 0 ] 2>/dev/null || [ "${lines_removed:-0}" -gt 0 ] 2>/dev/null; then
    lines_part="\033[32m+${lines_added}\033[00m \033[31m-${lines_removed}\033[00m"
  fi
fi

# ---------------------------------------------------------------------------
# Build model/ctx/cost/tokens/runtime segments (joined onto Line 2)
# Color for context usage: green (<50%), yellow (50-80%), red (>80%)
# ---------------------------------------------------------------------------

# Collect each segment into an array-like list using a separator
# We'll build the line piece by piece, inserting " | " between non-empty parts.
sep="\033[90m|\033[00m"

line3_parts=""
_append_part() {
  local part="$1"
  if [ -n "$part" ]; then
    if [ -n "$line3_parts" ]; then
      line3_parts="${line3_parts} ${sep} ${part}"
    else
      line3_parts="${part}"
    fi
  fi
}

# Model name in magenta
if [ -n "$model_name" ]; then
  _append_part "\033[35m${model_name}\033[00m"
fi

# Context % with traffic-light coloring
if [ -n "$used_pct" ]; then
  used_int=$(printf '%.0f' "$used_pct")
  if [ "$used_int" -ge 80 ]; then
    ctx_color="\033[31m"   # red
  elif [ "$used_int" -ge 50 ]; then
    ctx_color="\033[33m"   # yellow
  else
    ctx_color="\033[32m"   # green
  fi
  _append_part "${ctx_color}ctx:${used_int}%\033[00m"
fi

# Session cost (cyan)
_append_part "$cost_part"

# Cumulative tokens (dim gray)
_append_part "$tokens_part"

# Runtime (dim gray)
_append_part "$runtime_part"

# Changed lines (green+/red-)
_append_part "$lines_part"

# ---------------------------------------------------------------------------
# Build Line 2: directory  (+ optional model/ctx/cost/… segments on same line)
# ---------------------------------------------------------------------------
line2="\033[32m(${dir})\033[00m"
if [ -n "$line3_parts" ]; then
  line2="${line2} ${sep} ${line3_parts}"
fi

printf "\033[34m{ %s }\033[31m%s\033[33m [%s]\033[34m\n%b\n" \
  "$user" "$git_info" "$time_str" "$line2"
