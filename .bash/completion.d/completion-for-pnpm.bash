_pnpm_fzf_run_only_or_default() {
  local cur prev
  COMPREPLY=()
  cur="${COMP_WORDS[COMP_CWORD]}"
  prev="${COMP_WORDS[COMP_CWORD-1]}"

  # jq が無ければそのまま既存 completion に委譲
  if ! command -v jq >/dev/null 2>&1; then
    type _pnpm_completion >/dev/null 2>&1 && _pnpm_completion
    return 0
  fi

  # 対象:
  #   pnpm run <TAB>        → fzf
  #   pnpm run bu<TAB>      → fzf + query=bu
  if [[ "${COMP_WORDS[1]}" == "run" && $COMP_CWORD -eq 2 ]]; then
    [[ -f package.json ]] || return 0

    local scripts
    scripts="$(jq -r '.scripts | keys[]' package.json 2>/dev/null)" || return 0
    [[ -z "$scripts" ]] && return 0

    local choice
    choice="$(printf '%s\n' "$scripts" \
      | fzf --height=40% --reverse --query "$cur")"

    [[ -n "$choice" ]] && COMPREPLY=("$choice")
    return 0
  fi

  if type _pnpm_completion >/dev/null 2>&1; then
    _pnpm_completion
  fi
}

# 既存 completion は解除せず、上書き登録してフォールバックする
complete -o default -F _pnpm_fzf_run_only_or_default pnpm
