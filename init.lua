-- ============================================================================
--  nvim-rc — point d'entrée
-- ----------------------------------------------------------------------------
--  LANCEMENT : cette config est isolée, elle ne remplace PAS ~/.config/nvim.
--      NVIM_APPNAME=nvim-rc nvim
--  (l'env var fait pointer config/data/state/cache sur nvim-rc ; on peut en
--   faire un alias `nvrc` dans ~/.bashrc si souhaité).
--
--  PRÉ-REQUIS SYSTÈME (à installer via pacman/yay, hors Mason) :
--      neovim>=0.12  git  ripgrep  fd  fzf  wl-clipboard  curl  unzip  tar
--      un compilateur C (gcc/clang)  tree-sitter-cli (>=0.26.1, PAS npm)
--      node/npm (pour vtsls/prettier/eslint)  python3 (pour basedpyright/ruff)
--
--  Ordre de chargement EXPLICITE et volontaire :
--    1. options   -> réglages + touche leader (DOIT précéder lazy)
--    2. lazy      -> gestionnaire de plugins, charge lua/plugins/*
--    3. autocmds  -> automatismes explicites
--    4. keymaps   -> raccourcis globaux (après lazy, pour que tout existe)
-- ============================================================================

require('config.options')
require('config.lazy')
require('config.autocmds')
require('config.keymaps')
