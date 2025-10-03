#!/bin/bash

set -e  # Exit on error

INSTALL_DIR="$HOME/.bash-butler/bin"

echo "Creating Directory: $INSTALL_DIR"
mkdir -p "$INSTALL_DIR"

echo "Building Binary..."
go build -o "$INSTALL_DIR/bash-butler" ./cmd/bash-butler/

echo "Making Binary Executable..."
chmod +x "$INSTALL_DIR/bash-butler"

echo "Binary moved to: $INSTALL_DIR"

# Add to PATH if not already added
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    echo 'export PATH="$HOME/.bash-butler/bin:$PATH"' >> ~/.bashrc
    echo 'export PATH="$HOME/.bash-butler/bin:$PATH"' >> ~/.zshrc
    echo "Added $INSTALL_DIR to PATH. Restart your terminal or run:"
    echo "source ~/.bashrc  (or source ~/.zshrc for zsh users)"
else
    echo "$INSTALL_DIR is already in PATH."
fi

