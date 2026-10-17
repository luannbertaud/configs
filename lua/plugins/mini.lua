-- Plugin : nvim-mini/mini.nvim — modules mini.pairs et mini.icons.
--   mini.pairs : ferme/ouvre automatiquement (), [], {}, "", ''.
--   mini.icons : jeu d'icônes centralisé (utilisé par oil, lualine, fzf-lua…).
-- Chargé quand : au démarrage (les icônes doivent être prêtes pour l'UI).
return {
  'nvim-mini/mini.nvim',
  lazy = false,
  config = function()
    require('mini.pairs').setup()
    require('mini.icons').setup()
    -- Fournit aussi les icônes aux plugins qui attendent nvim-web-devicons.
    require('mini.icons').mock_nvim_web_devicons()
  end,
}
