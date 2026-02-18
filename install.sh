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

URL="https://github.com/ZenDeveloper7/flavour-gen/releases/latest/download/$BINARY"

echo "Downloading $URL..."
curl -sSL "$URL" -o "/tmp/$BINARY" || { echo "Download failed"; exit 1; }

chmod +x "/tmp/$BINARY"

DEST="/usr/local/bin/flavour-gen"
if [ -w "$(dirname "$DEST")" ]; then
  mv "/tmp/$BINARY" "$DEST"
else
  sudo mv "/tmp/$BINARY" "$DEST"
fi

echo "Installed to $DEST"
