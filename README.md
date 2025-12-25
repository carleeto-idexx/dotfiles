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

To initialize these dotfiles on a new machine:

```bash
sh -c "$(curl -fsLS get.chezmoi.io)" -- init --apply carleeto-idexx
```

### Prerequisites

Ensure you have the following installed:
- [Homebrew](https://brew.sh/)
- `chezmoi`
- `wezterm`
- `aerospace`

## ⌨️ Key Features & Bindings

### WezTerm
- `CMD + d`: Split horizontal
- `CMD + Shift + d`: Split vertical
- `CMD + k`: Clear scrollback

### Zsh
- `z`: Fast directory jumping (zoxide)
- `git`: Enhanced with custom aliases and plugins
