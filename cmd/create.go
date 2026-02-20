package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"

	"github.com/ZenDeveloper7/flavour-gen/pkg/config"
	flavour "github.com/ZenDeveloper7/flavour-gen/pkg/flavour"
	"github.com/ZenDeveloper7/flavour-gen/pkg/icon"
	"github.com/ZenDeveloper7/flavour-gen/pkg/keystore"
)

var (
	inputPath string
	logoPath  string
	bgColor   string
	outputDir string
	dryRun    bool
	verbose   bool
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
	Short: "Create Android app flavors from client data",
	RunE:  runCreate,
}

func init() {
	createCmd.Flags().StringVar(&inputPath, "input", "", "Input folder with data.json [required]")
	createCmd.Flags().StringVar(&logoPath, "logo", "", "Logo image file (PNG, JPEG, WEBP, etc.) [optional]")
	createCmd.Flags().StringVar(&bgColor, "bg-color", "", "Background color #RRGGBB (default: white)")
	createCmd.Flags().StringVar(&outputDir, "output-dir", "./output", "Output directory")
	createCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview without writing files")
	createCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose logging")

	createCmd.MarkFlagRequired("input")
}

func runCreate(cmd *cobra.Command, args []string) error {
	if verbose {
		infoC.Println("[INFO] Starting flavour generation")
	}

	// Check prerequisites
	if err := checkPrerequisites(); err != nil {
		return err
	}

	// Validate input folder
	inputDir, err := os.Stat(inputPath)
	if err != nil {
		return fmt.Errorf("input folder: %w", err)
	}
	if !inputDir.IsDir() {
		return errors.New("input must be a folder")
	}

	// Load client data (supports single or array)
	clients, err := config.LoadClientData(filepath.Join(inputPath, "data.json"))
	if err != nil {
		return fmt.Errorf("load data.json: %w", err)
	}

	if len(clients) == 0 {
		return errors.New("no clients found in data.json")
	}

	// Find Android project root from output-dir
	androidProject := filepath.Dir(outputDir)
	if !filepath.IsAbs(outputDir) {
		absOutput, _ := filepath.Abs(outputDir)
		androidProject = filepath.Dir(absOutput)
	}
	if androidProject == "." || androidProject == "/" {
		return errors.New("output-dir must be inside an Android project")
	}

	// Set templates dir to Android project's app folder
	flavour.SetTemplatesDir(filepath.Join(androidProject, "app"))

	// Process each client
	successCount := 0
	for i, client := range clients {
		if len(clients) > 1 {
			infoC.Printf("\n[INFO] Processing client %d/%d: %s\n", i+1, len(clients), client.AppName)
		}

		// Validate required fields
		if client.ThemeID <= 0 {
			warnC.Printf("[WARN] Client '%s': theme_id is required, skipping\n", client.AppName)
			continue
		}
		if client.AppLogo == "" && logoPath == "" {
			warnC.Printf("[WARN] Client '%s': app_logo is required, skipping\n", client.AppName)
			continue
		}

		// Check google-services.json (mandatory)
		gsFile := filepath.Join(inputPath, "google-services.json")
		if _, err := os.Stat(gsFile); err != nil {
			warnC.Printf("[WARN] Client '%s': google-services.json is required, skipping\n", client.AppName)
			continue
		}

		// Extract education_number from google-services.json
		eduNum, err := config.ExtractEducationNumber(gsFile)
		if err != nil {
			warnC.Printf("[WARN] Client '%s': failed to extract education_number: %v\n", client.AppName, err)
			continue
		}
		client.EducationNumber = eduNum
		if verbose {
			infoC.Printf("[INFO] Client '%s': Education number: %d\n", client.AppName, eduNum)
		}

		// Resolve logo path
		logoFilePath := resolveLogoPath(inputPath, client.AppLogo, logoPath)
		if logoFilePath == "" {
			warnC.Printf("[WARN] Client '%s': logo not found, skipping\n", client.AppName)
			continue
		}
		// Note: Any image format supported by the imaging library is accepted
		// (PNG, JPEG, WEBP, GIF, TIFF, BMP, etc.)

		// Validate theme exists
		availableThemes, err := flavour.ListThemes()
		if err != nil {
			return fmt.Errorf("list themes: %w", err)
		}
		if !slices.Contains(availableThemes, client.ThemeID) {
			warnC.Printf("[WARN] Client '%s': theme %d not found, skipping\n", client.AppName, client.ThemeID)
			continue
		}

		// Extract education_number from google-services.json if available
		gsFile = filepath.Join(inputPath, "google-services.json")
		if _, err := os.Stat(gsFile); err == nil {
			if eduNum, err := config.ExtractEducationNumber(gsFile); err == nil {
				client.EducationNumber = eduNum
				if verbose {
					infoC.Printf("[INFO] Client '%s': Education number: %d\n", client.AppName, eduNum)
				}
			}
		}

		// Create output directory for this client
		clientOutputDir := outputDir

		if !dryRun {
			if err := os.MkdirAll(filepath.Join(clientOutputDir, "app/keystore"), 0755); err != nil {
				warnC.Printf("[WARN] Client '%s': failed to create directories: %v\n", client.AppName, err)
				continue
			}
		}

		// Generate icons
		bg, err := icon.GetBackgroundColor(bgColor, false, logoFilePath)
		if err != nil {
			warnC.Printf("[WARN] Client '%s': background color: %v\n", client.AppName, err)
			continue
		}

		if _, err := icon.GenerateAppIcon(logoFilePath, client.ArchiveBasename, clientOutputDir, bg, dryRun); err != nil {
			warnC.Printf("[WARN] Client '%s': app icon: %v\n", client.AppName, err)
			continue
		}

		if _, err := icon.GenerateNotificationIcon(logoFilePath, client.ArchiveBasename, clientOutputDir, dryRun); err != nil {
			warnC.Printf("[WARN] Client '%s': notification icon: %v\n", client.AppName, err)
			continue
		}

		if verbose {
			infoC.Printf("[INFO] Client '%s': Icons generated\n", client.AppName)
		}

		// Duplicate theme and create gradle
		var themeErr error
		_, themeErr = flavour.DuplicateTheme(client.ThemeID, client.ArchiveBasename, clientOutputDir, &client, dryRun)
		if themeErr != nil {
			warnC.Printf("[WARN] Client '%s': duplicate theme: %v\n", client.AppName, err)
			continue
		}

		if verbose {
			infoC.Printf("[INFO] Client '%s': Gradle created\n", client.AppName)
		}

		// Update build_type.gradle and flavours.gradle
		if !dryRun {
			if err := flavour.AddFlavorToBuildType(client.ArchiveBasename, androidProject, dryRun); err != nil {
				warnC.Printf("[WARN] Client '%s': build_type.gradle: %v\n", client.AppName, err)
			}
			if err := flavour.AddFlavorToFlavours(client.ArchiveBasename, androidProject, dryRun); err != nil {
				warnC.Printf("[WARN] Client '%s': flavours.gradle: %v\n", client.AppName, err)
			}
		}

		// Generate keystore
		keystorePath, err := keystore.Generate(&client, clientOutputDir, dryRun)
		if err != nil {
			warnC.Printf("[WARN] Client '%s': keystore: %v\n", client.AppName, err)
			continue
		}
		if verbose && keystorePath != "" {
			infoC.Printf("[INFO] Client '%s': Keystore: %s\n", client.AppName, keystorePath)
		}

		// Copy google-services.json
		gsDest := filepath.Join(clientOutputDir, "app/src", client.ArchiveBasename, "google-services.json")
		if _, err := os.Stat(gsFile); err == nil && !dryRun {
			if data, err := os.ReadFile(gsFile); err == nil {
				os.WriteFile(gsDest, data, 0644)
			}
		}

		boldC.Print("✅ ")
		successC.Printf("Flavor created: %s\n", client.ArchiveBasename)
		successCount++
	}

	// Summary
	if len(clients) > 1 {
		infoC.Printf("\n[INFO] Processed %d/%d clients successfully\n", successCount, len(clients))
	}

	if successCount == 0 {
		return errors.New("no flavors were created")
	}

	if dryRun {
		warnC.Println("(dry-run mode)")
	}

	return nil
}

// checkPrerequisites verifies required tools are installed
func checkPrerequisites() error {
	if _, err := exec.LookPath("keytool"); err != nil {
		return errors.New("keytool not found. Install Java JDK")
	}
	if _, err := exec.LookPath("gradle"); err != nil {
		warnC.Println("[WARN] Gradle not found (optional)")
	}
	return nil
}

// resolveLogoPath finds the logo file path
func resolveLogoPath(inputDir, appLogo, cliLogo string) string {
	if cliLogo != "" {
		return cliLogo
	}
	if appLogo != "" {
		path := filepath.Join(inputDir, appLogo)
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	return ""
}
