-- Plugin : williamboman/mason-lspconfig.nvim — pont entre Mason et les serveurs LSP.
-- `ensure_installed` = liste EXPLICITE des serveurs LSP à installer.
-- automatic_enable = false : l'activation est faite explicitement, dans
--   plugins/lsp/nvim-lspconfig.lua (rien de caché).
-- Chargé quand : au démarrage.
return {
    'williamboman/mason-lspconfig.nvim',
    dependencies = { 'williamboman/mason.nvim' },
    opts = {
        ensure_installed = {
            'lua_ls', -- Lua
            'gopls', -- Go
            'basedpyright', -- Python
            'vtsls', -- TypeScript / JavaScript
            'bashls', -- Shell
            'html', -- HTML
            'cssls', -- CSS
            'jsonls', -- JSON
            'yamlls', -- YAML
            'emmet_language_server', -- Emmet (HTML/CSS)
        },
        automatic_enable = false,
    },
}
