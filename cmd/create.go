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
	flavour "github.com/ZenDeveloper7/flavour-gen/pkg/flavour"
	"github.com/ZenDeveloper7/flavour-gen/pkg/icon"
	"github.com/ZenDeveloper7/flavour-gen/pkg/keystore"
)

var (
	inputPath  string
	logoPath   string
	bgColor    string
	autoBG     bool
	outputDir  string
	dryRun     bool
	verbose    bool
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
	Short: "Create a new Android app flavor",
	RunE:  runCreate,
}

func init() {
	createCmd.Flags().StringVar(&inputPath, "input", "", "Input folder with data.json and google-services.json [required]")
	createCmd.Flags().StringVar(&logoPath, "logo", "", "Logo PNG file [optional - uses app_logo from data.json]")
	createCmd.Flags().StringVar(&bgColor, "bg-color", "", "Background color #RRGGBB (auto-detected if empty)")
	createCmd.Flags().BoolVar(&autoBG, "auto-bg", true, "Auto-detect background from logo")
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

	// Load and validate data.json
	clientData, err := loadClientData(filepath.Join(inputPath, "data.json"))
	if err != nil {
		return fmt.Errorf("load data.json: %w", err)
	}

	// Validate required fields
	if clientData.ThemeID <= 0 {
		return errors.New("theme_id is required in data.json")
	}
	if clientData.AppLogo == "" {
		return errors.New("app_logo is required in data.json")
	}

	// Resolve logo path - from data.json or --logo flag
	logoPath = resolveLogoPath(inputPath, clientData.AppLogo, logoPath)
	if logoPath == "" {
		return errors.New("logo not found. Use --logo or specify app_logo in data.json")
	}
	if ext := strings.ToLower(filepath.Ext(logoPath)); ext != ".png" {
		return errors.New("logo must be PNG")
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

	// Validate theme exists
	availableThemes, err := flavour.ListThemes()
	if err != nil {
		return fmt.Errorf("list themes: %w", err)
	}
	if !slices.Contains(availableThemes, clientData.ThemeID) {
		return fmt.Errorf("theme %d not found in project. available: %v", clientData.ThemeID, availableThemes)
	}

	// Check for google-services.json
	gsFile := filepath.Join(inputPath, "google-services.json")
	_, err = os.Stat(gsFile)
	hasGS := err == nil

	if verbose {
		infoC.Printf("[INFO] Client: %s (%s)\n", clientData.AppName, clientData.ArchiveBasename)
		infoC.Printf("[INFO] Theme: %d | Output: %s | GS: %v\n", clientData.ThemeID, outputDir, hasGS)
	}

	// Create output directory structure
	if !dryRun {
		os.MkdirAll(filepath.Join(outputDir, "app/keystore"), 0755)
		os.MkdirAll(filepath.Join(outputDir, "app/flavours"), 0755)
		os.MkdirAll(filepath.Join(outputDir, "app/src", clientData.ArchiveBasename), 0755)
	}

	// Process icons
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
		infoC.Printf("[INFO] Icons: %s, %s\n", iconPaths.AdaptiveXML, notifPath)
	}

	// Duplicate theme and create gradle
	themeFiles, err := flavour.DuplicateTheme(clientData.ThemeID, clientData.ArchiveBasename, outputDir, clientData, dryRun)
	if err != nil {
		return fmt.Errorf("duplicate theme: %w", err)
	}

	if verbose {
		infoC.Printf("[INFO] Gradle: %s\n", themeFiles.GradleFile)
	}

	// Update build_type.gradle and flavours.gradle
	if !dryRun {
		if err := flavour.AddFlavorToBuildType(clientData.ArchiveBasename, androidProject, dryRun); err != nil {
			warnC.Printf("[WARN] build_type.gradle: %v\n", err)
		}
		if err := flavour.AddFlavorToFlavours(clientData.ArchiveBasename, androidProject, dryRun); err != nil {
			warnC.Printf("[WARN] flavours.gradle: %v\n", err)
		}
	}

	// Generate keystore
	keystorePath, err := keystore.Generate(clientData, outputDir, dryRun)
	if err != nil {
		return fmt.Errorf("keystore: %w", err)
	}
	if verbose && keystorePath != "" {
		infoC.Printf("[INFO] Keystore: %s\n", keystorePath)
	}

	// Copy google-services.json
	if hasGS && !dryRun {
		dst := filepath.Join(outputDir, "app/src", clientData.ArchiveBasename, "google-services.json")
		if data, err := os.ReadFile(gsFile); err == nil {
			os.WriteFile(dst, data, 0644)
		}
	}

	// Success
	boldC.Print("✅ ")
	successC.Printf("Flavor created: %s\n", clientData.ArchiveBasename)
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

// loadClientData loads and computes client data from JSON file
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
