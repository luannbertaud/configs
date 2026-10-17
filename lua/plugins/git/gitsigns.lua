-- Plugin : lewis6991/gitsigns.nvim — signes Git dans la marge.
-- Chargé quand : à l'ouverture d'un fichier ; les raccourcis le chargent aussi.
return {
    'lewis6991/gitsigns.nvim',
    event = { 'BufReadPre', 'BufNewFile' },
    opts = {},
    keys = {
        { ']c', function() require('gitsigns').nav_hunk('next') end, desc = 'Git : hunk suivant' },
        { '[c', function() require('gitsigns').nav_hunk('prev') end, desc = 'Git : hunk précédent' },
        { '<leader>gs', function() require('gitsigns').stage_hunk() end, desc = 'Git : indexer le hunk' },
        { '<leader>gr', function() require('gitsigns').reset_hunk() end, desc = 'Git : annuler le hunk' },
        { '<leader>gp', function() require('gitsigns').preview_hunk() end, desc = 'Git : aperçu du hunk' },
        {
            '<leader>gb',
            function() require('gitsigns').blame_line({ full = true }) end,
            desc = 'Git : blame de la ligne',
        },
    },
}
