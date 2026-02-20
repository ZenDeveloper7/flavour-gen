package cmd

import (
	"runtime"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// Version info (set at build time via ldflags or hardcoded here)
var (
	Version   = "dev"     // Will be overridden by git tag in release builds
	BuildTime = "unknown" // Will be set via ldflags: -X cmd.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)
	GitCommit = "unknown" // Will be set via ldflags: -X cmd.GitCommit=$(git rev-parse --short HEAD)
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  `Display the version of flavour-gen, build time, and git commit hash.`,
	Run:   runVersion,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func runVersion(cmd *cobra.Command, args []string) {
	boldC := color.New(color.Bold)
	infoC := color.New(color.FgCyan)

	boldC.Println("Flavour Generator CLI")
	infoC.Printf("  Version:   %s\n", Version)
	infoC.Printf("  Build:     %s\n", BuildTime)
	infoC.Printf("  Commit:    %s\n", GitCommit)
	infoC.Printf("  Go:        %s\n", runtime.Version())
	infoC.Printf("  OS/Arch:   %s/%s\n", runtime.GOOS, runtime.GOARCH)
}

// GetVersion returns the version string for other uses
func GetVersion() string {
	return Version
}

// GetBuildInfo returns full build info
func GetBuildInfo() map[string]string {
	return map[string]string{
		"version":    Version,
		"buildTime":  BuildTime,
		"gitCommit":  GitCommit,
		"goVersion":  runtime.Version(),
		"os":         runtime.GOOS,
		"arch":       runtime.GOARCH,
		"buildDate":  time.Now().Format("2006-01-02"),
	}
}
