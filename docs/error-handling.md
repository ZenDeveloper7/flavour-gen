# Error Handling

## Common Errors

### keytool not found

**Error:**
```
keytool not found. Install Java JDK
```

**Solution:** Install Java JDK which includes keytool.

```bash
# Ubuntu/Debian
sudo apt install openjdk-17-jdk

# macOS
brew install openjdk

# Windows
# Download from Oracle or use Chocolatey
choco install openjdk17
```

---

### theme_id is required

**Error:**
```
theme_id is required in data.json
```

**Solution:** Add `theme_id` to your data.json:

```json
{
  "theme_id": 1
}
```

---

### app_logo is required

**Error:**
```
app_logo is required in data.json
```

**Solution:** Add `app_logo` to your data.json:

```json
{
  "app_logo": "app_logo.png"
}
```

Ensure the file exists in the input folder.

---

### Theme X not found

**Error:**
```
theme 1 not found in project. available: []
```

**Solution:** Ensure the theme exists in your Android project's app folder:

```
app/
└── flavours/
    └── appx_theme1.gradle   # Must exist
```

---

### Logo must be PNG

**Error:**
```
logo must be PNG
```

**Solution:** Convert your logo to PNG format:

```bash
# Using ImageMagick
convert input.jpg output.png

# Using ffmpeg
ffmpeg -i input.jpg output.png
```

---

### output-dir must be inside an Android project

**Error:**
```
output-dir must be inside an Android project
```

**Solution:** The output directory should be inside or under your Android project:

```bash
# Correct
--output-dir ./AndroidProject/app/output

# Incorrect (not in project)
--output-dir ./random-folder
```

---

### Cannot write to output directory

**Error:**
```
Cannot write to output directory
```

**Solution:** Check permissions or use a writable location:

```bash
chmod 755 ./output
# or
mkdir -p ./output
```

---

## Verbose Mode

For debugging, use verbose mode:

```bash
flavour-gen create --input ./client --output-dir ./output -v
```

This shows detailed logs of what the CLI is doing.

---

## Dry Run

Test without creating files:

```bash
flavour-gen create --input ./client --output-dir ./output --dry-run -v
```

This validates inputs and shows what would be created without writing files.
