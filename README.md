# Flavour Generator CLI

[![Release](https://img.shields.io/github/v/release/ZenDeveloper7/flavour-gen-releases?include_prereleases&label=release)](https://github.com/ZenDeveloper7/flavour-gen-releases/releases)
[![Go](https://img.shields.io/badge/Go-1.21-blue.svg)](https://golang.org)
[![Platform](https://img.shields.io/badge/platform-linux%20%7C%20macOS%20%7C%20Windows-lightgrey.svg)](https://github.com/ZenDeveloper7/flavour-gen-releases)

A Go-based CLI tool to generate complete Android app flavor packages from client data. Automatically creates icons, keystores, Gradle configuration, and updates build files.

## Features

- 🎨 **Icon Generation** - Adaptive launcher icons + legacy PNGs, notification icons (white on transparent)
- 🔐 **Keystore Generation** - Creates JKS keystore and updates `keystore.gradle`
- 📄 **Gradle Updates** - Duplicates theme, updates `build_type.gradle` and `flavours.gradle`
- 🔍 **Smart Detection** - Auto-detects background color from logo corners
- 🧪 **Dry-Run Mode** - Preview actions without writing files
- 📁 **Simple Input** - Just drop a folder with `data.json` and `google-services.json`

## Prerequisites

- **For prebuilt binary:** `curl`, `bash`, `sudo` (optional)
- **For building from source:** Go 1.21+, Java JDK (`keytool`)

## Installation

### Quick Install (All Platforms)

```bash
curl -sSL https://raw.githubusercontent.com/ZenDeveloper7/flavour-gen-releases/master/install.sh | bash
```

### Manual Download

Download from [Releases](https://github.com/ZenDeveloper7/flavour-gen-releases/releases):

- **Linux:** `flavour-gen-linux-amd64`
- **macOS:** `flavour-gen-darwin-amd64`
- **Windows:** `flavour-gen-windows-amd64.exe`

### Build from Source

```bash
git clone https://github.com/ZenDeveloper7/flavour-gen.git
cd flavour-gen
go build -o flavour-gen .
```

## Quick Start

```bash
flavour-gen create --input ./client-folder --output-dir /path/to/AndroidProject/app/output
```

## Input Folder Structure

```
client-folder/
├── data.json              # Client configuration (required)
├── google-services.json  # Firebase config (optional)
└── app_logo.png          # Logo image (required)
```

## data.json Fields

| Field | Required | Description |
|-------|----------|-------------|
| `app_name` | ✅ | Display name of the app |
| `archivebasename` | ✅ | Archive name (e.g., `physics_setu`) |
| `package_name` | ✅ | Android package (e.g., `com.ydcfzb.zgizxw`) |
| `version_name` | ✅ | Version string (e.g., `1.0.0`) |
| `version_code` | ✅ | Version code (integer) |
| `theme_id` | ✅ | Theme ID to use (e.g., `1`) |
| `app_logo` | ✅ | Path to logo file relative to input folder |
| `education_number` | - | Education number for gradle comment |
| `base_url` | ✅ | Production API base URL |
| `test_base_url` | ✅ | Test API base URL |
| `firebase_url` | - | Firebase database URL |
| `dynamic_link_domain` | - | Dynamic Links domain |
| `dynamic_link_prefix` | - | Dynamic Links prefix |

**Example:**

```json
{
  "app_name": "Physics Setu",
  "archivebasename": "physics_setu",
  "package_name": "com.ydcfzb.zgizxw",
  "version_name": "1.0.0",
  "version_code": 0,
  "theme_id": 1,
  "app_logo": "app_logo.png",
  "education_number": 20,
  "base_url": "https://physicssetuapi.akamai.net.in/",
  "test_base_url": "https://physicssetuapi.akamai.net.in/",
  "firebase_url": "https://physicssetuappx.firebaseio.com/",
  "dynamic_link_domain": "https://physicssetu.classx.co.in/",
  "dynamic_link_prefix": "https://physicssetu.page.link"
}
```

## Usage

```bash
# Basic usage
flavour-gen create --input ./client-data --output-dir ./app/output

# With custom logo (overrides app_logo in data.json)
flavour-gen create --input ./client-data --logo ./custom-logo.png --output-dir ./app/output

# With custom background color
flavour-gen create --input ./client-data --bg-color "#FF5722" --output-dir ./app/output

# Dry run (preview without creating files)
flavour-gen create --input ./client-data --output-dir ./app/output --dry-run -v
```

## Output Structure

```
output/
└── app/
    ├── keystore/
    │   └── <archivebasename>.jks          # Generated keystore
    ├── flavours/
    │   └── <archivebasename>.gradle       # Flavor gradle
    └── src/
        └── <archivebasename>/
            ├── google-services.json         # Copied from input
            └── res/
                ├── drawable/
                │   └── app_logo.png        # 512x512 logo
                ├── drawable-*/
                │   └── ic_notification_icon.png
                ├── mipmap-*/
                │   └── ic_launcher.png
                └── mipmap-anydpi-v26/
                    └── ic_launcher.xml
```

## What Gets Updated

The CLI automatically updates these files in your Android project:

1. **`app/flavours/<name>.gradle`** - Created from theme template with replaced values
2. **`app/build_type.gradle`** - Added `productFlavors.<name>.signingConfig` entry
3. **`app/flavours.gradle`** - Added `apply from: './flavours/<name>.gradle'`
4. **`app/keystore.gradle`** - Added signing config block

## Error Handling

| Error | Solution |
|-------|----------|
| `keytool not found` | Install Java JDK |
| `theme_id is required` | Add `theme_id` to data.json |
| `app_logo is required` | Add `app_logo` to data.json |
| `Theme X not found` | Ensure theme X exists in `app/flavours/appx_themeX.gradle` |
| `Logo must be PNG` | Convert logo to PNG format |

## Development

```bash
# Run tests
go test ./...

# Build
go build -o flavour-gen .

# Dry run with verbose output
./flavour-gen create --input ./test-data --output-dir ./output --dry-run -v
```

## License

MIT
