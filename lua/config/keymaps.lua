-- ============================================================================
--  Raccourcis GLOBAUX
-- ----------------------------------------------------------------------------
--  Règle appliquée dans toute la config :
--    - ici            : raccourcis toujours disponibles (éditeur, fenêtres…).
--    - fichier plugin : raccourcis qui CHARGENT un plugin paresseux (via
--                       `keys = {…}` dans le spec — voir fzf-lua.lua, gitsigns).
--    - autocmds.lua   : raccourcis LSP posés à l'attachement du serveur.
-- ============================================================================

local map = vim.keymap.set

-- Confort -------------------------------------------------------------------
map('n', '<leader>w', '<cmd>write<cr>', { desc = 'Enregistrer le fichier' })
map('n', '<leader>q', '<cmd>quit<cr>', { desc = 'Fermer la fenêtre' })

-- Navigation entre fenêtres (splits) ----------------------------------------
map('n', '<C-h>', '<C-w>h', { desc = 'Fenêtre : gauche' })
map('n', '<C-j>', '<C-w>j', { desc = 'Fenêtre : bas' })
map('n', '<C-k>', '<C-w>k', { desc = 'Fenêtre : haut' })
map('n', '<C-l>', '<C-w>l', { desc = 'Fenêtre : droite' })

-- Buffers -------------------------------------------------------------------
map('n', '<S-l>', '<cmd>bnext<cr>', { desc = 'Buffer suivant' })
map('n', '<S-h>', '<cmd>bprevious<cr>', { desc = 'Buffer précédent' })
map('n', '<leader>bd', '<cmd>bdelete<cr>', { desc = 'Fermer le buffer' })

-- Explorateur de fichiers (oil.nvim) ----------------------------------------
-- `-`        : ouvre oil sur le dossier PARENT (façon vim-vinegar).
-- `<leader>e`: ouvre oil dans une fenêtre FLOTTANTE sur le dossier courant.
-- (oil s'édite comme un buffer : modifie les noms/chemins puis `:w` pour appliquer.)
map('n', '-', '<cmd>Oil<cr>', { desc = 'Oil : dossier parent' })
map('n', '<leader>e', function()
  require('oil').toggle_float()
end, { desc = 'Explorateur de fichiers (oil)' })

-- Formatage (conform.nvim) --------------------------------------------------
map({ 'n', 'v' }, '<leader>f', function()
  require('conform').format({ async = true, lsp_format = 'fallback' })
end, { desc = 'Formater le buffer' })

-- Commentaires TODO / FIX ---------------------------------------------------
map('n', ']t', function() require('todo-comments').jump_next() end, { desc = 'Todo suivant' })
map('n', '[t', function() require('todo-comments').jump_prev() end, { desc = 'Todo précédent' })
