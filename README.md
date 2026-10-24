# configs

Gestionnaire des configurations personnelles versionnées dans
[`luannbertaud/configs`](https://github.com/luannbertaud/configs). Chaque
configuration vit sur sa **branche orpheline** (`nvim`, `xkb`, `aliases`…) et
se déploie dans un dossier de `~/.config`. La branche `main` ne contient que cet
outil : il liste, installe, met à jour et désinstalle ces configurations.

## Lancer

```sh
curl -fsSL https://raw.githubusercontent.com/luannbertaud/configs/main/run.sh | sh
```

`run.sh` détecte la plateforme, télécharge le binaire de la dernière *release*
dans un dossier temporaire et l'exécute. Rien n'est installé de façon
permanente.

Avec des arguments (notez le `-s --` pour `curl | sh`) :

```sh
curl -fsSL .../run.sh | sh -s -- --list
curl -fsSL .../run.sh | sh -s -- --install nvim
```

## Usage

| Commande                    | Effet                                                      |
| --------------------------- | ---------------------------------------------------------- |
| `configs`                   | TUI interactif                                             |
| `configs --list [--json]`   | Liste les configurations et leur état                      |
| `configs --install <nom>`   | Installe (ou met à jour) une configuration                 |
| `configs --install-missing` | Installe toutes les configurations absentes                |
| `configs --update [nom]`    | Met à jour une configuration, ou toutes celles à update    |
| `configs --uninstall <nom>` | Désinstalle une configuration                              |
| `configs --version`         | Affiche la version                                         |

Dans le TUI : `↑/↓` sélectionner, `entrée`/`a` agir sur la sélection
(installer ou mettre à jour), `d` supprimer (avec confirmation), `i` installer
les manquantes, `u` tout mettre à jour, `r` actualiser, `q` quitter. La fenêtre
est centrée : la liste à gauche, le détail de la configuration sélectionnée à
droite.

## États

| État         | Signification                                                        |
| ------------ | -------------------------------------------------------------------- |
| `absent`     | le dossier cible n'existe pas                                        |
| `désync`     | présent mais non géré : pas un dépôt git, mauvais dépôt/branche, ou fichiers suivis modifiés |
| `à update`   | clone géré et propre, mais `HEAD` diffère du tip de la branche        |
| `à jour`     | clone géré et propre, aligné sur le tip de la branche                 |

Une mise à jour qui ne peut pas être un fast-forward passe en **erreur** (affichée
pour la session). Aucune suppression automatique : la désinstallation est
toujours confirmée.

Le déploiement éventuel hors `~/.config` (par ex. une disposition clavier vers
`/usr/share/X11/xkb`) reste manuel : `configs` ne gère que le dossier cible.

## Comment ça marche

1. **Découverte** : un clone partiel éphémère de toutes les branches
   (`git clone --depth 1 --no-single-branch --filter=blob:none --no-checkout`)
   fournit les tips et permet de lire chaque `manifest.json` à la demande. Une
   branche sans manifeste est ignorée. Le dépôt est codé en dur.
2. **État** : `configs` compare `HEAD` local au tip distant et vérifie que le
   dossier est un clone de ce dépôt sur la bonne branche, sans fichier suivi
   modifié (`git status --porcelain --untracked-files=no`).
3. **Installation** : `git clone --single-branch --branch <branche> <url-https>`
   puis exécution de `post_install`, s'il est déclaré.
4. **Mise à jour** : `git pull --ff-only <url-https> <branche>` (l'URL HTTPS
   canonique est toujours utilisée, jamais l'`origin`), puis `post_install` si
   `HEAD` a changé. Un `origin` en SSH reste donc sans effet.
5. **Désinstallation** : `post_uninstall` s'il existe, puis retrait du dossier
   cible.

## Manifeste — `manifest.json`

À la racine de chaque branche-paquet :

```json
{
  "name": "nvim",
  "description": "Configuration Neovim (lazy.nvim, LSP natif, Mason)",
  "target": "~/.config/nvim",
  "prereq": ["nvim", "git", "rg", "fd", "fzf"],
  "post_install": "scripts/post-install.sh",
  "post_uninstall": "scripts/post-uninstall.sh",
  "hidden": false
}
```

| Champ            | Règle                                                                 |
| ---------------- | --------------------------------------------------------------------- |
| `name`           | facultatif ; défaut = nom de branche (identité = branche)             |
| `description`    | libellé affiché                                                       |
| `target`         | obligatoire ; `~/…`, relatif au `$HOME`, ou absolu                    |
| `prereq`         | commandes vérifiées sur le `PATH` ; si absentes, avertir et continuer |
| `post_install`   | chemin relatif optionnel ; exécuté via `sh` (cwd = cible) à l'installation et après un update qui change `HEAD` |
| `post_uninstall` | chemin relatif optionnel ; exécuté via `sh` avant le retrait du dossier |
| `hidden`         | `true` pour masquer du TUI (reste installable par nom)                |

Les scripts reçoivent `CONFIGS_TARGET` et `CONFIGS_COMMIT` ; `CONFIGS_PREV`
est fourni lors d'une mise à jour. Ils doivent être idempotents.

## Ajouter une configuration

1. Créer une branche orpheline portant les fichiers de la config.
2. Y déposer un `manifest.json` (voir ci-dessus).
3. Optionnel : les scripts `post_install` / `post_uninstall`.

## Développement

```sh
go build ./...
go vet ./...
gofmt -l .
```

Une release est publiée en poussant un tag `v*` (workflow
`.github/workflows/release.yml`).
