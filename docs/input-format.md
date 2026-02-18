# Input Format

The input folder should contain:

```
input-folder/
├── data.json              # Required
├── google-services.json   # Optional (provides education_number)
└── app_logo.png         # Required (or use --logo flag)
```

## data.json

### Required Fields

| Field | Type | Description |
|-------|------|-------------|
| `app_name` | string | Display name of the app |
| `archivebasename` | string | Archive name (no spaces, use underscores) |
| `package_name` | string | Android package identifier |
| `theme_id` | integer | Theme ID to duplicate |
| `app_logo` | string | Path to logo file (relative to input folder) |
| `base_url` | string | API base URL (used for both production and test) |

### Optional Fields

| Field | Type | Description |
|-------|------|-------------|
| `firebase_url` | string | Firebase database URL |
| `dynamic_link_domain` | string | Dynamic Links domain |
| `dynamic_link_prefix` | string | Dynamic Links prefix |

### Default Values

If not provided:

| Field | Default |
|-------|---------|
| `version_name` | "1.0.0" |
| `version_code` | 0 |
| `test_base_url` | Same as `base_url` |
| `education_number` | Extracted from google-services.json |

## Example

```json
{
  "app_name": "My App",
  "archivebasename": "my_app",
  "package_name": "com.mycompany.myapp",
  "theme_id": 1,
  "app_logo": "app_logo.png",
  "base_url": "https://api.example.com/",
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

Optional but provides:
- `education_number` - Extracted from `project_info.project_number`

Copy from Firebase console or existing flavour.

## Computed Values

These are automatically calculated:

| Field | Formula |
|-------|---------|
| `identity` | base64(package_name) |
| `dot_count` | count of "." in package_name |
| `alt_app_name` | lowercase(app_name with spaces as underscores) |
| `download_folder | archivebasename |
| `test_base_url` | base_url (if not provided) |
| `education_number` | project_number from google-services.json |
