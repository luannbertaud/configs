# Aliases

Alias shell (bash/zsh), installés dans `${XDG_CONFIG_HOME:-~/.config}/aliases`
et sourcés depuis `~/.bashrc` / `~/.zshrc`.

## Installation autonome (curl)

```sh
curl -fsSL https://raw.githubusercontent.com/luannbertaud/configs/aliases/install.sh | sh
```

ou, depuis un shell où les alias sont déjà chargés :

```sh
aliases-update
```

## Via le gestionnaire `configs`

```sh
configs --install aliases
```

Le gestionnaire clone la branche, puis `scripts/post-install.sh` réutilise
`install.sh` en mode local : `aliases.sh` est déjà dans le clone, il n'est pas
retéléchargé.

- Mise à jour : `configs --update aliases`
- Désinstallation : `configs --uninstall aliases`

## Shell courant seulement, rien écrit

```sh
eval "$(curl -fsSL https://raw.githubusercontent.com/luannbertaud/configs/aliases/aliases.sh)"
```

## Désinstallation autonome

```sh
curl -fsSL https://raw.githubusercontent.com/luannbertaud/configs/aliases/install.sh | sh -s -- --uninstall
```

ou `aliases-uninstall`.

Dans un clone, `install.sh --uninstall` ne retire que la ligne de source ; c'est
le gestionnaire qui supprime ensuite le dossier du paquet.
