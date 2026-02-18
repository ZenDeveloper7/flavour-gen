# Installation

## Prerequisites

- **Runtime:** `curl`, `bash`
- **Optional:** Java JDK (for keystore generation)

## Installation Methods

### 1. Quick Install (Recommended)

```bash
curl -sSL https://raw.githubusercontent.com/ZenDeveloper7/flavour-gen/master/install.sh | bash
```

This downloads the latest release and installs to `/usr/local/bin/flavour-gen`.

### 2. Homebrew (macOS/Linux)

```bash
brew tap ZenDeveloper7/tap
brew install ZenDeveloper7/tap/flavour-gen
```

### 3. Go Install

If you have Go installed:

```bash
go install github.com/ZenDeveloper7/flavour-gen@latest
```

### 4. Manual Download

Download from [Releases](https://github.com/ZenDeveloper7/flavour-gen/releases):

```bash
# Linux
curl -LO https://github.com/ZenDeveloper7/flavour-gen/releases/latest/download/flavour-gen-linux-amd64
chmod +x flavour-gen-linux-amd64
sudo mv flavour-gen-linux-amd64 /usr/local/bin/flavour-gen

# macOS
curl -LO https://github.com/ZenDeveloper7/flavour-gen/releases/latest/download/flavour-gen-darwin-amd64
chmod +x flavour-gen-darwin-amd64
sudo mv flavour-gen-darwin-amd64 /usr/local/bin/flavour-gen

# Windows
# Download from releases page
```

## Verification

```bash
flavour-gen --version
flavour-gen --help
```

## Updating

```bash
# Quick install
curl -sSL https://raw.githubusercontent.com/ZenDeveloper7/flavour-gen/master/install.sh | bash

# Homebrew
brew upgrade flavour-gen

# Go
go install github.com/ZenDeveloper7/flavour-gen@latest
```

## Uninstalling

```bash
# If installed to /usr/local/bin
sudo rm /usr/local/bin/flavour-gen

# Homebrew
brew uninstall flavour-gen

# Go
rm $(go env GOPATH)/bin/flavour-gen
```
