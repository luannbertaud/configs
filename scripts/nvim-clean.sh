#!/usr/bin/env bash
# ============================================================================
#  nvim-clean — efface les DONNÉES d'une config Neovim isolée (NVIM_APPNAME).
#
#  Supprime  : data/ , state/ , cache/ propres à la config.
#  Conserve  : ~/.config/<nom>  (ta configuration, versionnée).
#  Ne touche PAS : $XDG_RUNTIME_DIR (partagé entre configs, éphémère).
#
#  Usage :
#     nvim-clean.sh <nom> [--dry-run] [-y|--yes]
#
#  Exemples :
#     nvim-clean.sh nvim-rc
#     nvim-clean.sh nvim-rc --dry-run
# ============================================================================
set -euo pipefail

die() { printf 'erreur: %s\n' "$*" >&2; exit 1; }

usage() {
  cat <<'EOF'
Usage: nvim-clean.sh <nom-config> [--dry-run] [-y|--yes]

  <nom-config>   valeur de NVIM_APPNAME (ex: nvim-rc)
  --dry-run      affiche ce qui serait supprimé, sans rien supprimer
  -y, --yes      ne pas demander de confirmation

Supprime les données de la config (data, state, cache) et conserve
~/.config/<nom-config>.
EOF
}

[ $# -ge 1 ] || { usage; exit 1; }
name="$1"; shift

dry_run=0
assume_yes=0
while [ $# -gt 0 ]; do
  case "$1" in
    --dry-run) dry_run=1 ;;
    -y|--yes)  assume_yes=1 ;;
    -h|--help) usage; exit 0 ;;
    *) die "option inconnue: $1" ;;
  esac
  shift
done

# --- Garde-fou : valider le nom pour ne jamais viser un chemin dangereux -----
case "$name" in
  "" | "." | ".." | */*) die "nom de config invalide: '$name'" ;;
esac

config_home="${XDG_CONFIG_HOME:-$HOME/.config}"
data_home="${XDG_DATA_HOME:-$HOME/.local/share}"
state_home="${XDG_STATE_HOME:-$HOME/.local/state}"
cache_home="${XDG_CACHE_HOME:-$HOME/.cache}"

keep="$config_home/$name"
targets=(
  "$data_home/$name"
  "$state_home/$name"
  "$cache_home/$name"
)

printf 'Config CONSERVÉE : %s\n' "$keep"
[ -d "$keep" ] || printf '  (attention : ce dossier n’existe pas)\n'

# --- Rassembler ce qui existe réellement ------------------------------------
existing=()
printf '\nÀ supprimer :\n'
for t in "${targets[@]}"; do
  if [ -e "$t" ]; then
    existing+=("$t")
    printf '  - %s\n' "$t"
  else
    printf '  - %s  (absent)\n' "$t"
  fi
done

if [ "${#existing[@]}" -eq 0 ]; then
  printf '\nRien à supprimer.\n'
  exit 0
fi

if [ "$dry_run" -eq 1 ]; then
  printf '\n--dry-run : aucune suppression effectuée.\n'
  exit 0
fi

if [ "$assume_yes" -ne 1 ]; then
  printf '\nSupprimer ces %d dossier(s) ? [o/N] ' "${#existing[@]}"
  read -r reply
  case "$reply" in
    o | O | y | Y) ;;
    *) printf 'Annulé.\n'; exit 0 ;;
  esac
fi

for t in "${existing[@]}"; do
  rm -rf -- "$t"
  printf 'supprimé : %s\n' "$t"
done

printf '\nTerminé. Relance `NVIM_APPNAME=%s nvim` pour une installation propre.\n' "$name"
