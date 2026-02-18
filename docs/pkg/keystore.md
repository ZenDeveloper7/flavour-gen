# Keystore Package

`pkg/keystore/keystore.go` handles Java keystore generation.

## Overview

Creates a JKS keystore file for Android app signing and updates `keystore.gradle`.

## Functions

### `Generate(cd *config.ClientData, outputDir string, dryRun bool) (string, error)`

Generates a Java keystore (JKS) for the flavor.

**Process:**
1. Create keystore directory
2. Run `keytool` command to generate JKS
3. Add entry to project's `keystore.gradle`

**Keytool Command:**
```bash
keytool -genkeypair \
  -v -storetype JKS -keyalg RSA \
  -keysize 2048 -validity 10000 \
  -keystore <name>.jks \
  -alias <name> \
  -storepass <name> \
  -keypass <name> \
  -dname "CN=<name>, OU=Development, O=Company, L=City, ST=State, C=US"
```

**Defaults:**
- Algorithm: RSA
- Key size: 2048
- Validity: 10000 days
- Store password: archivebasename
- Key password: archivebasename

## keystore.gradle Update

Adds a new signing config block:

```groovy
<name> {
    storeFile file("keystore/<name>.jks")
    storePassword "<name>"
    keyAlias "<name>"
    keyPassword "<name>"
}
```

## Output

```
output/app/keystore/
└── <archivebasename>.jks
```

## Prerequisites

- Java JDK installed (keytool command available)

## Related Files

- `app/keystore.gradle` - Updated with new signing config
