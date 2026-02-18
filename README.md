# Flavour Generator CLI

[![Release](https://img.shields.io/github/v/release/ZenDeveloper7/flavour-gen?include_prereleases&label=release)](https://github.com/ZenDeveloper7/flavour-gen/releases/latest)
[![Go](https://img.shields.io/badge/Go-1.21-blue.svg)](https://golang.org)
[![License](https://img.shields.io/github/license/ZenDeveloper7/flavour-gen.svg)](https://github.com/ZenDeveloper7/flavour-gen/blob/master/LICENSE)
[![Platform](https://img.shields.io/badge/platform-linux%20%7C%20macOS%20%7C%20Windows-lightgrey.svg)](https://github.com/ZenDeveloper7/flavour-gen)

A Go-based CLI tool to generate complete Android app flavor packages, including icons, keystores, and Gradle configuration.

## Features

- Create full flavor package from a JSON client data file and a PNG logo
- Auto‑detect background color from logo corners (optional)
- Generate adaptive launcher icons and legacy PNGs in all densities
- Create notification icons (WEBP)
- Duplicate theme templates and substitute Gradle placeholders
- Generate Java keystore (JKS) for signing
- Dry‑run mode to preview output

## Prerequisites

- **For prebuilt binary:** `curl`, `bash`, `sudo` (optional) – no Go needed.
- **For building from source:** Go 1.21+, Java JDK (`keytool`), Gradle (optional for validation).

## Installation

### Quick Install (All Platforms)

```bash
curl -sSL https://raw.githubusercontent.com/ZenDeveloper7/flavour-gen-releases/master/install.sh | bash
```

This downloads and installs the latest binary to `/usr/local/bin/flavour-gen`.

### Manual Download

Download the latest release for your platform:

- **Linux:** [flavour-gen-linux-amd64](https://github.com/ZenDeveloper7/flavour-gen-releases/raw/master/flavour-gen-linux-amd64)
- **macOS:** [flavour-gen-darwin-amd64](https://github.com/ZenDeveloper7/flavour-gen-releases/raw/master/flavour-gen-darwin-amd64)
- **Windows:** [flavour-gen-windows-amd64.exe](https://github.com/ZenDeveloper7/flavour-gen-releases/raw/master/flavour-gen-windows-amd64.exe)

```bash
# Linux/macOS
chmod +x flavour-gen-linux-amd64
sudo mv flavour-gen-linux-amd64 /usr/local/bin/flavour-gen

# Windows
move flavour-gen-windows-amd64.exe flavour-gen.exe
```

### Build from Source

```bash
git clone https://github.com/ZenDeveloper7/flavour-gen.git
cd flavour-gen
go mod tidy
go build -o flavour-gen .
sudo mv flavour-gen /usr/local/bin/
```

### Go Install

If you have Go installed:

```bash
go install github.com/ZenDeveloper7/flavour-gen@latest
```

## Usage

```bash
./flavour-gen create --help
./flavour-gen create \
  --client-data example-client.json \
  --theme-id 1 \
  --logo path/to/logo.png \
  --output-dir ./output
```

### Client Data File

See `templates/client-data-template.json` for required fields. Example:

```json
{
  "app_name": "App Name",
  "archivebasename": "app_name",
  "package_name": "com.package.name",
  "version_name": "1.0.0",
  "version_code": 0,
  "base_url": "https://...",
  "test_base_url": "https://...",
  "firebase_url": "https://...firebaseio.com/",
  "dynamic_link_domain": "https://...",
  "dynamic_link_prefix": "https://..."
}
```

Computed fields (generated automatically): `IDENTITY`, `DOT_COUNT`, `ALT_APP_NAME`, `DOWNLOAD_FOLDER_NAME`.

## Output Structure

```
output/
└── app/
    ├── keystore/<archivebasename>.jks
    ├── flavours/<archivebasename>.gradle
    └── src/<archivebasename>/
        ├── google-services.json (if in theme template)
        └── res/
            ├── drawable(-hdpi|mdpi|xhdpi|xxhdpi|xxxhdpi)/
            │   └── ic_notification_icon.webp
            ├── mipmap-anydpi-v26/
            │   └── ic_launcher.xml
            └── mipmap-(hdpi|mdpi|xhdpi|xxhdpi|xxxhdpi)/
                └── ic_launcher.png
```

## Error Handling

- `keytool not installed` → Install Java JDK
- `gradle not installed` → Install Gradle (optional, only for validation)
- `Theme {id} not found` → Check templates folder for `appx_theme{id}_sample` and `appx_theme{id}.gradle`
- `Logo must be PNG`
- `Cannot write to output directory`

## Development

Run dry‑run to inspect actions without writing files:

```bash
./flavour-gen create --client-data example-client.json --theme-id 1 --logo logo.png --dry-run -v
```

## License

MIT. Contributions welcome!
