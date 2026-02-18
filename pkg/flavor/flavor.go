package flavor

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ZenDeveloper7/flavour-gen/pkg/config"
	"golang.org/x/exp/slices"
)

var (
	templatesDir = getTemplatesDir()
)

// SetTemplatesDir allows overriding the templates directory from outside the package
func SetTemplatesDir(dir string) {
	templatesDir = dir
}

func getTemplatesDir() string {
	// Allow env override
	if dir := os.Getenv("FLAVOUR_TEMPLATES"); dir != "" {
		return dir
	}
	// Default to current working directory
	cwd, _ := os.Getwd()
	return filepath.Join(cwd, "templates")
}

type ThemeFiles struct {
	ThemeSampleDir string
	GradleFile     string
}

// DuplicateTheme copies the theme sample folder and writes the Gradle flavor file.
func DuplicateTheme(themeID int, archiveBasename string, outputDir string, cd *config.ClientData, dryRun bool) (ThemeFiles, error) {
	var files ThemeFiles

	// Source paths in the Android project
	srcThemeDir := filepath.Join(templatesDir, "src", fmt.Sprintf("appx_theme%d_sample", themeID))
	srcGradle := filepath.Join(templatesDir, "flavours", fmt.Sprintf("appx_theme%d.gradle", themeID))

	// Validate source exists
	if _, err := os.Stat(srcThemeDir); os.IsNotExist(err) {
		return files, fmt.Errorf("theme sample folder not found: %s", srcThemeDir)
	}
	if _, err := os.Stat(srcGradle); os.IsNotExist(err) {
		return files, fmt.Errorf("theme gradle not found: %s", srcGradle)
	}

	// Destination paths
	dstSrcDir := filepath.Join(outputDir, "app/src", archiveBasename)
	dstGradle := filepath.Join(outputDir, "app/flavours", archiveBasename+".gradle")

	if !dryRun {
		// Create destination directories
		if err := os.MkdirAll(dstSrcDir, 0755); err != nil {
			return files, err
		}
		if err := os.MkdirAll(filepath.Dir(dstGradle), 0755); err != nil {
			return files, err
		}

		// Copy the entire theme sample folder
		if err := copyDir(srcThemeDir, dstSrcDir, cd); err != nil {
			return files, fmt.Errorf("copy theme folder: %w", err)
		}

		// Read the theme gradle
		gradleContent, err := os.ReadFile(srcGradle)
		if err != nil {
			return files, fmt.Errorf("read theme gradle: %w", err)
		}

		text := string(gradleContent)

		// Replace the flavor name (e.g., appx_theme1 -> knr_logics_clone)
		// This handles both the flavor block name and archivesBaseName
		text = strings.ReplaceAll(text, fmt.Sprintf("appx_theme%d", themeID), cd.ArchiveBasename)

		// Replace applicationId - find the line and replace the value
		text = replaceGradleLine(text, "applicationId", cd.PackageName)

		// Replace versionName
		text = replaceGradleLine(text, "versionName", cd.VersionName)

		// Replace versionCode (integer, no quotes)
		text = replaceGradleLineInt(text, "versionCode", cd.VersionCode)

		// Replace app_name in resValue
		text = replaceResValue(text, "app_name", cd.AppName)

		// Replace buildConfigField String values - these have escaped quotes
		text = replaceBuildConfigString(text, "BASE_URL", cd.BaseURL)
		text = replaceBuildConfigString(text, "TEST_BASE_URL", cd.TestBaseURL)
		text = replaceBuildConfigString(text, "FIREBASE_URL", cd.FirebaseURL)
		text = replaceBuildConfigString(text, "APP_NAME", cd.AppName)
		text = replaceBuildConfigString(text, "ALT_APP_NAME", cd.AltAppName)
		text = replaceBuildConfigString(text, "DOWNLOAD_FOLDER_NAME", cd.DownloadFolder)
		text = replaceBuildConfigString(text, "IDENTITY", cd.Identity)
		text = replaceBuildConfigString(text, "DYNAMIC_LINK_DOMAIN", cd.DynamicLinkDomain)
		text = replaceBuildConfigString(text, "DYNAMIC_LINK_PREFIX", cd.DynamicLinkPrefix)

		// Replace buildConfigField int values
		text = replaceBuildConfigInt(text, "DOT_COUNT", cd.DotCount)

		if err := os.WriteFile(dstGradle, []byte(text), 0644); err != nil {
			return files, fmt.Errorf("write gradle: %w", err)
		}
	}

	files.ThemeSampleDir = dstSrcDir
	files.GradleFile = dstGradle

	return files, nil
}

// replaceGradleLine replaces a simple key "value" line
func replaceGradleLine(text, key, newValue string) string {
	// Find line containing the key and replace the quoted value
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), key+" ") {
			// Replace the value in quotes
			old := fmt.Sprintf(`%s "%s"`, key, extractValue(text, key))
			if old == fmt.Sprintf(`%s ""`, key) {
				// Key not found with quoted value
				continue
			}
			lines[i] = strings.ReplaceAll(line, old, fmt.Sprintf(`%s "%s"`, key, newValue))
			break
		}
	}
	return strings.Join(lines, "\n")
}

