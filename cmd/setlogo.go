package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ZenDeveloper7/flavour-gen/pkg/icon"
)

var (
	logoFile string
	flavourName string
)

var setLogoCmd = &cobra.Command{
	Use:   "set-logo",
	Short: "Update logo for an existing flavour",
	Args:  cobra.ExactArgs(2),
	RunE:  runSetLogo,
}

func init() {
	setLogoCmd.Flags().StringVar(&bgColor, "bg-color", "", "Background color #RRGGBB (default: white)")
	setLogoCmd.Flags().StringVar(&outputDir, "output-dir", "", "Output directory (required)")
	setLogoCmd.MarkFlagRequired("output-dir")
}

func runSetLogo(cmd *cobra.Command, args []string) error {
	logoFile = args[0]
	flavourName = args[1]

	if verbose {
		infoC.Printf("[INFO] Setting logo for flavour: %s\n", flavourName)
	}

	// Validate logo file exists
	if _, err := os.Stat(logoFile); err != nil {
		return fmt.Errorf("logo file not found: %w", err)
	}
	// Note: Any image format supported by the imaging library is accepted
	// (PNG, JPEG, WEBP, GIF, TIFF, BMP, etc.)

	// Find Android project from output-dir or use current directory
	// For set-logo, we need to find the project automatically
	// Check current dir first, then walk up to find app/ folder
	projectRoot, err := findProjectRoot()
	if err != nil {
		return fmt.Errorf("could not find Android project. Use --output-dir or run from project root: %w", err)
	}

	if verbose {
		infoC.Printf("[INFO] Found project: %s\n", projectRoot)
	}

	// Check if flavour exists
	flavourSrcDir := filepath.Join(projectRoot, "app/src", flavourName)
	if _, err := os.Stat(flavourSrcDir); os.IsNotExist(err) {
		return fmt.Errorf("flavour '%s' not found in app/src/", flavourName)
	}

	// Check if flavours.gradle has this flavour
	flavoursGradle := filepath.Join(projectRoot, "app/flavours.gradle")
	if _, err := os.Stat(flavoursGradle); err == nil {
		content, _ := os.ReadFile(flavoursGradle)
		if !strings.Contains(string(content), fmt.Sprintf("flavours/%s.gradle", flavourName)) {
			return fmt.Errorf("flavour '%s' not found in flavours.gradle", flavourName)
		}
	}

	if verbose {
		infoC.Printf("[INFO] Flavour exists, generating icons...\n")
	}

	// Get background color
	bg, err := icon.GetBackgroundColor(bgColor, false, logoFile)
	if err != nil {
		return fmt.Errorf("background color: %w", err)
	}

	// Output directory
	outputDir := filepath.Join(projectRoot, "app")

	// Generate app icon
	iconPaths, err := icon.GenerateAppIcon(logoFile, flavourName, outputDir, bg, dryRun)
	if err != nil {
		return fmt.Errorf("app icon: %w", err)
	}

	// Generate notification icon
	notifPath, err := icon.GenerateNotificationIcon(logoFile, flavourName, outputDir, dryRun)
	if err != nil {
		return fmt.Errorf("notification icon: %w", err)
	}

	// Also update ic_launcher_playstore.png
	// It's already created in GenerateAppIcon

	boldC.Print("✅ ")
	successC.Printf("Logo updated for: %s\n", flavourName)

	if verbose {
		infoC.Printf("[INFO] Updated: %s\n", iconPaths.AppLogo)
		infoC.Printf("[INFO] Updated: %s\n", notifPath)
		infoC.Printf("[INFO] Updated: %s\n", iconPaths.Playstore)
	}

	return nil
}

// findProjectRoot looks for the Android project root
func findProjectRoot() (string, error) {
	// First try output-dir if provided
	if outputDir != "" {
		if filepath.IsAbs(outputDir) {
			return filepath.Dir(filepath.Dir(outputDir)), nil
		}
		absOutput, _ := filepath.Abs(outputDir)
		return filepath.Dir(filepath.Dir(absOutput)), nil
	}

	// If no output-dir, error - user should specify it
	return "", errors.New("please specify --output-dir")
}
