-- ============================================================================
--  Options de l'éditeur
--  Tout est écrit explicitement : aucun réglage « hérité » ou magique.
-- ============================================================================

-- Touche leader -------------------------------------------------------------
-- DOIT être défini avant le chargement des plugins (donc avant config.lazy).
vim.g.mapleader = ' '
vim.g.maplocalleader = ' '

local opt = vim.opt

-- Affichage -----------------------------------------------------------------
opt.number = true -- numéros de ligne absolus
opt.relativenumber = true -- numéros relatifs (saut sur N lignes)
opt.signcolumn = 'yes' -- colonne de signes toujours réservée
opt.cursorline = true -- surligne la ligne du curseur
opt.termguicolors = true -- couleurs 24 bits (requis par les thèmes)
opt.showmode = false -- le mode est déjà affiché par lualine
opt.wrap = false -- pas de retour à la ligne visuel
opt.scrolloff = 4 -- 4 lignes de contexte au-dessus/en dessous
opt.sidescrolloff = 8
opt.splitright = true -- nouveau split à droite…
opt.splitbelow = true -- …et en dessous
opt.laststatus = 3 -- une seule barre de statut globale (lualine)

-- Caractères invisibles -----------------------------------------------------
opt.list = true
opt.listchars = { tab = '» ', trail = '·', nbsp = '␣' }

-- Indentation ---------------------------------------------------------------
opt.expandtab = true -- Tab insère des espaces
opt.tabstop = 4 -- une tabulation = 4 espaces
opt.shiftwidth = 4 -- indentation = 4 espaces
opt.softtabstop = 4
opt.smartindent = true
opt.shiftround = true -- indente par multiples de shiftwidth

-- Recherche -----------------------------------------------------------------
opt.ignorecase = true -- insensible à la casse…
opt.smartcase = true -- …sauf quand on tape une majuscule
opt.incsearch = true -- recherche incrémentale
opt.hlsearch = true

-- Fichiers & historique -----------------------------------------------------
opt.undofile = true -- historique d'annulation persistant
opt.swapfile = false -- pas de fichier .swp
opt.backup = false
opt.updatetime = 250 -- délai CursorHold / écriture du swap
opt.timeoutlen = 300 -- délai d'attente des mappings

-- Édition -------------------------------------------------------------------
opt.mouse = 'a' -- souris active partout
opt.clipboard = 'unnamedplus' -- presse-papier système (via wl-clipboard)
opt.completeopt = { 'menuone', 'noselect' } -- complétion (voir plugin blink.cmp)
opt.confirm = true -- demande confirmation plutôt que d'échouer

-- Diagnostics (LSP + lint) --------------------------------------------------
vim.diagnostic.config({
    virtual_text = true, -- message en fin de ligne
    signs = true, -- marqueurs dans la marge
    underline = true, -- souligne la zone problématique
    update_in_insert = false, -- pas de MAJ pendant la saisie (bruit)
    severity_sort = true, -- erreurs avant warnings, etc.
    float = { border = 'rounded', source = true },
})
