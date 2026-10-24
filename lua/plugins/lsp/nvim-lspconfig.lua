-- Plugin : neovim/nvim-lspconfig — fournit les DÉFINITIONS des serveurs LSP.
-- Il ne fait plus l'activation lui-même : on utilise l'API native de Neovim 0.11+.
-- Chargé quand : au démarrage (avant tout attachement de serveur).
return {
    'neovim/nvim-lspconfig',
    dependencies = {
        'williamboman/mason.nvim',
        'williamboman/mason-lspconfig.nvim',
    },
    config = function()
        -- Capacités de complétion LSP fournies par blink.cmp.
        vim.lsp.config('*', {
            capabilities = require('blink.cmp').get_lsp_capabilities(),
        })

        -- Réglage d'un serveur précis : lua_ls connaît les API Neovim.
        vim.lsp.config('lua_ls', {
            settings = {
                Lua = {
                    diagnostics = { globals = { 'vim' } },
                    workspace = { checkThirdParty = false },
                    telemetry = { enable = false },
                },
            },
        })

        -- Activation EXPLICITE des serveurs (API native Neovim 0.12).
        vim.lsp.enable({
            'lua_ls',
            'gopls',
            'basedpyright',
            'vtsls',
            'bashls',
            'html',
            'cssls',
            'jsonls',
            'yamlls',
            'emmet_language_server',
        })
    end,
}
