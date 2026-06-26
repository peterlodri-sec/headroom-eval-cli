#!/bin/bash
# One-line installer for headroom-eval CLI
set -e

echo "🐋 headroom-eval-cli installer"
echo ""

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
case "$ARCH" in x86_64) ARCH="amd64";; aarch64|arm64) ARCH="arm64";; esac

VERSION="${1:-latest}"
if [ "$VERSION" = "latest" ]; then
    URL="https://github.com/peterlodri-sec/headroom-eval-cli/releases/latest/download/headroom-eval-${OS}-${ARCH}"
else
    URL="https://github.com/peterlodri-sec/headroom-eval-cli/releases/download/${VERSION}/headroom-eval-${OS}-${ARCH}"
fi

DEST="${2:-/usr/local/bin/headroom-eval}"

echo "  OS:    $OS"
echo "  Arch:  $ARCH"
echo "  From:  $URL"
echo "  To:    $DEST"
echo ""

if command -v curl >/dev/null 2>&1; then
    curl -fsSL "$URL" -o "$DEST"
elif command -v wget >/dev/null 2>&1; then
    wget -q "$URL" -O "$DEST"
else
    echo "Need curl or wget."
    exit 1
fi

chmod +x "$DEST"
echo "✅ headroom-eval installed to $DEST"
echo "   Run: headroom-eval"
