-- Plugin : mfussenegger/nvim-lint — lance des linters externes, affichés en diagnostics.
-- Chargé quand : à l'ouverture d'un fichier.
return {
    'mfussenegger/nvim-lint',
    event = { 'BufReadPre', 'BufNewFile' },
    config = function()
        local lint = require('lint')

        -- Linters par type de fichier (les binaires viennent de Mason / du système).
        lint.linters_by_ft = {
            python = { 'ruff' },
            sh = { 'shellcheck' },
            bash = { 'shellcheck' },
            javascript = { 'eslint' },
            typescript = { 'eslint' },
            javascriptreact = { 'eslint' },
            typescriptreact = { 'eslint' },
        }

        -- Déclenche le lint aux moments pertinents.
        vim.api.nvim_create_autocmd({ 'BufWritePost', 'BufReadPost', 'InsertLeave' }, {
            desc = 'Lance les linters du buffer',
            group = vim.api.nvim_create_augroup('nvim_rc_lint', { clear = true }),
            callback = function() require('lint').try_lint() end,
        })
    end,
}
