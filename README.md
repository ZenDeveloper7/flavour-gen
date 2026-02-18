# Flavour Generator CLI

[![Release](https://img.shields.io/github/v/release/ZenDeveloper7/flavour-gen?include_prereleases&label=release)](https://github.com/ZenDeveloper7/flavour-gen/releases)
[![Go](https://img.shields.io/badge/Go-1.21-blue.svg)](https://golang.org)
[![Platform](https://img.shields.io/badge/platform-linux%20%7C%20macOS%20%7C%20Windows-lightgrey.svg)](https://github.com/ZenDeveloper7/flavour-gen)

A Go-based CLI tool to generate complete Android app flavor packages from client data. Supports single or multiple clients at once. Automatically creates icons, keystores, Gradle configuration, and updates build files.

## Features

- 🎨 **Icon Generation** - Adaptive launcher icons, legacy PNGs, notification icons
- 🔐 **Keystore Generation** - Creates JKS keystore and updates `keystore.gradle`
- 📄 **Gradle Updates** - Duplicates theme, updates build files
- 👥 **Multi-Client Support** - Create multiple flavors from a single JSON array
- 🧪 **Dry-Run Mode** - Preview actions without writing files
- 📁 **Simple Input** - Just drop a folder with `data.json`

## Prerequisites

- **Runtime:** `curl`, `bash`
- **Optional:** Java JDK (for keystore generation)

## Installation

### Quick Install

```bash
curl -sSL https://raw.githubusercontent.com/ZenDeveloper7/flavour-gen/master/install.sh | bash
```

### Other Methods

- **Homebrew:** `brew tap ZenDeveloper7/tap && brew install flavour-gen`
- **Go Install:** `go install github.com/ZenDeveloper7/flavour-gen@latest`
- **Manual:** Download from [Releases](https://github.com/ZenDeveloper7/flavour-gen/releases)

## Quick Start

### Single Client

```bash
flavour-gen create \
  --input ./client-folder \
  --output-dir /path/to/AndroidProject/app/output
```

### Multiple Clients

Create multiple flavors at once:

```bash
flavour-gen create \
  --input ./clients-folder \
  --output-dir /path/to/AndroidProject/app/output
```

## Input Format

### Single Client (data.json)

```json
{
  "app_name": "My App",
  "archivebasename": "my_app",
  "package_name": "com.mycompany.myapp",
  "theme_id": 1,
  "app_logo": "app_logo.png",
  "base_url": "https://api.example.com/"
}
```

### Multiple Clients (data.json)

```json
[
  {
    "app_name": "Client App 1",
    "archivebasename": "client_app_1",
    "package_name": "com.client1.app",
    "theme_id": 1,
    "app_logo": "logo1.png",
    "base_url": "https://api.client1.com/"
  },
  {
    "app_name": "Client App 2",
    "archivebasename": "client_app_2",
    "package_name": "com.client2.app",
    "theme_id": 2,
    "app_logo": "logo2.png",
    "base_url": "https://api.client2.com/"
  }
]
```

### Input Folder Structure

```
input-folder/
├── data.json              # Single client or array (required)
├── google-services.json   # Required (provides education_number)
├── logo1.png            # Client 1 logo
└── logo2.png            # Client 2 logo
```

## Usage

```bash
# Basic usage
flavour-gen create --input ./client-data --output-dir ./app/output

# Dry run
flavour-gen create --input ./client-data --output-dir ./app/output --dry-run -v

# Custom background color
flavour-gen create --input ./client-data --bg-color "#FF5722" --output-dir ./app/output
```

## Output Structure

```
output/
└── app/
    ├── keystore/
    │   ├── client1.jks
    │   └── client2.jks
    ├── flavours/
    │   ├── client1.gradle
    │   └── client2.gradle
    └── src/
        ├── client1/
        │   ├── google-services.json
        │   ├── ic_launcher-playstore.png
        │   └── res/...
        └── client2/
            ├── google-services.json
            ├── ic_launcher-playstore.png
            └── res/...
```

## Error Handling

| Error | Solution |
|-------|----------|
| `keytool not found` | Install Java JDK |
| `theme_id is required` | Add `theme_id` to data.json |
| `app_logo is required` | Add `app_logo` to data.json |
| `Theme X not found` | Ensure theme exists in project's flavours folder |

## Documentation

- [Installation](./docs/installation.md)
- [Quick Start](./docs/quickstart.md)
- [Input Format](./docs/input-format.md)
- [Error Handling](./docs/error-handling.md)
- [Package Docs](./docs/pkg/)

## License

MIT
