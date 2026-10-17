-- Plugin : nvim-neo-tree/neo-tree.nvim — arbre de fichiers dans un panneau à gauche.
-- Chargé quand : au premier appui sur <leader>te / <leader>tE (voir `keys`).
-- Dépendances : nui.nvim (interface) + plenary.nvim (backend).
--   Les icônes viennent de mini.icons (« mock nvim-web-devicons »), donc
--   nvim-web-devicons est inutile ici.
return {
    'nvim-neo-tree/neo-tree.nvim',
    branch = 'v3.x',
    dependencies = {
        'nvim-lua/plenary.nvim',
        'MunifTanjim/nui.nvim',
    },
    -- `keys` / `cmd` = laziness ciblée : ces déclencheurs chargent le plugin.
    -- v3.x n'expose qu'une commande unique : `:Neotree <action>` (ex. `:Neotree toggle`).
    cmd = 'Neotree',
    keys = {
        { '<leader>te', '<cmd>Neotree toggle<cr>', desc = 'Toggle : explorateur de fichiers' },
        { '<leader>tE', '<cmd>Neotree reveal<cr>', desc = 'Toggle : révéler le fichier courant' },
    },
    opts = {
        close_if_last_window = false, -- garder Neovim ouvert si l'arbre est la dernière fenêtre
        enable_git_status = false, -- doublon inutile avec gitsigns
        enable_diagnostics = false, -- gérés globalement par vim.diagnostic
        filesystem = {
            follow_current_file = { enabled = true }, -- suit le fichier du buffer courant
            use_libuv_file_watcher = true, -- rafraîchit l'arbre sur changement disque
            hijack_netrw_behavior = 'open_default', -- `nvim .` ouvre neo-tree (une fois le plugin chargé)
            filtered_items = {
                hide_dotfiles = false, -- montre les fichiers cachés (comme oil le faisait)
                hide_gitignored = false,
            },
        },
        default_component_configs = {
            indent = { with_expanders = true }, -- flèches de pliage des dossiers
        },
        window = {
            width = 30,
            mappings = {
                ['<space>'] = 'none', -- libère <space> (leader) à l'intérieur de l'arbre
            },
        },
    },
}
