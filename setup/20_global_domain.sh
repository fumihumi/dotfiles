defaults write NSGlobalDomain AppleShowAllExtensions -bool true    # 全ての拡張子のファイルを表示する
defaults write NSGlobalDomain AppleShowScrollBars -string "Always"    # スクロールバーを常時表示する
defaults write NSGlobalDomain com.apple.springing.enabled -bool false
defaults write NSGlobalDomain KeyRepeat -int 2
defaults write NSGlobalDomain InitialKeyRepeat -int 15   # キーリピート開始までのタイミング
defaults write NSGlobalDomain NSWindowResizeTime -float 0.001    # コンソールアプリケーションの画面サイズ変更を高速にする
defaults write NSGlobalDomain WebKitDeveloperExtras -bool true    # Safari のコンテキストメニューに Web インスペクタを追加する
defaults write NSGlobalDomain NSAutomaticCapitalizationEnabled -bool false # 自動大文字の無効化

defaults write com.apple.iphonesimulator AllowFullscreenMode -bool YES
