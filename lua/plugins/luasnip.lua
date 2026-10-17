-- Plugin : L3MON4D3/LuaSnip — moteur de snippets (utilisé par blink.cmp).
-- Dépendance : friendly-snippets (voir plugins/friendly-snippets.lua).
-- Chargé quand : avec blink.cmp (dépendance).
return {
  'L3MON4D3/LuaSnip',
  dependencies = { 'rafamadriz/friendly-snippets' },
  config = function()
    -- Charge les snippets « VS Code » fournis par friendly-snippets.
    require('luasnip.loaders.from_vscode').lazy_load()
  end,
}
