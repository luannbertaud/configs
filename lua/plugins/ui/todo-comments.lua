-- Plugin : folke/todo-comments.nvim — surligne TODO / FIX / HACK… et navigue.
-- Dépendance : plenary.nvim (recherche des todos).
-- Chargé quand : VeryLazy, ou au raccourci ]t / [t ci-dessous.
return {
    'folke/todo-comments.nvim',
    dependencies = { 'nvim-lua/plenary.nvim' },
    event = 'VeryLazy',
    keys = {
        { ']t', function() require('todo-comments').jump_next() end, desc = 'TODO suivant' },
        { '[t', function() require('todo-comments').jump_prev() end, desc = 'TODO précédent' },
    },
    opts = {},
}
