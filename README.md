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

### New Machine Setup
To initialize these dotfiles on a new machine:

```bash
# Using HTTPS (Ensure you have a Personal Access Token ready)
sh -c "$(curl -fsLS get.chezmoi.io)" -- init --apply carleeto-idexx

# OR using SSH (If your SSH keys are already configured)
sh -c "$(curl -fsLS get.chezmoi.io)" -- init --apply git@github.com:carleeto-idexx/dotfiles.git
```

> [!IMPORTANT]
> Since this is a **private repository**, you must ensure your GitHub credentials (PAT or SSH keys) are configured on the machine before running the command.

> [!TIP]
> This command will install `chezmoi`, clone the entire repository (including the bootstrap script), and automatically run `run_once_before_bootstrap.sh` to install CLI dependencies.

### Existing Machine (Sync Changes)
If you already have `chezmoi` installed and want to pull the latest changes:

```bash
chezmoi update
```
*This command pulls the latest commits and runs `chezmoi apply`.*

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

## 🧪 Testing & Simulation

To simulate setting up these dotfiles on a completely new machine without formatting your computer:

```bash
# 1. Remove the local source repository
rm -rf "$(chezmoi source-path)"

# 2. Clear the chezmoi state (forces run_once_ scripts to re-run)
chezmoi state delete --bucket=scriptState

# 3. Optional: Delete the chezmoi binary if installed in ~/bin
rm -f ~/bin/chezmoi

# 4. Run the Quick Start command from the top of this README
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
