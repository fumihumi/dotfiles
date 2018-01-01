"vim-markdown
let g:vim_markdown_override_foldtext = 0
let g:vim_markdown_folding_level = 6
let g:vim_markdown_conceal = 0
let g:vim_markdown_conceal_code_blocks = 0

set nofoldenable

"previm
" let g:previm_enable_realtime = 1
let g:previm_open_cmd = ''
nnoremap [previm] <Nop>
nmap <Space>p [previm]
nnoremap <silent> [previm]o :<C-u>PrevimOpen<CR>
nnoremap <silent> [previm]r :call previm#refresh()<CR>
