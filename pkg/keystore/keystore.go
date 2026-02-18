package keystore

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ZenDeveloper7/flavour-gen/pkg/config"
)

// Generate creates a Java keystore (JKS) for the given client.
// It calls the external `keytool` command and adds an entry to keystore.gradle.
// Returns the path to the keystore file.
func Generate(cd *config.ClientData, outputDir string, dryRun bool) (string, error) {
	keystoreDir := filepath.Join(outputDir, "app/keystore")
	keystorePath := filepath.Join(keystoreDir, cd.ArchiveBasename+".jks")
	
	if dryRun {
		return keystorePath, nil
	}
	
	if err := os.MkdirAll(keystoreDir, 0755); err != nil {
		return "", err
	}
	
	// Build keytool command
	cmd := exec.Command("keytool",
		"-genkeypair",
		"-v",
		"-storetype", "JKS",
		"-keyalg", "RSA",
		"-keysize", "2048",
		"-validity", "10000",
		"-keystore", keystorePath,
		"-alias", cd.ArchiveBasename,
		"-storepass", cd.ArchiveBasename,
		"-keypass", cd.ArchiveBasename,
		"-dname", fmt.Sprintf("CN=%s, OU=Development, O=Company, L=City, ST=State, C=US", cd.ArchiveBasename),
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("keytool failed: %w", err)
	}
	
	// Add entry to keystore.gradle
	// Find the project root (parent of outputDir)
	projectRoot := filepath.Dir(outputDir)
	keystoreGradlePath := filepath.Join(projectRoot, "app/keystore.gradle")
	
	if _, err := os.Stat(keystoreGradlePath); err == nil {
		// Read existing keystore.gradle
		content, err := os.ReadFile(keystoreGradlePath)
		if err == nil {
			// Check if this signing config already exists
			configName := cd.ArchiveBasename
			if !strings.Contains(string(content), configName+" {") {
				// Add new signing config
				newEntry := fmt.Sprintf(`
        %s {
            storeFile file("keystore/%s.jks")
            storePassword "%s"
            keyAlias "%s"
            keyPassword "%s"
        }
`, configName, configName, configName, configName, configName)
				
				// Find position to insert - before the closing brace of signingConfigs
				existing := string(content)
				closingPos := strings.LastIndex(existing, "    }")
				if closingPos > 0 {
					newContent := existing[:closingPos] + newEntry + existing[closingPos:]
					err = os.WriteFile(keystoreGradlePath, []byte(newContent), 0644)
					if err != nil {
						return keystorePath, fmt.Errorf("update keystore.gradle: %w", err)
					}
				}
			}
		}
	}
	
	return keystorePath, nil
}
