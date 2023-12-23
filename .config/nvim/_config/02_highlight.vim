
" ref: http://potappo.hatenablog.jp/entry/2013/03/31/010746
" MEMO: スペースハイライトするとチカチカするためこれはoff
"" "行頭のスペースの連続をハイライトさせる
"" "Tab文字も区別されずにハイライトされるので、区別したいときはTab文字の表示を別に
"" "設定する必要がある。
"" function! SOLSpaceHilight()
""     syntax match SOLSpace "^\s\+" display containedin=ALL
""     highlight SOLSpace term=underline ctermbg=LightGray
"" endf

"全角スペースをハイライトさせる。
function! JISX0208SpaceHilight()
    syntax match JISX0208Space "　" display containedin=ALL
    highlight JISX0208Space term=underline ctermbg=LightCyan
endf

"syntaxの有無をチェックし、新規バッファと新規読み込み時にハイライトさせる
if has("syntax")
    syntax on
        augroup invisible
        autocmd! invisible
        "autocmd BufNew,BufRead * call SOLSpaceHilight()
        autocmd BufNew,BufRead * call JISX0208SpaceHilight()
    augroup END
endif

" ref: http://baqamore.hatenablog.com/entry/2016/11/15/220358
" コメント中の特定の単語を強調表示する
augroup HilightsForce
  autocmd!
  autocmd WinEnter,BufRead,BufNew,Syntax * :silent! call matchadd('Todo', '\(TODO\):')
  autocmd WinEnter,BufRead,BufNew,Syntax * :silent! call matchadd('Note', '\(NOTE\|MEMO\|INFO\):')
  autocmd WinEnter,BufRead,BufNew,Syntax * :silent! call matchadd('Debug', '\(DEBUG\):')
  "TODO:
  "NOTE:
  "DEBUG:
  "MEMO: colros: https://jonasjacek.github.io/colors/
  autocmd WinEnter,BufRead,BufNew,Syntax * highlight Todo ctermbg=Red guibg=Red guifg=White
  "autocmd WinEnter,BufRead,BufNew,Syntax * highlight Note ctermbg=14 guibg=Aqua guifg=White
  " FIXME: term256 color
  " autocod WinEnter,BufRead,BufNew,Syntax * highlight Note ctermbg=Aqua guibg=Aqua guifg=White
  autocmd WinEnter,BufRead,BufNew,Syntax * highlight Debug ctermbg=Yellow guibg=Yellow guifg=White
augroup END

