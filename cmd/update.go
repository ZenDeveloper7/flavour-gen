package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Check for and install updates",
	Long: `Check the latest release on GitHub and update flavour-gen.

This command:
- Fetches the latest release from GitHub
- Compares with current version
- Downloads and installs if newer version available
- Works on Linux, macOS, and Windows`,
	RunE: runUpdate,
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

// GitHubRelease represents a GitHub release API response
type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
	HTMLURL string `json:"html_url"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

func runUpdate(cmd *cobra.Command, args []string) error {
	boldC := color.New(color.Bold)
	infoC := color.New(color.FgCyan)
	successC := color.New(color.FgGreen)
	warnC := color.New(color.FgYellow)

	boldC.Println("\n🔄 Flavour Gen Updater")
	boldC.Println("======================\n")

	// Get current version
	currentVersion := GetVersion()
	infoC.Printf("Current version: %s\n", currentVersion)

	// Detect OS and architecture
	osName := runtime.GOOS
	arch := runtime.GOARCH
	if arch == "amd64" && runtime.GOARCH == "x86_64" {
		arch = "amd64" // standardize
	}

	infoC.Printf("Platform: %s/%s\n", osName, arch)

	// Fetch latest release from GitHub
	infoC.Println("\nChecking for updates...")
	latestRelease, err := fetchLatestRelease()
	if err != nil {
		return fmt.Errorf("failed to fetch latest release: %w", err)
	}

	if latestRelease == nil {
		infoC.Println("No releases found on GitHub.")
		return nil
	}

	latestVersion := latestRelease.TagName
	infoC.Printf("Latest version: %s\n", latestVersion)
	infoC.Printf("Release page: %s\n", latestRelease.HTMLURL)

	// Compare versions - simple string comparison for now
	if currentVersion == latestVersion {
		successC.Println("\n✅ You're already on the latest version!")
		return nil
	}

	// Warn if current version is "dev"
	if currentVersion == "dev" {
		warnC.Println("\n⚠️  Current version is 'dev' (not a release build).")
	}

	// Find matching binary for current platform
	binaryName := fmt.Sprintf("flavour-gen-%s-%s", osName, arch)
	if osName == "windows" {
		binaryName += ".exe"
	}

	var downloadURL string
	for _, asset := range latestRelease.Assets {
		if asset.Name == binaryName {
			downloadURL = asset.BrowserDownloadURL
			break
		}
	}

	if downloadURL == "" {
		return fmt.Errorf("no binary found for %s", binaryName)
	}

	// Ask for confirmation
	fmt.Println()
	boldC.Println("Update available!")
	fmt.Printf("  From: %s\n", currentVersion)
	fmt.Printf("  To:   %s\n", latestVersion)
	fmt.Printf("  URL:  %s\n", downloadURL)
	fmt.Println()

	// In a real implementation, we would prompt for confirmation here
	// For now, auto-proceed (user can cancel with Ctrl+C)
	fmt.Println("Downloading and installing update...")

	// Download new binary
	tmpFile := filepath.Join(os.TempDir(), binaryName)
	if err := downloadFile(downloadURL, tmpFile); err != nil {
		return fmt.Errorf("download failed: %w", err)
	}

	// Make executable (skip on Windows)
	if osName != "windows" {
		if err := os.Chmod(tmpFile, 0755); err != nil {
			return fmt.Errorf("chmod failed: %w", err)
		}
	}

	// Determine install location
	installDir, err := getInstallDir()
	if err != nil {
		return fmt.Errorf("cannot determine install location: %w", err)
	}

	installPath := filepath.Join(installDir, "flavour-gen")
	if osName == "windows" {
		installPath += ".exe"
	}

	// Replace existing binary
	if err := os.Rename(tmpFile, installPath); err != nil {
		// If rename fails (different filesystems), try copy then delete
		if err2 := copyFile(tmpFile, installPath); err2 != nil {
			return fmt.Errorf("install failed (rename: %v, copy: %v)", err, err2)
		}
		os.Remove(tmpFile)
	}

	successC.Printf("\n✅ Updated successfully to %s!\n", latestVersion)
	infoC.Printf("Installed to: %s\n", installPath)

	return nil
}

func fetchLatestRelease() (*GitHubRelease, error) {
	repoURL := "https://api.github.com/repos/ZenDeveloper7/flavour-gen/releases/latest"
	resp, err := http.Get(repoURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil // No releases
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var release GitHubRelease
	if err := json.Unmarshal(body, &release); err != nil {
		return nil, err
	}

	return &release, nil
}

func downloadFile(url, destPath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func getInstallDir() (string, error) {
	// Check common binary paths
	possibleDirs := []string{
		"/usr/local/bin",
		"/usr/bin",
		filepath.Join(os.Getenv("HOME"), ".local", "bin"),
	}

	// If a custom install location is known, could be stored in config
	// For now, try to find existing binary and replace in place
	existing, err := exec.LookPath("flavour-gen")
	if err == nil {
		// Replace at existing location
		dir := filepath.Dir(existing)
		if dir != "" && dir != "." {
			return dir, nil
		}
	}

	// Fallback to /usr/local/bin (may need sudo)
	for _, dir := range possibleDirs {
		if _, err := os.Stat(dir); err == nil {
			// Return first existing directory; handle write errors during install
			return dir, nil
		}
	}

	return "", fmt.Errorf("no suitable install directory found. Install manually from: https://github.com/ZenDeveloper7/flavour-gen/releases")
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}
