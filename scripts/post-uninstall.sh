#!/bin/sh
# post_uninstall du paquet nvim : efface data/state/cache via nvim-clean.sh.
# Le gestionnaire retire ensuite le clone de ~/.config/nvim.
set -eu

ROOT=${CONFIGS_TARGET:-$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)}
NAME=$(basename "$ROOT")
exec bash "$ROOT/scripts/nvim-clean.sh" "$NAME" -y
