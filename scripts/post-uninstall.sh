#!/bin/sh
# Adaptateur pour le gestionnaire `configs` : réutilise install.sh --uninstall,
# qui retire la ligne de source sans supprimer le clone (le gestionnaire le fait).
set -eu

ROOT=${CONFIGS_TARGET:-$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)}
CONFIGS_TARGET="$ROOT" exec sh "$ROOT/install.sh" --uninstall
