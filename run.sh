#!/bin/sh
# run.sh — télécharge la dernière release de `configs` et l'exécute.
#
# Usage :
#   curl -fsSL https://raw.githubusercontent.com/luannbertaud/configs/main/run.sh | sh
#
# Les arguments sont transmis au binaire ; avec `curl | sh`, utilisez `-s --` :
#   curl -fsSL .../run.sh | sh -s -- --list
#   curl -fsSL .../run.sh | sh -s -- --install nvim
#
# Sans argument, lance le TUI interactif.
set -eu

REPO="luannbertaud/configs"
BIN="configs"

# --- détection de la plateforme -------------------------------------------
os="$(uname -s | tr '[:upper:]' '[:lower:]')"
arch="$(uname -m)"
case "$arch" in
  x86_64 | amd64) arch="amd64" ;;
  aarch64 | arm64) arch="arm64" ;;
  *) echo "run.sh: architecture non supportée: $arch" >&2; exit 1 ;;
esac
case "$os" in
  linux | darwin) ;;
  *) echo "run.sh: OS non supporté: $os (utilise le binaire Windows manuellement)" >&2; exit 1 ;;
esac

asset="${BIN}_${os}_${arch}"
url="https://github.com/${REPO}/releases/latest/download/${asset}"

# --- téléchargement dans un dossier temporaire ----------------------------
tmp="$(mktemp -d "${TMPDIR:-/tmp}/configs.XXXXXX")"
trap 'rm -rf "$tmp"' EXIT
out="${tmp}/${BIN}"

echo "run.sh: téléchargement de ${asset} (dernière release)..." >&2
if ! curl -fsSL --proto '=https' "$url" -o "$out"; then
  echo "run.sh: échec du téléchargement depuis $url" >&2
  echo "run.sh: une release publiée avec cet asset existe-t-elle ?" >&2
  exit 1
fi
chmod +x "$out"

# --- exécution ------------------------------------------------------------
"$out" "$@"
