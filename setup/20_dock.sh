defaults write com.apple.dock autohide -bool true
defaults write com.apple.dock autohide-delay -float 1
defaults write com.apple.dock autohide-time-modifier -int 0
defaults write com.apple.dock persistent-apps -array
defaults write com.apple.dock magnification -bool true
defaults write com.apple.dock tilesize -int 30
defaults write com.apple.dock largesize -int 50
defaults write com.apple.dock magnification -bool true
defaults write com.apple.dock orientation -string "left"
killall Dock
