# Input Format

The input folder should contain:

```
input-folder/
├── data.json              # Required
├── google-services.json   # Optional
└── app_logo.png         # Required (or use --logo flag)
```

## data.json

### Required Fields

| Field | Type | Description |
|-------|------|-------------|
| `app_name` | string | Display name of the app |
| `archivebasename` | string | Archive name (no spaces, use underscores) |
| `package_name` | string | Android package identifier |
| `version_name` | string | Version string (e.g., "1.0.0") |
| `version_code` | integer | Version code (integer) |
| `theme_id` | integer | Theme ID to duplicate |
| `app_logo` | string | Path to logo file (relative to input folder) |
| `base_url` | string | Production API base URL |
| `test_base_url` | string | Test API base URL |

### Optional Fields

| Field | Type | Description |
|-------|------|-------------|
| `education_number` | integer | Education number for gradle comment |
| `firebase_url` | string | Firebase database URL |
| `dynamic_link_domain` | string | Dynamic Links domain |
| `dynamic_link_prefix` | string | Dynamic Links prefix |

## Example

```json
{
  "app_name": "My App",
  "archivebasename": "my_app",
  "package_name": "com.ydcfzb.zgizxw",
  "version_name": "1.0.0",
  "version_code": 0,
  "theme_id": 1,
  "app_logo": "app_logo.png",
  "education_number": 20,
  "base_url": "https://api.example.com/",
  "test_base_url": "https://api.example.com/",
  "firebase_url": "https://myapp.firebaseio.com/",
  "dynamic_link_domain": "https://example.com/",
  "dynamic_link_prefix": "https://myapp.page.link"
}
```

## Logo Requirements

- **Format:** PNG
- **Size:** Recommended 512x512 or larger
- **Placement:** In input folder (path relative to data.json)

## google-services.json

Optional but required if using Firebase features.

Copy from Firebase console or existing flavor.

## Computed Values

These are automatically calculated from input:

| Field | Formula |
|-------|---------|
| `identity` | base64(package_name) |
| `dot_count` | count of "." in package_name |
| `alt_app_name` | lowercase(app_name with spaces as underscores) |
| `download_folder | archivebasename |
