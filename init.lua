-- ============================================================================
--  nvim — point d'entrée
-- ----------------------------------------------------------------------------
--  LANCEMENT : configuration principale de Neovim (dans ~/.config/nvim).
--      nvim
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
