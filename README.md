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

## Étendre la configuration

Chaque plugin = **un fichier** dans `lua/plugins/<domaine>/`, avec un domaine
parmi : `theme`, `ui`, `editing`, `completion`, `lsp`, `git`, `finder`.

### Ajouter un plugin

1. Créer `lua/plugins/<domaine>/<nom>.lua` qui `return` **un** spec lazy.nvim.
2. Le charger paresseusement selon le besoin : `event`, `ft`, `cmd` ou `keys`
   (exemples déjà présents : `keys` dans fzf-lua/neo-tree, `ft` dans render-markdown).
3. Relancer Neovim, puis `:Lazy` pour l'installer.

```lua
-- lua/plugins/ui/mon-plugin.lua
return {
  'auteur/mon-plugin.nvim',
  event = 'VeryLazy',   -- ou ft = { 'markdown' }, cmd = { 'MonPlugin' }, keys = { … }
  opts = {},
}
```

> ⚠️ Un **nouveau sous-dossier** doit être ajouté à la main dans
> `lua/config/lazy.lua` (`spec = { { import = 'plugins.mondomaine' }, … }`) :
> lazy.nvim n'importe pas les sous-dossiers tout seul.

### Ajouter un thème

```lua
-- lua/plugins/theme/montheme.lua
return {
  'auteur/montheme.nvim',
  lazy = false,      -- un thème se charge au démarrage…
  priority = 1000,   -- …et en premier (évite un flash de couleurs par défaut)
  config = function()
    require('montheme').setup({})
    vim.cmd.colorscheme('montheme')
  end,
}
```

Puis supprime `lua/plugins/theme/everforest.lua` (sinon les deux thèmes se chargent).
Optionnel : mets à jour `install.colorscheme` dans `lua/config/lazy.lua`.

### Ajouter un serveur LSP

Deux endroits, dans l'ordre — avec le **même nom de serveur** :

1. `lua/plugins/lsp/mason-lspconfig.lua` → ajouter le nom à `ensure_installed`.
2. `lua/plugins/lsp/nvim-lspconfig.lua` → ajouter le nom à `vim.lsp.enable({ … })`
   (et une entrée `vim.lsp.config('nom', { … })` pour le paramétrer).

Relance, puis `:Mason` pour vérifier l'installation. Les raccourcis LSP
(`gd`, `K`, `<leader>rn`…) sont posés automatiquement à l'attachement du serveur
(voir `lua/config/autocmds.lua`).

Les **formatters / linters** suivent le même principe :
`lua/plugins/lsp/conform.lua` (`formatters_by_ft`) et
`lua/plugins/lsp/nvim-lint.lua` (`linters_by_ft`) pour l'usage, et
`lua/plugins/lsp/mason.lua` pour installer les binaires.

## Structure

```
init.lua            point d'entrée (lancement + pré-requis en commentaire)
lua/config/         options · keymaps · autocmds · lazy
lua/plugins/        1 fichier = 1 plugin, rangé par domaine :
  theme/            colorscheme (everforest)
  ui/               lualine · which-key · neo-tree · todo-comments · render-markdown
  editing/          mini.nvim (paires, icônes, défilement animé) · nvim-treesitter
  completion/       blink.cmp · LuaSnip · friendly-snippets
  lsp/              mason · mason-lspconfig · nvim-lspconfig · conform · nvim-lint · lazydev
  git/              gitsigns
  finder/           fzf-lua
lazy-lock.json      versions figées des plugins (versionné)
.stylua.toml        style Lua (appliqué par conform.nvim au format-on-save)
scripts/            outillage (nettoyage)
```
