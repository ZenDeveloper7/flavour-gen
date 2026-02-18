package keystore

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ZenDeveloper7/flavour-gen/pkg/config"
)

// Generate creates a Java keystore (JKS) and adds entry to keystore.gradle
func Generate(cd *config.ClientData, outputDir string, dryRun bool) (string, error) {
	keystoreDir := filepath.Join(outputDir, "app/keystore")
	keystorePath := filepath.Join(keystoreDir, cd.ArchiveBasename+".jks")

	if dryRun {
		return keystorePath, nil
	}

	if err := os.MkdirAll(keystoreDir, 0755); err != nil {
		return "", err
	}

	// Generate keystore using keytool
	cmd := exec.Command("keytool",
		"-genkeypair", "-v", "-storetype", "JKS", "-keyalg", "RSA",
		"-keysize", "2048", "-validity", "10000",
		"-keystore", keystorePath,
		"-alias", cd.ArchiveBasename,
		"-storepass", cd.ArchiveBasename,
		"-keypass", cd.ArchiveBasename,
		"-dname", fmt.Sprintf("CN=%s, OU=%s, O=%s, L=Delhi, ST=Delhi, C=IN", cd.ArchiveBasename,cd.ArchiveBasename,cd.ArchiveBasename),
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("keytool failed: %w", err)
	}

	// Add entry to keystore.gradle
	if err := addToKeystoreGradle(cd.ArchiveBasename, filepath.Dir(outputDir)); err != nil {
		fmt.Printf("[WARN] Failed to update keystore.gradle: %v\n", err)
	}

	return keystorePath, nil
}

// addToKeystoreGradle adds a new signing config entry to keystore.gradle
func addToKeystoreGradle(archiveBasename, projectRoot string) error {
	keystoreGradlePath := filepath.Join(projectRoot, "app/keystore.gradle")
	if _, err := os.Stat(keystoreGradlePath); err != nil {
		return fmt.Errorf("keystore.gradle not found: %w", err)
	}

	content, err := os.ReadFile(keystoreGradlePath)
	if err != nil {
		return err
	}

	// Skip if already exists
	if strings.Contains(string(content), archiveBasename+" {") {
		return nil
	}

	// Add new signing config after existing ones
	newEntry := fmt.Sprintf(`
    %s {
        storeFile file("keystore/%s.jks")
        storePassword "%s"
        keyAlias "%s"
        keyPassword "%s"
    }
`, archiveBasename, archiveBasename, archiveBasename, archiveBasename, archiveBasename)

	lines := strings.Split(string(content), "\n")
	var newLines []string
	added := false

	for i, line := range lines {
		newLines = append(newLines, line)
		// After closing brace of a signing config, add our new config
		if strings.Contains(line, "        }") && !added {
			if i+1 < len(lines) && strings.Contains(lines[i+1], archiveBasename) {
				continue
			}
			newLines = append(newLines, newEntry)
			added = true
		}
	}

	if !added {
		newLines = append(newLines, newEntry)
	}

	return os.WriteFile(keystoreGradlePath, []byte(strings.Join(newLines, "\n")), 0644)
}
