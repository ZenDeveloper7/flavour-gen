package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	flavour "github.com/ZenDeveloper7/flavour-gen/pkg/flavour"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check system and project status",
	Long: `Display system information, check prerequisites, detect Android projects,
and list available themes. Useful for debugging setup issues.`,
	RunE: runStatus,
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func runStatus(cmd *cobra.Command, args []string) error {
	boldC := color.New(color.Bold)
	successC := color.New(color.FgGreen)
	errorC := color.New(color.FgRed)
	warnC := color.New(color.FgYellow)
	infoC := color.New(color.FgCyan)

	// Header
	boldC.Println("\n🔍 Flavour Generator Status")
	boldC.Println("============================")

	// Version Section
	boldC.Println("\n📦 Version Information")
	infoC.Printf("  Version:   %s\n", GetVersion())
	infoC.Printf("  Go:        %s\n", runtime.Version())
	infoC.Printf("  OS/Arch:   %s/%s\n", runtime.GOOS, runtime.GOARCH)

	// Prerequisites Section
	boldC.Println("\n✅ Prerequisites")
	
	// Check keytool
	if _, err := exec.LookPath("keytool"); err != nil {
		errorC.Println("  ✗ keytool    - Install Java JDK")
	} else {
		successC.Println("  ✓ keytool    - Java Keystore tool found")
	}

	// Check gradle
	if _, err := exec.LookPath("gradle"); err != nil {
		warnC.Println("  ⚠ gradle     - Not found (optional, for advanced builds)")
	} else {
		successC.Println("  ✓ gradle     - Gradle build tool found")
	}

	// Check working directory
	boldC.Println("\n📁 Working Directory")
	cwd, _ := os.Getwd()
	infoC.Printf("  Current:   %s\n", cwd)

	// Try to detect Android project
	boldC.Println("\n🔎 Android Project Detection")
	appPath := filepath.Join(cwd, "app")
	if _, err := os.Stat(appPath); err == nil {
		successC.Printf("  ✓ Found app/ directory at: %s\n", appPath)
		
		// Check for key gradle files
		checkFile := func(name string) {
			path := filepath.Join(appPath, name)
			if _, err := os.Stat(path); err == nil {
				successC.Printf("  ✓ Found %s\n", name)
			} else {
				warnC.Printf("  ⚠ Missing %s\n", name)
			}
		}
		
		checkFile("build.gradle")
		checkFile("flavours.gradle")
		checkFile("build_type.gradle")
		checkFile("keystore.gradle")
	} else {
		warnC.Println("  ⚠ No app/ directory found in current path")
		infoC.Println("  Tip: Run from Android project root or use --output-dir")
	}

	// Available Themes
	boldC.Println("\n🎨 Available Themes")
	themes, err := flavour.ListThemes()
	if err != nil {
		errorC.Printf("  ✗ Error: %v\n", err)
	} else if len(themes) == 0 {
		warnC.Println("  ⚠ No themes found")
		infoC.Println("  Tip: Ensure flavours/ directory exists with appx_theme*.gradle files")
	} else {
		successC.Printf("  ✓ Found %d theme(s):\n", len(themes))
		for _, themeID := range themes {
			infoC.Printf("    • Theme %d\n", themeID)
		}
	}

	// Configuration Section
	boldC.Println("\n⚙️ Configuration")
	if dir := os.Getenv("FLAVOUR_TEMPLATES"); dir != "" {
		infoC.Printf("  FLAVOUR_TEMPLATES: %s\n", dir)
	} else {
		infoC.Println("  FLAVOUR_TEMPLATES: not set (using default)")
	}

	// Summary
	boldC.Println("\n📊 Summary")
	isReady := true
	if _, err := exec.LookPath("keytool"); err != nil {
		isReady = false
		errorC.Println("  ✗ Not ready: keytool is required")
	}
	if len(themes) == 0 {
		isReady = false
		errorC.Println("  ✗ Not ready: No themes available")
	}
	
	if isReady {
		successC.Println("  ✓ Ready to generate flavours!")
	} else {
		errorC.Println("  ✗ Setup incomplete - see issues above")
	}

	fmt.Println()
	return nil
}
