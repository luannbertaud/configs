-- Plugin : williamboman/mason.nvim — installe des outils externes.
-- Emplacement : ~/.local/share/nvim/mason/ (bin/ et packages/).
--
-- NOTE IMPORTANTE : mason.nvim n'a PAS d'option `ensure_installed`.
-- On installe donc EXPLICITEMENT la liste `tools` ci-dessous au démarrage
-- (no-op si le paquet est déjà installé). Ce sont les formatters/linters NON-LSP ;
-- les serveurs LSP sont gérés séparément par mason-lspconfig.nvim.
-- Chargé quand : au démarrage, via mason-lspconfig (dépendance).

-- Outils NON-LSP à garantir (binaires utilisés par conform.nvim / nvim-lint).
local tools = {
    'stylua', -- Lua : format
    'ruff', -- Python : lint + fix
    'prettier', -- JS/TS, CSS, HTML, JSON, YAML, Markdown : format
    'eslint_d', -- JS/TS : lint (utilisé par nvim-lint)
    'shfmt', -- Shell : format
    'shellcheck', -- Shell : lint
}

return {
    'williamboman/mason.nvim',
    opts = {
        ui = { border = 'rounded' },
    },
    config = function(_, opts)
        require('mason').setup(opts)

        -- Installe les outils manquants (la liste `tools` est celle ci-dessus).
        local registry = require('mason-registry')
        registry.refresh(function()
            for _, name in ipairs(tools) do
                local ok, pkg = pcall(registry.get_package, name)
                if ok and not pkg:is_installed() then pkg:install() end
            end
        end)
    end,
}
