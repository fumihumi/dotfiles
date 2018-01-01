" ref: https://gist.github.com/ony/efd1de13181f4bcebaf5795a4ee46905
" Na隱ｰve autoload of packs
" See https://github.com/k-takata/minpac/issues/28
" And https://github.com/junegunn/vim-plug/blob/9813d5e/plug.vim#L271-L276
function! minautopac#add(repo, ...)
  let l:opts = get(a:000, 0, {})
  if has_key(l:opts, 'for')
    let l:name = substitute(a:repo, '^.*/', '', '')
    let l:ft = l:opts.for  " TODO: support array
    execut printf('autocmd FileType %s packadd %s', l:ft, l:name)
  endif
endfunction

if exists('*minpac#init')
  call minpac#init()

  command! -nargs=+ -bar Plugin call minpac#add(<args>) | call minautopac#add(<args>)

  call minpac#add('k-takata/minpac', {'type': 'opt'})
  call minpac#add('vim-jp/syntax-vim-ex')
  "for vim-lsp-settings
  call minpac#add('prabirshrestha/async.vim')
  call minpac#add('prabirshrestha/asyncomplete.vim')
  call minpac#add('prabirshrestha/asyncomplete-lsp.vim')
  call minpac#add('prabirshrestha/vim-lsp')
  call minpac#add('mattn/vim-lsp-settings')
  call minpac#add('mattn/vim-lsp-icons')
  " vim-airline
  call minpac#add('vim-airline/vim-airline')
  call minpac#add('vim-airline/vim-airline-themes')
  "
  call minpac#add('tpope/vim-pathogen')
  call minpac#add('scrooloose/nerdtree')
  call minpac#add('thinca/vim-quickrun')
  call minpac#add('bronson/vim-trailing-whitespace')
  call minpac#add('Yggdroot/indentLine')
  call minpac#add('cohama/lexima.vim')
  call minpac#add('simeji/winresizer')
  call minpac#add('tyru/open-browser.vim')
  call minpac#add('godlygeek/tabular')
  call minpac#add('plasticboy/vim-markdown')
  call minpac#add('previm/previm')

  " ColorScheme
  call minpac#add('nanotech/jellybeans.vim')
  call minpac#add('tomasr/molokai')
  " fzf
  call minpac#add('junegunn/fzf')
  call minpac#add('junegunn/fzf.vim')

  " 一部のfiletypeで利用するplugins
  Plugin 'pangloss/vim-javascript', {'type': 'opt', 'for': 'javascript,javascript.jsx'}
  Plugin 'maxmellon/vim-jsx-pretty', {'type': 'opt', 'for': 'javascript,javascript.jsx'}
  Plugin 'alampros/vim-styled-jsx', {'type': 'opt', 'for': 'javascript,javascript.jsx'}
  Plugin 'moll/vim-node', {'type': 'opt', 'for': 'javascript'}
  Plugin 'leafgarland/typescript-vim', {'type': 'opt', 'for': 'typescript.tsx'}
  Plugin 'ianks/vim-tsx', {'type': 'opt', 'for': 'typescript.tsx'}
  Plugin 'styled-components/vim-styled-components', {'type': 'opt', 'for': 'javascript.jsx,typescript.tsx'}
  Plugin 'posva/vim-vue', {'type': 'opt', 'for': 'vue'}
  Plugin 'mattn/vim-goimports', {'type': 'opt', 'for': 'go'}
endif

command! -bar PackUpdate packadd minpac | runtime minautopac.vim | call minpac#update()
command! -bar PackClean packadd minpac | runtime minautopac.vim | call minpac#clean()
"command! PackUpdate packadd minpac | source $MYVIMRC | call minpac#update('', {'do': 'call minpac#status()'})
"command! PackClean  packadd minpac | source $MYVIMRC | call minpac#clean()
command! PackStatus packadd minpac | source $MYVIMRC | call minpac#status()

runtime! OPT ftdetect/*.vim

augroup MinAutoPac
  autocmd!
  runtime! myplugins.vim
augroup END
