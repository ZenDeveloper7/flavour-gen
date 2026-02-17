# Flavour Generator CLI

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

- Go 1.21+
- Java JDK (`keytool`)
- Gradle (optional for validation)

## Build

```bash
git clone https://github.com/ZenDeveloper7/flavour-gen.git
cd flavour-gen
go mod tidy
go build -o flavour-gen .
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
- `gradle not installed` → Install Gradle
- `Theme {id} not found` → Check templates folder for `appx_theme{id}_sample` and `appx_theme{id}.gradle`
- `Logo must be PNG`
- `Cannot write to output directory`

## Development

Run dry‑run to inspect actions without writing files:

```bash
./flavour-gen create --client-data example-client.json --theme-id 1 --logo logo.png --dry-run -v
```

## License

MIT (or choose your own). Contributions welcome!
