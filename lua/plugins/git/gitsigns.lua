-- Plugin : lewis6991/gitsigns.nvim — signes Git dans la marge.
-- Chargé quand : à l'ouverture d'un fichier ; les raccourcis le chargent aussi.
return {
    'lewis6991/gitsigns.nvim',
    event = { 'BufReadPre', 'BufNewFile' },
    opts = {},
    keys = {
        { ']c', function() require('gitsigns').nav_hunk('next') end, desc = 'Hunk Git suivant' },
        { '[c', function() require('gitsigns').nav_hunk('prev') end, desc = 'Hunk Git précédent' },
        { '<leader>hp', function() require('gitsigns').preview_hunk() end, desc = 'Aperçu du hunk' },
        { '<leader>hr', function() require('gitsigns').reset_hunk() end, desc = 'Annuler le hunk' },
        { '<leader>hs', function() require('gitsigns').stage_hunk() end, desc = 'Indexer (stage) le hunk' },
    },
}
