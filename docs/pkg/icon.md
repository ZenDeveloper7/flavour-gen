# Icon Package

`pkg/icon/icon.go` handles icon generation for Android app flavors.

## Overview

Generates:
- **Launcher Icons** - Adaptive icons + legacy PNGs
- **Notification Icons** - White on transparent background
- **App Logo** - 512x512 PNG in drawable folder

## Functions

### `GetBackgroundColor(bgColor string, autoBG bool, logoPath string) (color.RGBA, error)`

Extracts background color for adaptive icon.

**Options:**
- If `bgColor` provided (#RRGGBB) → use it
- If `autoBG` is true → sample 4 corners of logo, average them
- Otherwise → white (#FFFFFF)

### `GenerateAppIcon(logoPath, baseName, outputDir string, bg color.RGBA, dryRun bool) (IconPaths, error)`

Creates launcher icons.

**Outputs:**
- `res/drawable/app_logo.png` - 512x512 logo
- `res/mipmap-anydpi-v26/ic_launcher.xml` - Adaptive icon XML
- `res/mipmap-{density}/ic_launcher.png` - Legacy icons

**Densities:**
| Density | Size |
|---------|------|
| mdpi | 48px |
| hdpi | 72px |
| xhdpi | 96px |
| xxhdpi | 144px |
| xxxhdpi | 192px |

### `GenerateNotificationIcon(logoPath, baseName, outputDir string, dryRun bool) (string, error)`

Creates notification icons.

**Process:**
1. Open logo image
2. Convert to grayscale
3. Create white logo on transparent background
4. Save as PNG at various densities

**Outputs:**
- `res/drawable-{density}/ic_notification_icon.png`

**Densities:**
| Density | Size |
|---------|------|
| mdpi | 24px |
| hdpi | 36px |
| xhdpi | 48px |
| xxhdpi | 72px |
| xxxhdpi | 96px |

## IconPaths Struct

```go
type IconPaths struct {
    AdaptiveXML   string
    LegacyMDpi    string
    LegacyHdpi    string
    LegacyXhdpi   string
    LegacyXXhdpi  string
    LegacyXXXhdpi string
    AppLogo       string
}
```

## Dependencies

- `github.com/disintegration/imaging` - Image processing
