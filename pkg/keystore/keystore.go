package keystore

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/ZenDeveloper7/flavour-gen/pkg/config"
)

// Generate creates a Java keystore (JKS) for the given client.
// It calls the external `keytool` command.
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
	return keystorePath, nil
}