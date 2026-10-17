-- Plugin : folke/lazydev.nvim — améliore lua_ls pour éditer CETTE config
--   (complétion correcte des API vim.*, runtime, plugins…).
-- Chargé quand : fichier de type lua uniquement.
return {
  'folke/lazydev.nvim',
  ft = 'lua',
  opts = {
    library = {
      { path = '${3rd}/luv/library', words = { 'vim%.uv' } },
    },
  },
}
