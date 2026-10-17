-- ============================================================================
--  Autocmds — automatismes, volontairement peu nombreux et explicites.
-- ============================================================================

-- 1) Flash bref du texte copié ----------------------------------------------
vim.api.nvim_create_autocmd('TextYankPost', {
    desc = 'Surligne brièvement le texte copié',
    group = vim.api.nvim_create_augroup('nvim_rc_yank', { clear = true }),
    callback = function() vim.hl.on_yank({ timeout = 150 }) end,
})

-- 2) Rouvrir un fichier à l'endroit où on l'avait quitté --------------------
vim.api.nvim_create_autocmd('BufReadPost', {
    desc = 'Restaure la dernière position du curseur',
    group = vim.api.nvim_create_augroup('nvim_rc_lastplace', { clear = true }),
    callback = function()
        local mark = vim.api.nvim_buf_get_mark(0, '"')
        local lcount = vim.api.nvim_buf_line_count(0)
        if mark[1] > 0 and mark[1] <= lcount then pcall(vim.api.nvim_win_set_cursor, 0, mark) end
    end,
})

-- 3) Raccourcis LSP -----------------------------------------------------------
--    Plus rien ici : on s'appuie sur les DÉFAUTS de Neovim (gr*, K, ]d…) et le
--    seul ajout (`gd`, définition) est global dans config/keymaps.lua.
