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
#
# Le binaire est mis en cache dans ${XDG_CACHE_HOME:-$HOME/.cache}/configs et
# réutilisé tant que la dernière release n'a pas changé. Supprimer ce dossier
# force un nouveau téléchargement.
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

# --- cache ----------------------------------------------------------------
cache_dir="${XDG_CACHE_HOME:-$HOME/.cache}/${BIN}"
out="${cache_dir}/${asset}"
tag_file="${cache_dir}/release-tag"

# Tag de la dernière release, sans l'API GitHub (petite requête HEAD).
latest_tag=""
latest_tag="$(curl -fsSIL -o /dev/null -w '%{url_effective}' \
  "https://github.com/${REPO}/releases/latest" 2>/dev/null |
  sed -n 's#.*/releases/tag/##p' || true)"

cached_tag=""
if [ -f "$tag_file" ]; then
  cached_tag="$(cat "$tag_file" 2>/dev/null || true)"
fi

# Cache valide : même release et binaire exécutable.
if [ -x "$out" ] && [ -n "$latest_tag" ] && [ "$latest_tag" = "$cached_tag" ]; then
  exec "$out" "$@"
fi

# GitHub injoignable : on garde la possibilité d'utiliser le cache.
if [ -z "$latest_tag" ]; then
  if [ -x "$out" ]; then
    echo "run.sh: GitHub injoignable, utilisation du binaire en cache (${cached_tag:-version inconnue})" >&2
    exec "$out" "$@"
  fi
  echo "run.sh: impossible de joindre GitHub et aucun binaire en cache" >&2
  exit 1
fi

# --- téléchargement -------------------------------------------------------
mkdir -p "$cache_dir"
tmp="$(mktemp "${cache_dir}/.${BIN}.XXXXXX")"
trap 'rm -f "$tmp"' EXIT INT TERM

echo "run.sh: téléchargement de ${asset} (release ${latest_tag})..." >&2
if ! curl -fsSL --proto '=https' "$url" -o "$tmp"; then
  if [ -x "$out" ]; then
    echo "run.sh: échec du téléchargement, utilisation du binaire en cache (${cached_tag:-version inconnue})" >&2
    exec "$out" "$@"
  fi
  echo "run.sh: échec du téléchargement depuis $url" >&2
  echo "run.sh: une release publiée avec cet asset existe-t-elle ?" >&2
  exit 1
fi

chmod +x "$tmp"
mv "$tmp" "$out"
printf '%s\n' "$latest_tag" > "$tag_file"
trap - EXIT INT TERM

# --- exécution ------------------------------------------------------------
exec "$out" "$@"
