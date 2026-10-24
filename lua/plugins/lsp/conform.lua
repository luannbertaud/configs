-- Plugin : stevearc/conform.nvim — formatage (formatters externes + LSP).
-- Chargé quand : avant l'écriture (format on save) ou au raccourci <leader>cf.
return {
    'stevearc/conform.nvim',
    event = { 'BufWritePre' },
    cmd = { 'ConformInfo' },
    keys = {
        {
            '<leader>cf',
            function() require('conform').format({ async = true, lsp_format = 'fallback' }) end,
            mode = { 'n', 'v' },
            desc = 'Code : formater le buffer',
        },
    },
    opts = {
        formatters_by_ft = {
            lua = { 'stylua' },
            python = { 'ruff_format' },
            go = { 'gofmt' },
            javascript = { 'prettier' },
            typescript = { 'prettier' },
            javascriptreact = { 'prettier' },
            typescriptreact = { 'prettier' },
            css = { 'prettier' },
            scss = { 'prettier' },
            html = { 'prettier' },
            json = { 'prettier' },
            yaml = { 'prettier' },
            markdown = { 'prettier' },
            sh = { 'shfmt' },
            bash = { 'shfmt' },
        },
        format_on_save = {
            timeout_ms = 500,
            lsp_format = 'fallback', -- utilise le LSP si aucun formatter n'est défini
        },
    },
}
