#!/bin/bash

# Bootstrap script for macOS and NixOS

echo "🚀 Starting setup..."

# Check for automation preference
FORCE_FLAG=""
if [ -t 0 ]; then
    echo "❓ Do you want a fully automated setup? (Overwrites local conflicts without asking)"
    read -p "Type 'yes' to automate, or press Enter to handle conflicts manually: " choice
    if [[ "$choice" == "yes" ]]; then
        echo "🤖 Automated mode: --force will be used."
        FORCE_FLAG="--force"
    else
        echo "✋ Manual mode: You will be prompted for any file conflicts."
    fi
else
    echo "🚫 Non-interactive environment detected. Proceeding without --force."
fi

if [[ "$OSTYPE" == "darwin"* ]]; then
    echo "🍎 macOS detected"
    if ! command -v brew &> /dev/null; then
        echo "Installing Homebrew..."
        /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
        # Source Homebrew for the current session
        eval "$(/opt/homebrew/bin/brew shellenv)" 2>/dev/null || eval "$(/usr/local/bin/brew shellenv)" 2>/dev/null
    fi
    # Packages to install
    packages=(chezmoi zoxide asdf fzf tig)
    
    # Special handle for oh-my-posh (official tap)
    if ! brew list oh-my-posh &>/dev/null; then
        brew install --formula jandedobbeleer/oh-my-posh/oh-my-posh
    fi

    for pkg in "${packages[@]}"; do
        if ! brew list "$pkg" &>/dev/null; then
            brew install "$pkg"
        fi
    done
elif grep -q "ID=nixos" /etc/os-release 2>/dev/null; then
    echo "❄️ NixOS detected"
    nix-env -iA nixos.chezmoi nixos.zoxide nixos.oh-my-posh nixos.asdf-vm nixos.fzf nixos.tig
else
    echo "⚠️ Unsupported OS. This script only supports macOS and NixOS."
    exit 1
fi

# Install full Brewfile if present (macOS only). This covers all formulae,
# casks, taps, and VS Code extensions captured at dump time. Idempotent.
if [[ "$OSTYPE" == "darwin"* ]] && [ -f "$PWD/Brewfile" ]; then
    echo "📦 Installing packages from Brewfile..."
    brew bundle --file="$PWD/Brewfile"
fi

# Clone zsh-autosuggestions if missing
if [ ! -d "$HOME/.zsh/zsh-autosuggestions" ]; then
    echo "📥 Cloning zsh-autosuggestions..."
    mkdir -p "$HOME/.zsh"
    git clone https://github.com/zsh-users/zsh-autosuggestions "$HOME/.zsh/zsh-autosuggestions"
fi

# Apply dotfiles
echo "✨ Applying dotfiles..."
if command -v chezmoi &> /dev/null; then
    chezmoi apply --source "$PWD" $FORCE_FLAG
else
    # Fallback if chezmoi was just installed and PATH isn't updated in this script's subshell
    PATH="$PATH:/opt/homebrew/bin:/usr/local/bin" chezmoi apply --source "$PWD" $FORCE_FLAG
fi

echo "✅ Setup complete!"

echo ""
echo "✨ Next Steps:"
echo "1. Run 'exec zsh' to start using your new environment immediately."
echo "2. Enjoy your new terminal setup!"
echo ""
