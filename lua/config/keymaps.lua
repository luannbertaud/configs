-- ============================================================================
--  Raccourcis GLOBAUX
-- ----------------------------------------------------------------------------
--  Règle appliquée dans toute la config :
--    - ici            : raccourcis toujours disponibles (éditeur, fenêtres…).
--    - fichier plugin : raccourcis qui CHARGENT un plugin paresseux (via
--                       `keys = {…}` dans le spec — voir finder/fzf-lua.lua).
--    - LSP            : on s'appuie sur les DÉFAUTS de Neovim (gr*, K, ]d…),
--                       documentés en français dans lua/plugins/ui/which-key.lua.
--
--  Voir tous les keymaps : <leader>fk (picker fzf-lua) ou <leader> (which-key).
-- ============================================================================

local map = vim.keymap.set

-- Confort -------------------------------------------------------------------
map('n', '<leader>w', '<cmd>write<cr>', { desc = 'Enregistrer le fichier' })
map('n', '<leader>q', '<cmd>quit<cr>', { desc = 'Fermer la fenêtre' })
map('n', '<leader>x', ':', { desc = 'Saisir une commande (:)' })

-- Navigation entre fenêtres (splits) ----------------------------------------
map('n', '<C-h>', '<C-w>h', { desc = 'Fenêtre : gauche' })
map('n', '<C-j>', '<C-w>j', { desc = 'Fenêtre : bas' })
map('n', '<C-k>', '<C-w>k', { desc = 'Fenêtre : haut' })
map('n', '<C-l>', '<C-w>l', { desc = 'Fenêtre : droite' })

-- Buffers -------------------------------------------------------------------
map('n', 'H', '<cmd>bprevious<cr>', { desc = 'Buffer précédent' })
map('n', 'L', '<cmd>bnext<cr>', { desc = 'Buffer suivant' })
map('n', '<leader>bd', '<cmd>bdelete<cr>', { desc = 'Fermer le buffer' })
map('n', '<leader>bn', '<cmd>bnext<cr>', { desc = 'Buffer suivant' })
map('n', '<leader>bp', '<cmd>bprevious<cr>', { desc = 'Buffer précédent' })

-- Copier le chemin du fichier courant ----------------------------------------
-- `ya` = chemin ABSOLU ; `yr` = chemin RELATIF à la racine du projet (dossier
-- git, sinon cwd). Le résultat part dans le presse-papier (+).
local function copy_path(absolute)
    local abs = vim.api.nvim_buf_get_name(0)
    local path = abs
    if not absolute then
        local root = vim.fs.root(0, { '.git', '.hg', 'pyproject.toml', 'package.json' }) or vim.uv.cwd()
        path = abs:sub(1, #root) == root and abs:sub(#root + 2) or vim.fn.fnamemodify(abs, ':.')
    end
    vim.fn.setreg('+', path)
    vim.notify('Copié : ' .. path)
end

map('n', '<leader>ya', function() copy_path(true) end, { desc = 'Copier : chemin absolu' })
map('n', '<leader>yr', function() copy_path(false) end, { desc = 'Copier : chemin relatif (racine)' })

-- LSP : complément aux DÉFAUTS de Neovim --------------------------------------
-- Neovim fournit déjà, SANS configuration (donc rien à réécrire ici) :
--   grn renommer · gra action de code · grr références · gri implémentation
--   grt type · grx codelens · gO symboles · K survol · <C-S> signature
--   ]d / [d diagnostics · <C-W>d diagnostic flottant · gq formater (LSP)
--   gc commentaire · gx ouvrir le lien du curseur
-- Seule la DÉFINITION n'a pas de défaut : on l'ajoute (habitude forte).
-- Défini GLOBALEMENT pour rester visible dans <leader>fk (picker de keymaps).
map('n', 'gd', vim.lsp.buf.definition, { desc = 'LSP : définition' })

-- Diagnostics ----------------------------------------------------------------
-- On retire les défauts <C-W>d / <C-W><C-D> (peu pratiques) au profit de
-- <leader>cw. La LISTE des diagnostics : <leader>cd (buffer) / <leader>cD (projet).
pcall(vim.keymap.del, 'n', '<C-W>d')
pcall(vim.keymap.del, 'n', '<C-W><C-D>')
map('n', '<leader>cw', vim.diagnostic.open_float, { desc = 'Code : diagnostic (fenêtre flottante)' })

-- Terminal intégré (sans plugin) ---------------------------------------------
-- `<leader>tt` ouvre/ferme un terminal en bas (hauteur 15) ; la coquille reste
-- vivante entre deux ouvertures. `<Esc><Esc>` rend la main à Neovim : ensuite
-- les raccourcis habituels marchent (<C-hjkl>, <leader>q pour fermer…).
-- L'explorateur est sous `<leader>te` / `<leader>tE` (voir ui/neo-tree.lua) et
-- le formatage sous `<leader>cf` (voir lsp/conform.lua).
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
