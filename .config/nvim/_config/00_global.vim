set number
set cursorline
" set cursorcolumn "コレ入れるとだいぶ重くなる
set laststatus=2   " ステータス行を常に表示
set cmdheight=2    " メッセージ表示欄を2行確保
set showmatch      " 対応する括弧を強調表示
set helpheight=999 " ヘルプを画面いっぱいに開く
set breakindent    " インデントされた行の折り返しを綺麗に

"" 不可視文字の表示記号指定
set list           " 不可視文字を表示
set listchars=tab:..,trail:-,nbsp:%,eol:↲,extends:❯,precedes:❮ ",space:_

""移動設定
set backspace=indent,eol,start "Backspaceキーの影響範囲に制限を設けない
set whichwrap=b,s,h,l,<,>,[,]  "行頭行末の左右移動で行をまたぐ
set scrolloff=8                "上下8行の視界を確保
set sidescrolloff=16           " 左右スクロール時の視界を確保
set sidescroll=1               " 左右スクロールは一文字づつ行う

""ファイル処理関連
set confirm "未保存時に終了前に保存確認
set hidden "未保存時に別のファイルを開くことが出来る
set autoread "外部でファイルに変更がされた場合は読みなおす

""検索設定
set hlsearch "検索文字列をハイライト
set incsearch "インクリメンタルサーチ
set ignorecase "大文字と小文字を区別しない
set smartcase "大文字と小文字が混在した言葉で検索を行った場合に限り、大文字と小文字を区別する
set wrapscan "最後尾まで検索を終えたら次の検索で先頭に移る

""タブ、インデント設定
set expandtab "タブ入力を複数の空白入力に置き換える
set tabstop=2 "画面上でタブ文字が占める幅
set shiftwidth=2 "自動インデントでずれる幅
set softtabstop=2 "連続した空白に対してタブキーやバックスペースキーでカーソルが動く幅
set autoindent "改行時に前の行のインデントを継続する
set smartindent "改行時に入力された行の末尾に合わせて次の行のインデントを増減する

set encoding=utf-8
set fileencodings=utf-8
set fileformats=unix,dos,mac

set nobackup
set noswapfile

set showcmd
set virtualedit=block

""OSのクリップボードをレジスタ指定無しで Yank, Put 出来るようにする
set clipboard=unnamed,unnamedplus

set wildmenu wildmode=list:longest,full "" コマンドラインモードでTABキーによるファイル名補完を有効にする
set history=1000 "" コマンドラインの履歴を1000件保存する
""ビープ音すべてを無効にする
set visualbell t_vb=
set fileencoding=utf-8

