-- Plugin : nvim-lualine/lualine.nvim — barre de statut.
-- Chargé quand : VeryLazy (dès que l'UI est prête, après le démarrage).
return {
  'nvim-lualine/lualine.nvim',
  event = 'VeryLazy',
  opts = {
    options = {
      theme = 'auto',            -- suit le thème actif (everforest)
      globalstatus = true,       -- une seule barre pour tout Neovim
      section_separators = '',   -- rendu sobre, sans séparateurs
      component_separators = '',
    },
  },
}
