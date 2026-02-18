# Configuration Package

`pkg/config/config.go` handles client data parsing and validation.

## Overview

Defines the `ClientData` struct and computes derived fields from input JSON.

## ClientData Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `app_name` | string | Yes | Display name |
| `archivebasename` | string | Yes | Archive name (e.g., `my_app`) |
| `package_name` | string | Yes | Android package |
| `theme_id` | string | Yes | Theme ID to use |
| `app_logo` | string | Yes | Path to logo file |
| `base_url` | string | Yes | API base URL |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `firebase_url` | string | - | Firebase URL |
| `dynamic_link_domain` | string | - | Dynamic Links domain |
| `dynamic_link_prefix` | string | - | Dynamic Links prefix |

## Computed Fields

These are automatically calculated:

| Field | Description | Formula |
|-------|-------------|---------|
| `identity` | Base64 encoded package | base64(package_name) |
| `dot_count` | Dot count in package | count("." in package_name) |
| `alt_app_name` | Lowercase with underscores | lowercase(spaces to underscores) |
| `download_folder` | Archive basename | archivebasename |
| `version_name` | Version string | "1.0.0" (default) |
| `version_code` | Version code | 0 (default) |
| `test_base_url` | Test API URL | base_url (default) |
| `education_number` | From google-services.json | project_number |

## Functions

### `Compute()`

Calculates derived fields from input data.

### `Validate() error`

Validates required fields. Returns error if any required field is missing.

### `ExtractEducationNumber(path string) (int, error)`

Extracts education_number from google-services.json file.

## Example

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
