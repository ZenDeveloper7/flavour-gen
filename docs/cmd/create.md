# Create Command

`cmd/create.go` implements the `create` subcommand that generates Android app flavors.

## Overview

This is the main command that:
1. Validates prerequisites (keytool)
2. Loads client data from input folder
3. Generates icons (launcher + notification)
4. Duplicates theme folder
5. Creates keystore
6. Updates Gradle files

## Usage

```bash
flavour-gen create --input ./client-folder --output-dir ./output
```

## Flags

| Flag | Required | Description |
|------|----------|-------------|
| `--input` | Yes | Input folder with data.json |
| `--output-dir` | No | Output directory (default: ./output) |
| `--logo` | No | Logo PNG file (uses app_logo from data.json if not provided) |
| `--bg-color` | No | Background color #RRGGBB |
| `--auto-bg` | No | Auto-detect background from logo (default: true) |
| `--dry-run` | No | Preview without creating files |
| `--verbose`, `-v` | No | Verbose logging |

## Input Folder Structure

```
client-folder/
├── data.json              # Required
├── google-services.json   # Optional
└── app_logo.png         # Required (or via --logo)
```

## Workflow

1. **Prerequisites Check** - Verify keytool is installed
2. **Validate Input** - Check data.json and required fields
3. **Detect Android Project** - Find project root from output-dir
4. **Validate Theme** - Ensure theme ID exists in project
5. **Create Output Structure** - Create necessary directories
6. **Generate Icons** - Create launcher and notification icons
7. **Duplicate Theme** - Copy theme folder and update Gradle
8. **Update Gradle Files** - Modify build_type.gradle, flavours.gradle
9. **Generate Keystore** - Create JKS file
10. **Copy Config Files** - Copy google-services.json

## Related Files

- `pkg/flavor/` - Theme duplication logic
- `pkg/icon/` - Icon generation
- `pkg/keystore/` - Keystore generation
- `pkg/config/` - Configuration parsing
