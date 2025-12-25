# 🌌 Dotfiles

My personal macOS development environment, managed with [chezmoi](https://www.chezmoi.io/).

## 🛠 Tech Stack

| Component | Tool | Description |
| :--- | :--- | :--- |
| **Window Manager** | [AeroSpace](https://nikitabobko.github.io/AeroSpace/) | i3-like tiling window manager for macOS |
| **Terminal** | [WezTerm](https://wezfurlong.org/wezterm/) | GPU-accelerated terminal (Dracula theme) |
| **Shell** | [Zsh](https://www.zsh.org/) | With `oh-my-posh`, `autosuggestions`, and `zoxide` |
| **Prompt** | [Oh My Posh](https://ohmyposh.dev/) | Using the `zen.toml` theme |
| **Version Control** | [Git](https://git-scm.com/) & [Tig](https://jonas.github.io/tig/) | Text-mode interface for Git |
| **Tool Manager** | [asdf](https://asdf-vm.com/) | Extendable version manager |
| **Security** | Bitwarden | Integrated SSH agent |

## 🚀 Quick Start

This section guides you through setting up these dotfiles on a new machine or updating an existing one.

### New Machine Setup
To initialize these dotfiles on a new machine, run the following command. This will install `chezmoi`, clone this repository, and execute the `run_once_bootstrap.sh` script to install all necessary dependencies.

```bash
sh -c "$(curl -fsLS get.chezmoi.io)" -- init --apply carleeto-idexx
```

> [!TIP]
> This command will install `chezmoi`, clone this repository, and run the `run_once_bootstrap.sh` script to install dependencies.

### Existing Machine (Sync Changes)
If you already have `chezmoi` installed and want to pull the latest changes (including new scripts like the bootstrap):

```bash
chezmoi update
```
*This command will perform a `git pull --rebase` in your source directory and `apply` any changes.*

> [!TIP]
> This command automatically installs all required dependencies (Homebrew, zoxide, etc.) on both macOS and NixOS.

## 🛠 Maintenance & Workflow

### 1. Update your Configuration
To modify a file (e.g., your `.zshrc`), use the `chezmoi edit` command (aliased to `ce`):

```bash
ce ~/.zshrc
```

### 2. Push Changes to GitHub
When you're happy with your local changes, sync them back to your repository and push to GitHub using the `csync` alias:

```bash
csync
```

### 3. Pull Changes from GitHub
If you've made changes on another machine, pull them down to your current machine:

```bash
chezmoi update
```

## ⌨️ Key Features & Bindings

- **🚀 Automated Setup**: One-command installation for macOS and NixOS.
- **🛡️ Identity Agnostic**: No hardcoded Git names, emails, or SSH keys. Configure each machine locally for maximum flexibility.
- **📦 Package Management**: Integrated support for Homebrew, asdf, and Nix.
- **🐚 Modern Shell**: Zsh with `oh-my-posh` (Zen theme), `zoxide`, and auto-suggestions.

### WezTerm
- `CMD + d`: Split horizontal
- `CMD + Shift + d`: Split vertical
- `CMD + k`: Clear scrollback

### Zsh
- `z`: Fast directory jumping (zoxide)
- `git`: Enhanced with custom aliases and plugins
