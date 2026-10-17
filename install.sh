#!/bin/sh
# Install or uninstall the aliases from github.com/luannbertaud/configs (branch: aliases).
#
# Standalone, downloads aliases.sh:
#   curl -fsSL <raw>/install.sh | sh
#   curl -fsSL <raw>/install.sh | sh -s -- --uninstall
#
# From a clone, or via the `configs` package manager, uses the local aliases.sh:
#   ./install.sh
#   ./install.sh --uninstall
set -eu

URL="https://raw.githubusercontent.com/luannbertaud/configs/aliases/aliases.sh"
DIR="${CONFIGS_TARGET:-${XDG_CONFIG_HOME:-$HOME/.config}/aliases}"
FILE="$DIR/aliases.sh"
MARKER="# Aliases from github.com/luannbertaud/configs (branch: aliases)"
LINE="[ -f \"$FILE\" ] && . \"$FILE\""

# Local mode when the package manager provides CONFIGS_TARGET, or when
# install.sh sits next to aliases.sh (running from a clone).
SELF_DIR=""
case "$0" in
    *install.sh)
        SELF_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" 2>/dev/null && pwd) || SELF_DIR=""
        ;;
esac
LOCAL=0
if [ -n "${CONFIGS_TARGET:-}" ] || { [ -n "$SELF_DIR" ] && [ -f "$SELF_DIR/aliases.sh" ]; }; then
    LOCAL=1
fi

# Existing rc files, or ~/.bashrc if there is none.
rc_files() {
    found=""
    for rc in "$HOME/.bashrc" "$HOME/.zshrc"; do
        [ -f "$rc" ] && found="$found $rc"
    done
    echo "${found:-$HOME/.bashrc}"
}

install() {
    if [ "$LOCAL" = 1 ]; then
        [ -f "$FILE" ] || { echo "aliases: $FILE not found" >&2; exit 1; }
        echo "Using $FILE"
    else
        mkdir -p "$DIR"
        curl -fsSL "$URL" -o "$FILE.tmp"
        mv "$FILE.tmp" "$FILE"
        echo "Downloaded $FILE"
    fi

    for rc in $(rc_files); do
        if [ -f "$rc" ] && grep -qxF "$MARKER" "$rc"; then
            echo "Already sourced in $rc"
        else
            printf '\n%s\n%s\n' "$MARKER" "$LINE" >> "$rc"
            echo "Added source line to $rc"
        fi
    done
    echo "Open a new shell or run: . \"$FILE\""
}

uninstall() {
    for rc in "$HOME/.bashrc" "$HOME/.zshrc"; do
        [ -f "$rc" ] && grep -qxF "$MARKER" "$rc" || continue
        tmp=$(mktemp)
        # Drop the marker, the line after it, and one blank line before it.
        awk -v m="$MARKER" '
            skip { skip = 0; next }
            $0 == m { skip = 1; if (n) n--; while (n-- > 0) print ""; n = 0; next }
            /^$/ { n++; next }
            { while (n-- > 0) print ""; n = 0; print }
            END { while (n-- > 0) print "" }
        ' "$rc" > "$tmp"
        cat "$tmp" > "$rc"
        rm -f "$tmp"
        echo "Removed source line from $rc"
    done

    # In a clone, the package manager removes the directory itself.
    if [ "$LOCAL" = 0 ]; then
        rm -f "$FILE"
        rmdir "$DIR" 2>/dev/null && echo "Removed $DIR" || true
    fi
    echo "Uninstalled (aliases stay defined in already open shells)"
}

case "${1:-}" in
    "") install ;;
    --uninstall) uninstall ;;
    *) echo "Usage: install.sh [--uninstall]" >&2; exit 1 ;;
esac
