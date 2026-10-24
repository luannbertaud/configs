# Shell aliases (bash/zsh).
# https://github.com/luannbertaud/configs/tree/aliases

alias ll='ls -lah'
alias n='nvim'

# Gestionnaire de paquets, lancé via curl. Les arguments sont transmis :
#   configs --list · configs --install nvim · configs --update
configs() {
    curl -fsSL https://raw.githubusercontent.com/luannbertaud/configs/main/run.sh | sh -s -- "$@"
}

alias gst='git status'
alias gadog='git log --all --decorate --oneline --graph'

# Meta: manage this config. Git clone (package manager): git pull + local
# install.sh. Standalone (curl): install.sh via curl. Session: reload only.
alias aliases-update='DIR="${XDG_CONFIG_HOME:-$HOME/.config}/aliases"; if [ -d "$DIR/.git" ]; then ( cd "$DIR" && git pull --ff-only https://github.com/luannbertaud/configs.git aliases ) && . "$DIR/aliases.sh"; elif [ -f "$DIR/aliases.sh" ]; then curl -fsSL https://raw.githubusercontent.com/luannbertaud/configs/aliases/install.sh | sh && . "$DIR/aliases.sh"; else eval "$(curl -fsSL https://raw.githubusercontent.com/luannbertaud/configs/aliases/aliases.sh)"; fi'
alias aliases-uninstall='DIR="${XDG_CONFIG_HOME:-$HOME/.config}/aliases"; if [ -d "$DIR/.git" ]; then sh "$DIR/install.sh" --uninstall && rm -rf -- "$DIR" && echo "aliases: désinstallé"; elif [ -f "$DIR/aliases.sh" ]; then curl -fsSL https://raw.githubusercontent.com/luannbertaud/configs/aliases/install.sh | sh -s -- --uninstall; else echo "aliases: non installé (mode session), rien à désinstaller" >&2; false; fi'
