-- Plugin : ibhagwan/fzf-lua — picker (fichiers, texte, buffers, aide…).
-- Dépendance EXTERNE : le binaire `fzf` (installé par le système, pas Mason).
-- Chargé quand : au premier usage d'un raccourci ci-dessous, ou de :FzfLua.
return {
    'ibhagwan/fzf-lua',
    cmd = 'FzfLua',
    keys = {
        { '<leader><leader>', function() require('fzf-lua').files() end, desc = 'Fichiers' },
        { '<leader>ff', function() require('fzf-lua').files() end, desc = 'Fichiers' },
        { '<leader>fg', function() require('fzf-lua').live_grep() end, desc = 'Chercher du texte' },
        { '<leader>fb', function() require('fzf-lua').buffers() end, desc = 'Buffers ouverts' },
        { '<leader>fr', function() require('fzf-lua').oldfiles() end, desc = 'Fichiers récents' },
        { '<leader>fh', function() require('fzf-lua').help_tags() end, desc = 'Aide (help tags)' },
        {
            '<leader>cd',
            function() require('fzf-lua').diagnostics_document() end,
            desc = 'Code : diagnostics du buffer',
        },
        {
            '<leader>cD',
            function() require('fzf-lua').diagnostics_workspace() end,
            desc = 'Code : diagnostics du projet',
        },
        { '<leader>fk', function() require('fzf-lua').keymaps() end, desc = 'Raccourcis (keymaps)' },
    },
    opts = {},
}
