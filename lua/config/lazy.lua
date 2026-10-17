-- ============================================================================
--  Bootstrap + configuration de lazy.nvim
-- ----------------------------------------------------------------------------
--  lazy.nvim est le seul plugin installé hors de lua/plugins/ : on le clone ici
--  de façon VISIBLE, dans le data dir propre à nvim-rc.
--  Emplacement des plugins : ~/.local/share/nvim-rc/lazy/
-- ============================================================================

local lazypath = vim.fn.stdpath('data') .. '/lazy/lazy.nvim'

-- Amorce : clone lazy.nvim s'il est absent (premier lancement seulement).
if not (vim.uv or vim.loop).fs_stat(lazypath) then
  vim.fn.system({
    'git',
    'clone',
    '--filter=blob:none',
    '--branch=stable',
    'https://github.com/folke/lazy.nvim.git',
    lazypath,
  })
  if vim.v.shell_error ~= 0 then
    vim.api.nvim_echo({
      { 'Échec du clone de lazy.nvim (git / connexion ?).\n', 'ErrorMsg' },
      { 'Appuie sur une touche pour quitter…' },
    }, true, {})
    vim.fn.getchar()
    os.exit(1)
  end
end

vim.opt.rtp:prepend(lazypath)

require('lazy').setup({
  -- Chaque fichier de lua/plugins/ contient UN SEUL plugin (voir son en-tête).
  spec = {
    { import = 'plugins' },
  },
  -- Thème appliqué pendant l'installation (évite le flash blanc).
  install = { colorscheme = { 'everforest' } },
  -- Pas de mise à jour automatique : tu décides (`:Lazy update`), pour rester
  -- déterministe vis-à-vis du lockfile.
  checker = { enabled = false },
  change_detection = { notify = false },
})
