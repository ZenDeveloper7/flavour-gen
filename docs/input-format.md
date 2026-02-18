# Input Format

The input folder should contain:

```
input-folder/
├── data.json              # Required (single client or array)
├── google-services.json   # Required (provides education_number)
└── app_logo.png         # Required (or use --logo flag)
```

## data.json

The JSON file supports **both single client and multiple clients**:

### Single Client

```json
{
  "app_name": "My App",
  "archivebasename": "my_app",
  "package_name": "com.mycompany.myapp",
  "theme_id": 1,
  "app_logo": "app_logo.png",
  "base_url": "https://api.example.com/"
}
```

### Multiple Clients (Array)

```json
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
```

### Required Fields

| Field | Type | Description |
|-------|------|-------------|
| `app_name` | string | Display name of the app |
| `archivebasename` | string | Archive name (no spaces, use underscores) |
| `package_name` | string | Android package identifier |
| `theme_id` | integer | Theme ID to duplicate |
| `app_logo` | string | Path to logo file (relative to input folder) |
| `base_url` | string | API base URL |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `firebase_url` | string | - | Firebase URL |
| `dynamic_link_domain` | string | - | Dynamic Links domain |
| `dynamic_link_prefix` | string | - | Dynamic Links prefix |

## Logo Requirements

- **Format:** PNG
- **Size:** Recommended 512x512 or larger
- **Placement:** In input folder (path relative to data.json)

For multiple clients, each client can have its own logo in the input folder.

## google-services.json

**Required** - provides:
- `education_number` - Extracted from `project_info.project_id` or `project_number`

Copy from Firebase console or existing flavour.

## Computed Values

These are automatically calculated:

| Field | Formula |
|-------|---------|
| `identity` | base64(package_name) |
| `dot_count` | count of "." in package_name |
| `alt_app_name` | lowercase(app_name with spaces as underscores) |
| `download_folder` | archivebasename |
| `version_name` | "1.0.0" (default) |
| `version_code` | 1 (default) |
| `test_base_url` | base_url (default) |
| `education_number` | from google-services.json |
