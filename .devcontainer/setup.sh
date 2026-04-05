#!/usr/bin/env bash
set -euo pipefail

echo "==> Configuring Git..."
git config --global user.name "0x1f6"
git config --global user.email "178943044+0x1f6@users.noreply.github.com"

echo "==> Installing CodeCrafters CLI..."
curl -fsSL https://codecrafters.io/install.sh | bash

echo "==> Verifying CodeCrafters CLI..."
codecrafters ping

echo "==> Setup complete!"
