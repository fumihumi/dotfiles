echo "$(brew --prefix bash)" | sudo tee -a /etc/shells
chsh -s $(brew --prefix bash)
