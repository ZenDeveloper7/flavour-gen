package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"

	"github.com/ZenDeveloper7/flavour-gen/pkg/config"
	"github.com/ZenDeveloper7/flavour-gen/pkg/flavor"
	"github.com/ZenDeveloper7/flavour-gen/pkg/icon"
	"github.com/ZenDeveloper7/flavour-gen/pkg/keystore"
)

var (
	clientDataPath string
	themeID       int
	logoPath      string
	bgColor       string
	autoBG        bool
	outputDir     string
	dryRun        bool
	verbose       bool
)

// Color helpers
var (
	infoC    = color.New(color.FgCyan)
	warnC    = color.New(color.FgYellow)
	successC = color.New(color.FgGreen)
	errorC   = color.New(color.FgRed)
	boldC    = color.New(color.Bold)
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new app flavor package",
	RunE:  runCreate,
}

func init() {
	createCmd.Flags().StringVar(&clientDataPath, "client-data", "", "Client data file (JSON) [required]")
	createCmd.Flags().IntVar(&themeID, "theme-id", 0, "Theme ID [required]")
	createCmd.Flags().StringVar(&logoPath, "logo", "", "Logo PNG file [required]")
	createCmd.Flags().StringVar(&bgColor, "bg-color", "", "Background color #RRGGBB (auto-detected if empty)")
	createCmd.Flags().BoolVar(&autoBG, "auto-bg", true, "Auto-detect background from logo")
	createCmd.Flags().StringVar(&outputDir, "output-dir", "./output", "Output directory")
	createCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview without writing files")
	createCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose logging")

	createCmd.MarkFlagRequired("client-data")
	createCmd.MarkFlagRequired("theme-id")
	createCmd.MarkFlagRequired("logo")
}

func runCreate(cmd *cobra.Command, args []string) error {
	if verbose {
		infoC.Println("[INFO] Starting flavor generation")
	}

	// Step 1: Prerequisites
	if err := checkPrerequisites(); err != nil {
		return fmt.Errorf("prerequisites check failed: %w", err)
	}

	// Step 2: Validate inputs
	clientData, err := loadClientData(clientDataPath)
	if err != nil {
		return fmt.Errorf("client data: %w", err)
	}

	// Validate theme exists
	availableThemes, err := flavor.ListThemes()
	if err != nil {
		return fmt.Errorf("list themes: %w", err)
	}
	if !slices.Contains(availableThemes, themeID) {
		return fmt.Errorf("theme %d not found. available: %v", themeID, availableThemes)
	}

	// Validate logo
	if _, err := os.Stat(logoPath); err != nil {
		return fmt.Errorf("logo file: %w", err)
	}
	if ext := strings.ToLower(filepath.Ext(logoPath)); ext != ".png" {
		return errors.New("logo must be PNG")
	}

	// Validate output directory
	if stat, err := os.Stat(outputDir); err == nil && !stat.IsDir() {
		return errors.New("output path exists and is not a directory")
	}

	if verbose {
		infoC.Printf("[INFO] Client: %s (%s)\n", clientData.AppName, clientData.ArchiveBasename)
		infoC.Printf("[INFO] Theme: %d\n", themeID)
		infoC.Printf("[INFO] Output: %s\n", outputDir)
	}

	// Step 3: Create folder structure
	if !dryRun {
		if err := os.MkdirAll(filepath.Join(outputDir, "app/keystore"), 0755); err != nil {
			return fmt.Errorf("create keystore dir: %w", err)
		}
		if err := os.MkdirAll(filepath.Join(outputDir, "app/flavours"), 0755); err != nil {
			return fmt.Errorf("create flavours dir: %w", err)
		}
		if err := os.MkdirAll(filepath.Join(outputDir, "app/src", clientData.ArchiveBasename), 0755); err != nil {
			return fmt.Errorf("create src dir: %w", err)
		}
	}

	// Step 4: Process icons
	if verbose {
		infoC.Println("[INFO] Processing icons...")
	}
	bg, err := icon.GetBackgroundColor(bgColor, autoBG, logoPath)
	if err != nil {
		return fmt.Errorf("background color: %w", err)
	}

	iconPaths, err := icon.GenerateAppIcon(logoPath, clientData.ArchiveBasename, outputDir, bg, dryRun)
	if err != nil {
		return fmt.Errorf("app icon: %w", err)
	}

	notifPath, err := icon.GenerateNotificationIcon(logoPath, clientData.ArchiveBasename, outputDir, dryRun)
	if err != nil {
		return fmt.Errorf("notification icon: %w", err)
	}

	if verbose {
		infoC.Printf("[INFO] Icons written: %s, %s\n", iconPaths.AdaptiveXML, notifPath)
	}

	// Step 5: Duplicate theme and generate Gradle files
	if verbose {
		infoC.Println("[INFO] Duplicating theme...")
	}
	themeFiles, err := flavor.DuplicateTheme(themeID, clientData.ArchiveBasename, outputDir, clientData, dryRun)
	if err != nil {
		return fmt.Errorf("duplicate theme: %w", err)
	}

	if verbose {
		infoC.Printf("[INFO] Theme files: %v\n", themeFiles.GradleFile)
	}

	// Step 6: Generate keystore
	if verbose {
		infoC.Println("[INFO] Generating keystore...")
	}
	keystorePath, err := keystore.Generate(clientData, outputDir, dryRun)
	if err != nil {
		return fmt.Errorf("keystore: %w", err)
	}
	if verbose && keystorePath != "" {
		infoC.Printf("[INFO] Keystore: %s\n", keystorePath)
	}

	// Step 7: Copy google-services.json if available and firebase URL provided
	if clientData.FirebaseURL != "" {
		src := filepath.Join("templates", fmt.Sprintf("appx_theme%d_sample", themeID), "google-services.json")
		dst := filepath.Join(outputDir, "app/src", clientData.ArchiveBasename, "google-services.json")
		if _, err := os.Stat(src); err == nil {
			if !dryRun {
				if err := copyFile(src, dst); err != nil {
					return fmt.Errorf("copy google-services.json: %w", err)
				}
			}
			if verbose {
				infoC.Println("[INFO] google-services.json copied")
			}
		} else if verbose {
			warnC.Println("[WARN] google-services.json not found in theme template")
		}
	}

	// Success output
	boldC.Print("✅ ")
	successC.Printf("Flavor created: %s\n", clientData.ArchiveBasename)
	if dryRun {
		warnC.Println("(dry-run mode — no files written)")
	}
	return nil
}

func checkPrerequisites() error {
	// Check keytool exists (required)
	if _, err := exec.LookPath("keytool"); err != nil {
		return errors.New("keytool not installed. Install Java JDK")
	}
	// Check gradle exists (optional - only warn)
	if _, err := exec.LookPath("gradle"); err != nil {
		warnC.Println("[WARN] Gradle not found (optional - only needed for validation)")
	}
	return nil
}

func loadClientData(path string) (*config.ClientData, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cd config.ClientData
	if err := json.Unmarshal(data, &cd); err != nil {
		return nil, err
	}
	cd.Compute()
	return &cd, nil
}

func copyFile(src, dst string) error {
	in, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, in, 0644)
}
