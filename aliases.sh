# Shell aliases (bash/zsh).
# https://github.com/luannbertaud/configs/tree/aliases

alias ll='ls -lah'

alias gst='git status'
alias gadog='git log --all --decorate --oneline --graph'

# Meta: manage this config. Installed: update + reload. Session mode (not installed): reload only.
alias aliases-update='if [ -f "${XDG_CONFIG_HOME:-$HOME/.config}/aliases/aliases.sh" ]; then curl -fsSL https://raw.githubusercontent.com/luannbertaud/configs/aliases/install.sh | sh && . "${XDG_CONFIG_HOME:-$HOME/.config}/aliases/aliases.sh"; else eval "$(curl -fsSL https://raw.githubusercontent.com/luannbertaud/configs/aliases/aliases.sh)"; fi'
alias aliases-uninstall='if [ -f "${XDG_CONFIG_HOME:-$HOME/.config}/aliases/aliases.sh" ]; then curl -fsSL https://raw.githubusercontent.com/luannbertaud/configs/aliases/install.sh | sh -s -- --uninstall; else echo "aliases: not installed (session mode), nothing to uninstall" >&2; false; fi'
