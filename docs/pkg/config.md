# Configuration Package

`pkg/config/config.go` handles client data parsing and validation.

## Overview

Defines the `ClientData` struct and computes derived fields from input JSON.

## ClientData Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `app_name` | string | Yes | Display name |
| `archivebasename` | string | Yes | Archive name (e.g., `physics_setu`) |
| `package_name` | string | Yes | Android package |
| `version_name` | string | Yes | Version string |
| `version_code` | int | Yes | Version code |
| `theme_id` | int | Yes | Theme ID to use |
| `app_logo` | string | Yes | Path to logo file |
| `education_number` | int | No | Education number for gradle |
| `base_url` | string | Yes | Production API URL |
| `test_base_url` | string | Yes | Test API URL |
| `firebase_url` | string | No | Firebase URL |
| `dynamic_link_domain` | string | No | Dynamic Links domain |
| `dynamic_link_prefix` | string | No | Dynamic Links prefix |

## Computed Fields

These are automatically calculated:

| Field | Description | Formula |
|-------|-------------|---------|
| `identity` | Base64 encoded package | base64(package_name) |
| `dot_count` | Dot count in package | count("." in package_name) |
| `alt_app_name` | Lowercase with underscores | lowercase(spaces to underscores) |
| `download_folder` | Archive basename | archivebasename |

## Functions

### `Compute()`

Calculates derived fields from input data.

### `Validate() error`

Validates required fields. Returns error if any required field is missing.

## Example

```json
{
  "app_name": "Physics Setu",
  "archivebasename": "physics_setu",
  "package_name": "com.ydcfzb.zgizxw",
  "theme_id": 1,
  "app_logo": "app_logo.png",
  "base_url": "https://api.example.com/"
}
```
