# Flavor Package

`pkg/flavor/flavor.go` handles theme duplication and Gradle file generation.

## Overview

This package:
1. Copies theme sample folders
2. Replaces placeholders in Gradle files
3. Updates build_type.gradle and flavours.gradle

## Functions

### `SetTemplatesDir(dir string)`

Sets the templates directory. Used by CLI to point to Android project.

### `ListThemes() ([]int, error)`

Scans the project's flavours folder for available theme IDs.

### `DuplicateTheme(themeID int, archiveBasename string, outputDir string, cd *config.ClientData, dryRun bool) (ThemeFiles, error)`

Copies theme folder and creates Gradle file.

**Steps:**
1. Validate source theme folder exists
2. Copy entire theme folder to output
3. Replace placeholders in Gradle file:
   - `applicationId` → package_name
   - `versionName` → version_name
   - `versionCode` → version_code
   - `app_name` → app_name
   - `buildConfigField` values (BASE_URL, etc.)
   - `//Education X` → education_number
4. Update strings.xml with client data
5. Write modified Gradle file

### `AddFlavorToBuildType(archiveBasename, projectRoot string, dryRun bool) error`

Adds signing config entry to `build_type.gradle` inside the `release` block.

### `AddFlavorToFlavours(archiveBasename, projectRoot string, dryRun bool) error`

Adds apply statement to `flavours.gradle`.

## Placeholder Replacement

The package replaces these Gradle values:

| Original | Replaced With |
|----------|---------------|
| `appx_theme{id}` | archivebasename |
| `applicationId "..."` | package_name |
| `versionName "..."` | version_name |
| `versionCode N` | version_code |
| `resValue "string", "app_name", "..."` | app_name |
| `buildConfigField("String", "BASE_URL", "...")` | base_url |
| `buildConfigField("int", "DOT_COUNT", N)` | dot_count |
| `//Education N` | education_number |

## Source Structure

Expected in Android project:
```
app/
├── flavours/
│   └── appx_theme{id}.gradle    # Source gradle
└── src/
    └── appx_theme{id}_sample/   # Source folder
```

## Output Structure

```
output/
└── app/
    ├── flavours/
    │   └── {archivebasename}.gradle
    └── src/
        └── {archivebasename}/
            └── (copied from theme sample)
```
