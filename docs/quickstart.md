# Quick Start

## Basic Usage

Generate a new Android app flavor:

```bash
flavour-gen create \
  --input ./client-data \
  --output-dir ./AndroidStudioProjects/appx-core-android/output
```

## Step-by-Step

### 1. Prepare Input Folder

Create a folder with required files:

```bash
mkdir my-client
cd my-client

# Create data.json
cat > data.json << 'EOF'
{
  "app_name": "My App",
  "archivebasename": "my_app",
  "package_name": "com.mycompany.myapp",
  "version_name": "1.0.0",
  "version_code": 0,
  "theme_id": 1,
  "app_logo": "app_logo.png",
  "base_url": "https://api.myapp.com/",
  "test_base_url": "https://api.myapp.com/"
}
EOF

# Add logo
cp /path/to/logo.png app_logo.png

# Add google-services.json (optional)
cp /path/to/google-services.json .
```

### 2. Run flavour-gen

```bash
flavour-gen create \
  --input ./my-client \
  --output-dir /path/to/project/app/output
```

### 3. Check Output

```bash
ls -la /path/to/project/app/output/
# app/
#   keystore/my_app.jks
#   flavours/my_app.gradle
#   src/my_app/
#     google-services.json
#     res/
#       drawable/app_logo.png
#       drawable-*/ic_notification_icon.png
#       mipmap-*/ic_launcher.png
```

## Dry Run Mode

Preview what would be created without writing files:

```bash
flavour-gen create \
  --input ./my-client \
  --output-dir /path/to/project/app/output \
  --dry-run -v
```

## Options

### Custom Logo

```bash
flavour-gen create \
  --input ./my-client \
  --logo /custom/path/logo.png \
  --output-dir ./output
```

### Custom Background Color

```bash
flavour-gen create \
  --input ./my-client \
  --bg-color "#FF5722" \
  --output-dir ./output
```

### Auto Background (Default)

Uses dominant color from logo corners:

```bash
flavour-gen create \
  --input ./my-client \
  --auto-bg \
  --output-dir ./output
```

## Next Steps

1. Open Android project in Android Studio
2. Build the new flavor: `./gradlew assemble<ArchiveBasename>Debug`
3. Find APK in `app/build/outputs/apk/`

## Troubleshooting

See [Error Handling](./error-handling.md) for common issues.
