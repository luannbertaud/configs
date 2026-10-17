# nvim-rc

Configuration Neovim **isolée** : elle vit dans `~/.config/nvim-rc` et ne touche
pas à `~/.config/nvim`. Elle se lance via la variable d'environnement `NVIM_APPNAME`.

## Pré-requis

`neovim>=0.12`, `git`, `ripgrep`, `fd`, `fzf`, `wl-clipboard`, `curl`, `tar`,
`unzip`, un compilateur C, `tree-sitter-cli>=0.26.1` (pas npm), `node`/`npm`, `python3`.

## Lancer

```sh
NVIM_APPNAME=nvim-rc nvim
```

Alias pratique dans `~/.bashrc` :

```sh
alias nvrc='NVIM_APPNAME=nvim-rc nvim'
```

## Où vont les fichiers

| Rôle | Chemin | Supprimable |
|---|---|---|
| Config (à garder) | `~/.config/nvim-rc` | non |
| Plugins | `~/.local/share/nvim-rc/lazy` | oui |
| LSP / formatters (Mason) | `~/.local/share/nvim-rc/mason` | oui |
| Parsers treesitter | `~/.local/share/nvim-rc/site` | oui |
| État (undo, shada, logs) | `~/.local/state/nvim-rc` | oui |
| Cache | `~/.cache/nvim-rc` | oui |

## Installation & mises à jour

Au **premier lancement**, `lazy.nvim` s'installe seul puis installe les plugins ;
Mason installe ensuite les outils et serveurs LSP en arrière-plan (laisse
quelques secondes, relance une fois si besoin).

| Commande | Quand |
|---|---|
| `:Lazy` | voir l'état, installer ce qui manque |
| `:Lazy sync` | après un clone, ou tout (re)mettre à jour |
| `:Lazy restore` | **nouvelle machine** : fige les versions du `lazy-lock.json` |
| `:Lazy update` | mettre à jour les plugins (met à jour le lockfile) |
| `:Mason` | voir les outils / serveurs LSP installés ou manquants |
| `:MasonUpdate` | rafraîchir le catalogue Mason |
| `:checkhealth` | diagnostiquer (lazy, mason, lsp, treesitter…) |

Reproduire ailleurs : `git clone` → lancer → `:Lazy restore` → redémarrer.

## Repartir d'une installation propre

```sh
scripts/nvim-clean.sh nvim-rc            # demande confirmation
scripts/nvim-clean.sh nvim-rc --dry-run  # affiche sans supprimer
scripts/nvim-clean.sh nvim-rc -y         # sans confirmation
```

Supprime `~/.local/share/nvim-rc`, `~/.local/state/nvim-rc` et `~/.cache/nvim-rc`.
**Conserve** `~/.config/nvim-rc`. Ne touche pas au runtime dir partagé.

## Structure

```
init.lua          point d'entrée (lancement + pré-requis en commentaire)
lua/config/       options · keymaps · autocmds · lazy
lua/plugins/      1 fichier = 1 plugin (commenté)
lazy-lock.json    versions figées des plugins (versionné)
scripts/          outillage (nettoyage)
```
