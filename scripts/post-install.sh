#!/bin/sh
# Adaptateur pour le gestionnaire `configs` : réutilise install.sh du clone,
# qui détecte le mode local et source aliases.sh sans le retélécharger.
set -eu

ROOT=${CONFIGS_TARGET:-$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)}
CONFIGS_TARGET="$ROOT" exec sh "$ROOT/install.sh"
