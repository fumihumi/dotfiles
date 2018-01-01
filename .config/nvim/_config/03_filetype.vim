""オートコメントアウト無効
au FileType * setlocal formatoptions-=ro

autocmd BufNewFile,BufRead *.rc setfiletype sh
autocmd BufRead,BufNewFile *.jsx set filetype=javascript.jsx
autocmd BufNewFile,BufRead *.tsx set filetype=typescript.tsx
autocmd BufNewFile,BufRead *.rb setfiletype ruby
autocmd BufNewFile,BufRead *.jbuilder setfiletype ruby

au FileType go setlocal sw=4 ts=4 sts=4 noet
au FileType go setlocal makeprg=go\ build\ ./... errorformat=%f:%l:\ %m
