-- Plugin : folke/which-key.nvim — affiche les raccourcis disponibles.
-- Chargé quand : VeryLazy.
-- Rôle ici : nommer les GROUPES du leader et DOCUMENTER en français les défauts
-- de Neovim qu'on utilise (LSP, diagnostics…), pour que tout soit lisible dans
-- le popup (et retrouvable via <leader>fk, le picker de keymaps).
return {
    'folke/which-key.nvim',
    event = 'VeryLazy',
    opts = {
        spec = {
            -- Groupes sous <leader> -------------------------------------------------
            { '<leader>b', group = 'Buffers' },
            { '<leader>c', group = 'Code' },
            { '<leader>f', group = 'Find / recherche' },
            { '<leader>g', group = 'Git' },
            { '<leader>t', group = 'Toggles' },
            { '<leader>y', group = 'Copier (chemins)' },

            -- Défauts Neovim documentés (aucun mapping créé, juste des libellés) ----
            { 'g', group = 'Aller / g-commandes' },
            { 'gr', group = 'LSP' },
            { 'grn', desc = 'LSP : renommer le symbole' },
            { 'gra', desc = 'LSP : action de code' },
            { 'grr', desc = 'LSP : références' },
            { 'gri', desc = 'LSP : implémentation' },
            { 'grt', desc = 'LSP : définition de type' },
            { 'grx', desc = 'LSP : codelens' },
            { 'gO', desc = 'LSP : symboles du document' },
            { 'K', desc = 'LSP : info au survol' },
            { '<C-S>', desc = 'LSP : aide sur la signature' },
            { 'gq', desc = 'Formater (LSP / formatexpr)' },
            { 'gc', desc = 'Commenter (ligne / texte)' },
            { 'gx', desc = 'Ouvrir le lien du curseur' },
            { '[', group = 'Précédent' },
            { ']', group = 'Suivant' },
            { ']d', desc = 'Diagnostic suivant' },
            { '[d', desc = 'Diagnostic précédent' },
        },
    },
}
