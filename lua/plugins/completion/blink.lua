-- Plugin : saghen/blink.cmp — moteur de complétion.
-- Version ÉPINGLÉE sur la branche stable 1.x (la 2.x est en développement).
-- Chargé quand : au démarrage — ses capacités LSP doivent être connues AVANT
--   que les serveurs ne s'attachent (cf. plugins/lsp/nvim-lspconfig.lua).
return {
    'saghen/blink.cmp',
    version = '1.*',
    lazy = false,
    dependencies = { 'L3MON4D3/LuaSnip' },
    opts = {
        keymap = { preset = 'default' },
        appearance = { nerd_font_variant = 'mono' },
        completion = { documentation = { auto_show = true } },
        sources = { default = { 'lsp', 'path', 'snippets', 'buffer' } },
        snippets = { preset = 'luasnip' },
        fuzzy = { implementation = 'prefer_rust_with_warning' },
    },
}
