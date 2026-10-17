-- ============================================================================
--  Autocmds — automatismes, volontairement peu nombreux et explicites.
-- ============================================================================

-- 1) Flash bref du texte copié ----------------------------------------------
vim.api.nvim_create_autocmd('TextYankPost', {
  desc = 'Surligne brièvement le texte copié',
  group = vim.api.nvim_create_augroup('nvim_rc_yank', { clear = true }),
  callback = function()
    vim.hl.on_yank({ timeout = 150 })
  end,
})

-- 2) Rouvrir un fichier à l'endroit où on l'avait quitté --------------------
vim.api.nvim_create_autocmd('BufReadPost', {
  desc = 'Restaure la dernière position du curseur',
  group = vim.api.nvim_create_augroup('nvim_rc_lastplace', { clear = true }),
  callback = function()
    local mark = vim.api.nvim_buf_get_mark(0, '"')
    local lcount = vim.api.nvim_buf_line_count(0)
    if mark[1] > 0 and mark[1] <= lcount then
      pcall(vim.api.nvim_win_set_cursor, 0, mark)
    end
  end,
})

-- 3) Quand un serveur LSP s'attache à un buffer -----------------------------
--    On pose ici les raccourcis buffer-local liés au LSP.
vim.api.nvim_create_autocmd('LspAttach', {
  desc = 'Raccourcis LSP au moment de l’attachement',
  group = vim.api.nvim_create_augroup('nvim_rc_lsp_attach', { clear = true }),
  callback = function(event)
    local function map(keys, func, desc)
      vim.keymap.set('n', keys, func, { buffer = event.buf, desc = 'LSP: ' .. desc })
    end

    map('gd', vim.lsp.buf.definition, 'Définition')
    map('gD', vim.lsp.buf.declaration, 'Déclaration')
    map('gr', vim.lsp.buf.references, 'Références')
    map('gi', vim.lsp.buf.implementation, 'Implémentation')
    map('K', vim.lsp.buf.hover, 'Info au survol')
    map('<leader>rn', vim.lsp.buf.rename, 'Renommer le symbole')
    map('<leader>ca', vim.lsp.buf.code_action, 'Action de code')
    map('<leader>d', vim.diagnostic.open_float, 'Diagnostic (fenêtre flottante)')
    map('[d', function() vim.diagnostic.jump({ count = -1 }) end, 'Diagnostic précédent')
    map(']d', function() vim.diagnostic.jump({ count = 1 }) end, 'Diagnostic suivant')
  end,
})
