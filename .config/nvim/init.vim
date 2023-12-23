if exists('g:vscode')
  source ~/.config/nvim/_config/00_global.vim
  finish
endif

let $VIM_CONFIG_HOME = $HOME."/.config/nvim"
let s:minpac_dir = expand($VIM_CONFIG_HOME.'/pack/minpac')

if isdirectory(s:minpac_dir) == 0
    call mkdir(s:minpac_dir .'/opt', "p")
endif
" minpac 本体
let s:minpac_repo_dir = s:minpac_dir .'/opt/minpac'
if isdirectory(s:minpac_dir .'/opt') == 0
    call mkdir(s:minpac_dir .'/opt', "p")
endif

let s:minpac_download = 0

if has('vim_starting')
  if !isdirectory(s:minpac_repo_dir)
    echo "Install minpac ..."
    execute '!git clone --depth 1 https://github.com/k-takata/minpac ' . s:minpac_repo_dir
    let s:minpac_download = 1
  endif
endif

set packpath^=~/.config/nvim
packadd minpac

call map(sort(split(globpath(&runtimepath, '_config/*.vim'))), { ->[execute('exec "so" v:val')] })

syntax enable

"" conflict Errors Highlight
match ErrorMsg '^\(<\|=\|>\)\{7\}\([^=].\+\)\?$'

"ファイルをひらいたとき最後にカーソルがあった場所に移動
augroup vimrcEx
  au BufRead * if line("'\"") > 0 && line("'\"") <= line("$") |
  \ exe "normal g`\"" | endif
augroup END

autocmd BufWritePre * :%s/\s\+$//ge

cabbr w!! w !sudo tee > /dev/null %

" git commit 時にはプラグインは読み込まない
if $HOME != $USERPROFILE && $GIT_EXEC_PATH != ''
  finish
end

set t_Co=256

autocmd ColorScheme * highlight Constant ctermfg=207
silent! colorscheme jellybeans


command! Soba source ~/.config/nvim/init.vim
