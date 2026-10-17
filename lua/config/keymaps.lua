-- ============================================================================
--  Raccourcis GLOBAUX
-- ----------------------------------------------------------------------------
--  Règle appliquée dans toute la config :
--    - ici            : raccourcis toujours disponibles (éditeur, fenêtres…).
--    - fichier plugin : raccourcis qui CHARGENT un plugin paresseux (via
--                       `keys = {…}` dans le spec — voir finder/fzf-lua.lua).
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

-- Explorateur de fichiers (neo-tree.nvim) -----------------------------------
-- Les raccourcis <leader>e / <leader>E vivent dans lua/plugins/ui/neo-tree.lua
-- (`keys`), car ils doivent CHARGER le plugin : c'est la règle de la config.

-- Terminal intégré (sans plugin) ---------------------------------------------
-- `<leader>tt` ouvre/ferme un terminal en bas (hauteur 15) ; la coquille reste
-- vivante entre deux ouvertures. `<Esc><Esc>` rend la main à Neovim : ensuite
-- les raccourcis habituels marchent (<C-hjkl>, <leader>q pour fermer…).
local term = { buf = nil, win = nil }

local function toggle_terminal()
    -- Déjà ouvert → on masque la fenêtre (le processus continue de tourner).
    if term.win and vim.api.nvim_win_is_valid(term.win) then
        vim.api.nvim_win_hide(term.win)
        term.win = nil
        return
    end
    if term.buf and vim.api.nvim_buf_is_valid(term.buf) then
        vim.cmd('botright split') -- réaffiche le terminal existant
        vim.api.nvim_win_set_buf(0, term.buf)
    else
        vim.cmd('botright new')
        term.buf = vim.api.nvim_get_current_buf()
        vim.fn.jobstart(vim.o.shell, { term = true }) -- démarre la coquille
        vim.bo[term.buf].bufhidden = 'hide'
    end
    term.win = vim.api.nvim_get_current_win()
    vim.wo.winfixheight = true
    vim.cmd('resize 15')
    vim.cmd('startinsert')
end

map('n', '<leader>tt', toggle_terminal, { desc = 'Terminal : ouvrir/fermer (bas)' })
map('t', '<Esc><Esc>', '<C-\\><C-n>', { desc = 'Terminal : revenir en mode normal' })

-- Formatage (conform.nvim) --------------------------------------------------
map(
    { 'n', 'v' },
    '<leader>f',
    function() require('conform').format({ async = true, lsp_format = 'fallback' }) end,
    { desc = 'Formater le buffer' }
)

-- Commentaires TODO / FIX ---------------------------------------------------
map('n', ']t', function() require('todo-comments').jump_next() end, { desc = 'Todo suivant' })
map('n', '[t', function() require('todo-comments').jump_prev() end, { desc = 'Todo précédent' })
