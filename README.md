# aliases

Shell aliases (bash/zsh), kept in `${XDG_CONFIG_HOME:-~/.config}/aliases/aliases.sh` and sourced from `~/.bashrc` / `~/.zshrc`.

Install (or update):

```sh
curl -fsSL https://raw.githubusercontent.com/luannbertaud/configs/aliases/install.sh | sh
```

Current shell only, nothing written:

```sh
eval "$(curl -fsSL https://raw.githubusercontent.com/luannbertaud/configs/aliases/aliases.sh)"
```

Uninstall:

```sh
curl -fsSL https://raw.githubusercontent.com/luannbertaud/configs/aliases/install.sh | sh -s -- --uninstall
```
