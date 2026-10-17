-- Plugin : stevearc/oil.nvim — éditer le système de fichiers comme un buffer.
-- Chargé quand : au démarrage (l'auteur déconseille le lazy-loading d'oil).
-- Le raccourci `-` est défini dans config/keymaps.lua.
return {
  'stevearc/oil.nvim',
  lazy = false,
  opts = {
    default_file_explorer = true,       -- `nvim .` ouvre oil
    view_options = { show_hidden = true },
  },
}
