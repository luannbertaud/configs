-- Plugin : nvim-treesitter/nvim-treesitter — coloration / indentation syntaxiques.
-- Branche `main` (réécriture pour Neovim 0.12). NE SUPPORTE PAS le lazy-loading.
-- Pré-requis EXTERNES : tree-sitter CLI (>= 0.26.1, PAS via npm), un compilateur
--   C, `tar` et `curl`.
-- Chargé quand : au démarrage ; `build` met à jour les parsers.
return {
    'nvim-treesitter/nvim-treesitter',
    branch = 'main',
    lazy = false,
    build = ':TSUpdate',
    config = function()
        require('nvim-treesitter').setup()

        -- Installe les parsers manquants (no-op si déjà présents).
        require('nvim-treesitter').install({
            'lua',
            'python',
            'javascript',
            'typescript',
            'tsx',
            'html',
            'css',
            'json',
            'yaml',
            'bash',
            'markdown',
            'markdown_inline',
            'regex',
            'query',
            'vim',
            'vimdoc',
            'diff',
            'gitcommit',
        })

        -- Coloration TS + indentation, activées pour chaque buffer.
        vim.api.nvim_create_autocmd('FileType', {
            desc = 'Démarre treesitter (coloration) et l’indentation',
            group = vim.api.nvim_create_augroup('nvim_rc_treesitter', { clear = true }),
            callback = function()
                pcall(vim.treesitter.start)
                vim.bo.indentexpr = "v:lua.require'nvim-treesitter'.indentexpr()"
            end,
        })
    end,
}
