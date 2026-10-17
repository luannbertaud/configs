-- Plugin : nvim-mini/mini.nvim — modules mini.pairs, mini.icons et mini.animate.
--   mini.pairs   : ferme/ouvre automatiquement (), [], {}, "", ''.
--   mini.icons   : jeu d'icônes centralisé (utilisé par neo-tree, lualine, fzf-lua…).
--   mini.animate : fait « glisser » les grands défilements (zt/zz/zb,
--                  <C-u>/<C-d>, <C-f>/<C-b>…) au lieu de sauter d'un coup.
-- Chargé quand : au démarrage (les icônes doivent être prêtes pour l'UI).
return {
    'nvim-mini/mini.nvim',
    lazy = false,
    config = function()
        require('mini.pairs').setup()
        require('mini.icons').setup()
        -- Fournit aussi les icônes aux plugins qui attendent nvim-web-devicons.
        require('mini.icons').mock_nvim_web_devicons()

        -- Animations (mini.animate). La DURÉE est réglable juste en dessous :
        -- `duration` en millisecondes, durée totale fixe quel que soit le saut.
        -- Baisse-la (ex. 100) pour plus sec, monte-la (ex. 300) pour plus fluide.
        local animate = require('mini.animate')
        local duration = 100

        -- Un « cran » de molette fait 'mousescroll' lignes (ver:3 par défaut).
        local mouse_lines = tonumber(vim.o.mousescroll:match('ver:(%d+)')) or 3

        animate.setup({
            scroll = {
                enable = true, -- zt/zz/zb, <C-u>/<C-d>, <C-f>/<C-b>…
                timing = animate.gen_timing.linear({ duration = duration, unit = 'total' }),
                -- On n'anime PAS les sauts ≤ un cran de molette : un scroll souris
                -- reste donc net, sans le tremblement dû au chevauchement des
                -- animations (chaque cran lançait la sienne). Les gros sauts clavier
                -- (<C-u>/<C-d>, <C-f>/<C-b>, zt/zz/zb) restent animés.
                subscroll = animate.gen_subscroll.equal({
                    predicate = function(total_scroll) return total_scroll > mouse_lines end,
                    max_output_steps = 60, -- fluidité (nb d'étapes intermédiaires)
                }),
            },
            cursor = {
                enable = false, -- déplacement du curseur
                timing = animate.gen_timing.linear({ duration = duration, unit = 'total' }),
            },
            open = { enable = false }, -- ouverture de fenêtre
            close = { enable = false }, -- fermeture de fenêtre
            resize = { enable = false }, -- redimensionnement
        })
    end,
}
