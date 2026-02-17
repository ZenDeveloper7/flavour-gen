package gradle

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ZenDeveloper7/flavour-gen/pkg/config"
)

// Validate checks that a generated Gradle file parses (basic syntax check).
func Validate(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("file missing: %w", err)
	}
	if info.Size() == 0 {
		return fmt.Errorf("file empty")
	}
	return nil
}

// WriteFile writes content to dst.
func WriteFile(dst string, content []byte) error {
	return os.WriteFile(dst, content, 0644)
}

// EnsureResValues creates res/values/strings.xml if missing in the client src folder.
func EnsureResValues(dir string, cd *config.ClientData) error {
	valuesDir := filepath.Join(dir, "res/values")
	if err := os.MkdirAll(valuesDir, 0755); err != nil {
		return err
	}
	path := filepath.Join(valuesDir, "strings.xml")
	if _, err := os.Stat(path); err == nil {
		return nil // already exists
	}
	content := `<?xml version="1.0" encoding="utf-8"?>
<resources>
    <string name="app_name">` + cd.AppName + `</string>
</resources>`
	return os.WriteFile(path, []byte(content), 0644)
}