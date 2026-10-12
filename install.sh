#!/bin/sh
# Install or uninstall the aliases from github.com/luannbertaud/configs (branch: aliases).
#   curl -fsSL <raw>/install.sh | sh
#   curl -fsSL <raw>/install.sh | sh -s -- --uninstall
set -eu

URL="https://raw.githubusercontent.com/luannbertaud/configs/aliases/aliases.sh"
DIR="${XDG_CONFIG_HOME:-$HOME/.config}/aliases"
FILE="$DIR/aliases.sh"
MARKER="# Aliases from github.com/luannbertaud/configs (branch: aliases)"
LINE="[ -f \"$FILE\" ] && . \"$FILE\""

# Existing rc files, or ~/.bashrc if there is none.
rc_files() {
    found=""
    for rc in "$HOME/.bashrc" "$HOME/.zshrc"; do
        [ -f "$rc" ] && found="$found $rc"
    done
    echo "${found:-$HOME/.bashrc}"
}

install() {
    mkdir -p "$DIR"
    curl -fsSL "$URL" -o "$FILE.tmp"
    mv "$FILE.tmp" "$FILE"
    echo "Downloaded $FILE"

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

    rm -f "$FILE"
    rmdir "$DIR" 2>/dev/null && echo "Removed $DIR" || true
    echo "Uninstalled (aliases stay defined in already open shells)"
}

case "${1:-}" in
    "") install ;;
    --uninstall) uninstall ;;
    *) echo "Usage: install.sh [--uninstall]" >&2; exit 1 ;;
esac
