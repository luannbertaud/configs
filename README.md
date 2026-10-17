# configs

Gestionnaire des configurations personnelles versionnées dans
[`luannbertaud/configs`](https://github.com/luannbertaud/configs). Chaque
configuration vit sur sa **branche orpheline** (`nvim`, `xkb`, `aliases`…) et
se déploie dans un dossier de `~/.config`. La branche `main` ne contient que cet
outil : il liste, installe et met à jour ces configurations.

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

| Commande                    | Effet                                                        |
| --------------------------- | ------------------------------------------------------------ |
| `configs`                   | TUI interactif                                               |
| `configs --list [--json]`   | Liste les configurations et leur état                        |
| `configs --install <nom>`   | Installe (ou met à jour) une configuration                   |
| `configs --install-missing` | Installe toutes les configurations manquantes                |
| `configs --update [nom]`    | Met à jour une configuration, ou toutes si le nom est omis   |
| `configs --version`         | Affiche la version                                           |

Dans le TUI : `↑/↓` naviguer, `entrée` détails, `i` installer les manquants,
`u` tout mettre à jour, `r` actualiser, `q` quitter.

**Hors périmètre :** la suppression (il n'y a jamais d'action destructive).
Le déploiement éventuel hors `~/.config` (par ex. une disposition clavier vers
`/usr/share/X11/xkb`) reste manuel.

## Comment ça marche

1. **Découverte** : `git ls-remote --heads` liste les branches du dépôt, puis
   `configs` récupère le `configs.json` de chacune. Une branche sans
   `configs.json` est ignorée (ce n'est pas une configuration).
2. **État** : une configuration est « installée » si son dossier cible existe.
   S'il existe mais n'est pas un clone git de la branche, elle est signalée en
   **désync** (jamais supprimée automatiquement).
3. **Installation / mise à jour** :
   - si la branche fournit un `install.sh`, `configs` le délègue
     (`curl … | sh`), le script étant idempotent et auto-localisé ;
   - sinon, `configs` utilise git : `git clone --branch <branche>
     --single-branch` à l'installation, `git -C <cible> pull --ff-only` à la
     mise à jour.

## Ajouter une configuration

1. Créer une branche orpheline portant les fichiers de la config.
2. Y déposer un `configs.json` à la racine :

   ```json
   {
     "desc": "Configuration Neovim",
     "target": "~/.config/nvim",
     "prereq": ["neovim", "git", "rg", "fd", "fzf"]
   }
   ```

   Champs : `desc` (libellé), `target` (dossier cible relatif au `$HOME`,
   `~` accepté, ou chemin absolu), `prereq` (commandes vérifiées via `PATH`),
   `hidden` (`true` pour ne pas l'afficher dans le TUI), `name` (facultatif,
   sinon le nom de la branche).
3. Optionnel : un `install.sh` idempotent à la racine si l'installation
   demande plus qu'un simple clone.

## Développement

```sh
go build ./...
go vet ./...
gofmt -l .
```

Une release est publiée en poussant un tag `v*` (workflow
`.github/workflows/release.yml`).
