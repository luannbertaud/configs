-- Plugin : sainnhe/everforest — LE thème de la config (unique et explicite).
-- Chargé quand : au démarrage, priorité haute (couleurs prêtes immédiatement).
return {
    'sainnhe/everforest',
    lazy = false,
    priority = 1000,
    config = function()
        vim.g.everforest_background = 'medium' -- 'hard' | 'medium' | 'soft'
        vim.g.everforest_enable_italic = true
        vim.g.everforest_better_performance = 1 -- thème allégé (moins de groupes)
        vim.g.everforest_transparent_background = 0 -- 0 = fond opaque
        vim.cmd.colorscheme('everforest')
    end,
}