// replaceGradleLineInt replaces a key value (integer, no quotes)
func replaceGradleLineInt(text, key string, newValue int) string {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), key+" ") {
			// Find and replace the integer value
			oldVal := extractIntValue(line, key)
			old := fmt.Sprintf(`%s %d`, key, oldVal)
			lines[i] = strings.ReplaceAll(line, old, fmt.Sprintf(`%s %d`, key, newValue))
			break
		}
	}
	return strings.Join(lines, "\n")
}

// replaceResValue replaces resValue "string", "key", "value"
func replaceResValue(text, key, newValue string) string {
	old := fmt.Sprintf(`resValue "string", "%s", "%s"`, key, extractResValue(text, key))
	return strings.ReplaceAll(text, old, fmt.Sprintf(`resValue "string", "%s", "%s"`, key, newValue))
}

// replaceBuildConfigString replaces buildConfigField("String", "KEY", "\"value\"")
func replaceBuildConfigString(text, key, newValue string) string {
	oldValue := extractBuildConfigString(text, key)
	if oldValue == "" {
		return text // Key not found
	}
	old := fmt.Sprintf(`buildConfigField("String", "%s", "\"%s\"")`, key, oldValue)
	return strings.ReplaceAll(text, old, fmt.Sprintf(`buildConfigField("String", "%s", "\"%s\"")`, key, newValue))
}

// replaceBuildConfigInt replaces buildConfigField("int", "KEY", value)
func replaceBuildConfigInt(text, key string, newValue int) string {
	old := fmt.Sprintf(`buildConfigField("int", "%s", %d)`, key, extractBuildConfigInt(text, key))
	if old == fmt.Sprintf(`buildConfigField("int", "%s", 0)`, key) && !containsBuildConfigInt(text, key) {
		return text // Key not found
	}
	return strings.ReplaceAll(text, old, fmt.Sprintf(`buildConfigField("int", "%s", %d)`, key, newValue))
}

// Helper functions to extract values
func extractValue(text, key string) string {
	// Find the line containing the key
	for _, line := range strings.Split(text, "\n") {
		if strings.Contains(line, key+" ") {
			// Extract value in quotes
			idx := strings.Index(line, `"`)
			if idx >= 0 {
				endIdx := strings.Index(line[idx+1:], `"`)
				if endIdx >= 0 {
					return line[idx+1 : idx+1+endIdx]
				}
			}
		}
	}
	return ""
}

func extractIntValue(text, key string) int {
	// First try to find in lines
	for _, line := range strings.Split(text, "\n") {
		if strings.Contains(line, key+" ") {
			var val int
			_, err := fmt.Sscanf(strings.TrimSpace(line), key+" %d", &val)
			if err == nil {
				return val
			}
		}
	}
	return 0
}

func extractResValue(text, key string) string {
	pattern := fmt.Sprintf(`resValue "string", "%s", "`, key)
	idx := strings.Index(text, pattern)
	if idx >= 0 {
		start := idx + len(pattern)
		endIdx := strings.Index(text[start:], `"`)
		if endIdx >= 0 {
			return text[start : start+endIdx]
		}
	}
	return ""
}

func extractBuildConfigString(text, key string) string {
	pattern := fmt.Sprintf(`buildConfigField("String", "%s", "\"`, key)
	idx := strings.Index(text, pattern)
	if idx >= 0 {
		start := idx + len(pattern)
		endIdx := strings.Index(text[start:], `\"`)
		if endIdx >= 0 {
			return text[start : start+endIdx]
		}
	}
	return ""
}

func extractBuildConfigInt(text, key string) int {
	pattern := fmt.Sprintf(`buildConfigField("int", "%s", %%d)`, key)
	for _, line := range strings.Split(text, "\n") {
		var val int
		if _, err := fmt.Sscanf(line, pattern, &val); err == nil {
			return val
		}
	}
	return 0
}

func containsBuildConfigInt(text, key string) bool {
	pattern := fmt.Sprintf(`buildConfigField("int", "%s"`, key)
	return strings.Contains(text, pattern)
}

func replacePlaceholders(text string, cd *config.ClientData) string {
	// Replace placeholder format ${VAR}
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

func copyDir(src, dst string, cd *config.ClientData) error {
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
			if err := copyDir(srcPath, dstPath, cd); err != nil {
				return err
			}
		} else {
			in, err := os.ReadFile(srcPath)
			if err != nil {
				return err
			}
			text := string(in)
			ext := strings.ToLower(filepath.Ext(entry.Name()))
			if slices.Contains([]string{".gradle", ".xml", ".json", ".properties", ".kt", ".java", ".txt", ".md"}, ext) {
				text = replacePlaceholders(text, cd)
			}
			if err := os.WriteFile(dstPath, []byte(text), 0644); err != nil {
				return err
			}
		}
	}
	return nil
}

// ListThemes scans the project's flavours folder for theme IDs.
func ListThemes() ([]int, error) {
	var themes []int

	// Check in flavours folder for gradle files
	flavoursDir := filepath.Join(templatesDir, "flavours")
	entries, err := os.ReadDir(flavoursDir)
	if err != nil {
		return themes, err
	}
	for _, entry := range entries {
		name := entry.Name()
		if strings.HasSuffix(name, ".gradle") && strings.HasPrefix(name, "appx_theme") {
			name := strings.TrimSuffix(name, ".gradle")
			var id int
			if n, err := fmt.Sscanf(name, "appx_theme%d", &id); err == nil && n == 1 {
				themes = append(themes, id)
			}
		}
	}
	return themes, nil
}
