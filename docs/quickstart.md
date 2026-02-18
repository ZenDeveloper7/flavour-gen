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
mkdir my-clients
cd my-clients

# Create data.json (single client)
cat > data.json << 'EOF'
{
  "app_name": "My App",
  "archivebasename": "my_app",
  "package_name": "com.mycompany "theme_id":.myapp",
  1,
  "app_logo": "app_logo.png",
  "base_url": "https://api.myapp.com/"
}
EOF

# Add logo
cp /path/to/logo.png app_logo.png
```

### 2. Run flavour-gen

```bash
flavour-gen create \
  --input ./my-clients \
  --output-dir /path/to/project/app/output
```

### 3. Check Output

```bash
ls -la /path/to/project/app/output/
```

---

## Multiple Clients

You can create multiple flavors at once using an array in data.json:

```bash
mkdir my-clients
cd my-clients

# Create data.json with multiple clients
cat > data.json << 'EOF'
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
EOF

# Add logos
cp /path/to/logo1.png logo1.png
cp /path/to/logo2.png logo2.png
```

Then run:

```bash
flavour-gen create \
  --input ./my-clients \
  --output-dir /path/to/project/app/output
```

This will create both flavors in the output directory.

---

## Dry Run Mode

Preview what would be created without writing files:

```bash
flavour-gen create \
  --input ./my-clients \
  --output-dir /path/to/project/app/output \
  --dry-run -v
```

---

## Options

### Custom Logo

```bash
flavour-gen create \
  --input ./my-clients \
  --logo /custom/path/logo.png \
  --output-dir ./output
```

### Custom Background Color

```bash
flavour-gen create \
  --input ./my-clients \
  --bg-color "#FF5722" \
  --output-dir ./output
```

---

## Next Steps

1. Open Android project in Android Studio
2. Build a flavor: `./gradlew assemble<ArchiveBasename>Debug`
3. Find APK in `app/build/outputs/apk/`

## Troubleshooting

See [Error Handling](./error-handling.md) for common issues.
