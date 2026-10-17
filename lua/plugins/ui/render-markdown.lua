-- Plugin : MeanderingProgrammer/render-markdown.nvim — rendu Markdown DANS Neovim
-- (titres, listes, cases à cocher, tableaux, blocs de code…), sans navigateur.
-- Chargé quand : ouverture d'un fichier Markdown (`ft`).
-- Dépendances : nvim-treesitter (parsers `markdown`/`markdown_inline`, déjà
--   installés) + un fournisseur d'icônes (mini.icons, déjà là).
-- Le fichier reste du texte brut : `:RenderMarkdown toggle` montre la source.
return {
    'MeanderingProgrammer/render-markdown.nvim',
    ft = { 'markdown' },
    dependencies = { 'nvim-treesitter/nvim-treesitter' },
    opts = {},
}
