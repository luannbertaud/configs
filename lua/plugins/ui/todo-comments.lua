-- Plugin : folke/todo-comments.nvim — surligne TODO / FIX / HACK… et navigue.
-- Dépendance : plenary.nvim (recherche des todos).
-- Chargé quand : VeryLazy ; les raccourcis ]t / [t sont dans config/keymaps.lua.
return {
    'folke/todo-comments.nvim',
    dependencies = { 'nvim-lua/plenary.nvim' },
    event = 'VeryLazy',
    opts = {},
}
