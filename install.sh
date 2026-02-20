#!/bin/bash
set -e
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"
case "$ARCH" in
  x86_64) ARCH="amd64" ;;
  arm64) ARCH="arm64" ;;
  *) echo "Unsupported arch: $ARCH"; exit 1 ;;
esac
BINARY="flavour-gen-$OS-$ARCH"
[ "$OS" = "windows" ] && BINARY="$BINARY.exe"
URL="https://raw.githubusercontent.com/$GITHUB_REPOSITORY/gh-pages/$BINARY"
echo "Downloading $URL..."
curl -sSL "$URL" -o "/tmp/$BINARY" || { echo "Download failed"; exit 1; }
chmod +x "/tmp/$BINARY"
DEST="/usr/local/bin/flavour-gen"
if [ -w "$(dirname "$DEST")" ]; then
  sudo mv "/tmp/$BINARY" "$DEST"
else
  mv "/tmp/$BINARY" "$DEST"
fi
echo "Installed to $DEST"
