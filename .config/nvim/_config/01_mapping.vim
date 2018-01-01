let mapleader = "\<Space>"
nmap <Esc><Esc> :nohlsearch<CR><Esc>

nmap ss :split<Return><C-w>w
nmap sv :vsplit<Return><C-w>w

nmap <Leader>w [window]
nnoremap [window]h <C-w>h
nnoremap [window]j <C-w>j
nnoremap [window]k <C-w>k
nnoremap [window]l <C-w>l

" 逕ｻ髱｢繝ｪ繧ｵ繧､繧ｺ
nnoremap [window]L <C-w>>
nnoremap [window]H <C-w><
nnoremap [window]J <C-w>+
nnoremap [window]K <C-w>-

nmap <Leader>b [buf]
nnoremap <silent> [buf]p :bprev<CR>
nnoremap <silent> [buf]n :bnext<CR>
nnoremap <silent> [buf]f :bfirst<CR>
nnoremap <silent> [buf]l :blast<CR>
nnoremap <silent> [buf]d :bdelete<CR>

