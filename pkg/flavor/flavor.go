package flavor

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/atlas/flavor-gen/pkg/config"
)

var (
	templatesDir = "templates"
)

// DuplicateTheme copies the theme sample folder and writes the Gradle flavor file.
// Returns the paths to created files.
func DuplicateTheme(themeID int, archiveBasename string, outputDir string, cd *config.ClientData, dryRun bool) (ThemeFiles, error) {
	var files ThemeFiles

	themeName := fmt.Sprintf("appx_theme%d_sample", themeID)
	srcThemeDir := filepath.Join(templatesDir, themeName)
	if _, err := os.Stat(srcThemeDir); os.IsNotExist(err) {
		return files, fmt.Errorf("theme sample folder not found: %s", srcThemeDir)
	}

	// Destination: app/src/<archiveBasename>/
	dstSrcDir := filepath.Join(outputDir, "app/src", archiveBasename)
	if !dryRun {
		if err := os.MkdirAll(dstSrcDir, 0755); err != nil {
			return files, err
		}
		// Copy entire folder contents
		if err := copyDir(srcThemeDir, dstSrcDir); err != nil {
			return files, err
		}
	}
	files.ThemeSampleDir = dstSrcDir

	// Copy the theme gradle file to app/flavours/
	srcGradle := filepath.Join(templatesDir, fmt.Sprintf("appx_theme%d.gradle", themeID))
	dstGradleDir := filepath.Join(outputDir, "app/flavours")
	if !dryRun {
		if err := os.MkdirAll(dstGradleDir, 0755); err != nil {
			return files, err
		}
		content, err := os.ReadFile(srcGradle)
		if err != nil {
			return files, fmt.Errorf("read theme gradle: %w", err)
		}
		// Replace placeholders
		text := replacePlaceholders(string(content), cd)
		dstGradle := filepath.Join(dstGradleDir, archiveBasename+".gradle")
		if err := os.WriteFile(dstGradle, []byte(text), 0644); err != nil {
			return files, err
		}
		files.GradleFile = dstGradle
	} else {
		files.GradleFile = filepath.Join(dstGradleDir, archiveBasename+".gradle")
	}

	return files, nil
}

func replacePlaceholders(text string, cd *config.ClientData) string {
	// Common placeholders
	text = strings.ReplaceAll(text, "${APP_NAME}", cd.AppName)
	text = strings.ReplaceAll(text, "${ARCHIVE_BASENAME}", cd.ArchiveBasename)
	text = strings.ReplaceAll(text, "${PACKAGE_NAME}", cd.PackageName)
	text = strings.ReplaceAll(text, "${VERSION_NAME}", cd.VersionName)
	text = strings.ReplaceAll(text, "${VERSION_CODE}", fmt.Sprintf("%d", cd.VersionCode))
	text = strings.ReplaceAll(text, "${BASE_URL}", cd.BaseURL)
	text = strings.ReplaceAll(text, "${TEST_BASE_URL}", cd.TestBaseURL)
	text = strings.ReplaceAll(text, "${FIREBASE_URL}", cd.FirebaseURL)
	text = strings.ReplaceAll(text, "${DYNAMIC_LINK_DOMAIN}", cd.DynamicLinkDomain)
	text = strings.ReplaceAll(text, "${DYNAMIC_LINK_PREFIX}", cd.DynamicLinkPrefix)
	text = strings.ReplaceAll(text, "${IDENTITY}", cd.Identity)
	text = strings.ReplaceAll(text, "${DOT_COUNT}", fmt.Sprintf("%d", cd.DotCount))
	text = strings.ReplaceAll(text, "${ALT_APP_NAME}", cd.AltAppName)
	text = strings.ReplaceAll(text, "${DOWNLOAD_FOLDER_NAME}", cd.DownloadFolder)

	return text
}

func copyDir(src, dst string) error {
	if _, err := os.Stat(dst); os.IsNotExist(err) {
		if err := os.MkdirAll(dst, 0755); err != nil {
			return err
		}
	}
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())
		if entry.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			// Copy file
			in, err := os.ReadFile(srcPath)
			if err != nil {
				return err
			}
			// Also replace placeholders in files that are likely templates
			text := string(in)
			// Heuristic: only replace in files with certain extensions
			ext := strings.ToLower(filepath.Ext(entry.Name()))
			if slices.Contains([]string{".gradle", ".xml", ".json", ".properties", ".kt", ".java", ".txt", ".md"}, ext) {
				// Try to load as JSON for google-services.json? For now keep as text replace
				text = replacePlaceholders(text, &config.ClientData{})
			}
			if err := os.WriteFile(dstPath, []byte(text), 0644); err != nil {
				return err
			}
		}
	}
	return nil
}

// ListThemes scans templates dir for theme IDs (based on appx_theme<N>_sample folders)
func ListThemes() ([]int, error) {
	var themes []int
	entries, err := os.ReadDir(templatesDir)
	if err != nil {
		return themes, err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			name := entry.Name()
			var id int
			if n, err := fmt.Sscanf(name, "appx_theme%d_sample", &id); err == nil && n == 1 {
				themes = append(themes, id)
			}
		}
	}
	return themes, nil
}

type ThemeFiles struct {
	ThemeSampleDir string
	GradleFile     string
}