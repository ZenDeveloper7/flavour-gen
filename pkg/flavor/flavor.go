package flavor

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
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

		// Read the theme gradle and modify it
		gradleContent, err := os.ReadFile(srcGradle)
		if err != nil {
			return files, fmt.Errorf("read theme gradle: %w", err)
		}
		
		// Replace values in gradle
		text := string(gradleContent)
		
		// Replace the flavor name (e.g., appx_theme1 -> knr_logics_clone)
		// This handles both the flavor block name and archivesBaseName
		text = strings.ReplaceAll(text, fmt.Sprintf("appx_theme%d", themeID), cd.ArchiveBasename)
		
		// Replace applicationId
		text = replaceGradleValue(text, "applicationId", cd.PackageName)
		
		// Replace versionName
		text = replaceGradleValue(text, "versionName", fmt.Sprintf("\"%s\"", cd.VersionName))
		
		// Replace versionCode
		text = replaceGradleValue(text, "versionCode", fmt.Sprintf("%d", cd.VersionCode))
		
		// Replace app_name in resValue
		text = replaceGradleResValue(text, "app_name", cd.AppName)
		
		// Replace URLs and other key values
		text = replaceGradleBuildConfig(text, "BASE_URL", cd.BaseURL)
		text = replaceGradleBuildConfig(text, "TEST_BASE_URL", cd.TestBaseURL)
		text = replaceGradleBuildConfig(text, "FIREBASE_URL", cd.FirebaseURL)
		text = replaceGradleBuildConfig(text, "APP_NAME", cd.AppName)
		text = replaceGradleBuildConfig(text, "ALT_APP_NAME", cd.AltAppName)
		text = replaceGradleBuildConfig(text, "DOWNLOAD_FOLDER_NAME", cd.DownloadFolder)
		text = replaceGradleBuildConfig(text, "IDENTITY", cd.Identity)
		text = replaceGradleBuildConfig(text, "DYNAMIC_LINK_DOMAIN", cd.DynamicLinkDomain)
		text = replaceGradleBuildConfig(text, "DYNAMIC_LINK_PREFIX", cd.DynamicLinkPrefix)
		text = replaceGradleBuildConfigInt(text, "DOT_COUNT", cd.DotCount)
		
		if err := os.WriteFile(dstGradle, []byte(text), 0644); err != nil {
			return files, fmt.Errorf("write gradle: %w", err)
		}
	}

	files.ThemeSampleDir = dstSrcDir
	files.GradleFile = dstGradle

	return files, nil
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
	
	// Also handle non-placeholder format (hardcoded theme values)
	// Replace applicationId, versionName, versionCode, app_name from theme1
	// These need to be replaced with client data values
	// This is done via regexp for more complex patterns in the gradle
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

// replaceGradleValue replaces a simple gradle property value
// e.g., replaceGradleValue(text, "applicationId", "com.new.package")
func replaceGradleValue(text, key, newValue string) string {
	// Match: key "value" 
	pattern := fmt.Sprintf(`(%s\s+")[^"]*`, key)
	re := regexp.MustCompile(pattern)
	return re.ReplaceAllString(text, fmt.Sprintf(`$1%s"`, newValue))
}

// replaceGradleResValue replaces resValue "string", "key", "value"
func replaceGradleResValue(text, key, newValue string) string {
	// Match: resValue "string", "key", "value"
	// Captures everything up to the value
	re := regexp.MustCompile(fmt.Sprintf(`(resValue\s+"string",\s+"%s",\s+")[^"]*`, key))
	return re.ReplaceAllString(text, fmt.Sprintf(`$1%s"`, newValue))
}

// replaceGradleBuildConfig replaces buildConfigField("TYPE", "KEY", "value")
func replaceGradleBuildConfig(text, key, newValue string) string {
	// Match: buildConfigField("String", "KEY", "value")
	pattern := fmt.Sprintf(`(buildConfigField\("String",\s+"%s",\s+")[^"]*`, key)
	re := regexp.MustCompile(pattern)
	return re.ReplaceAllString(text, fmt.Sprintf(`$1%s"`, newValue))
}

// replaceGradleBuildConfigInt replaces buildConfigField("int", "KEY", value)
func replaceGradleBuildConfigInt(text, key string, newValue int) string {
	pattern := fmt.Sprintf(`(buildConfigField\("int",\s+"%s",\s+)\d+`, key)
	re := regexp.MustCompile(pattern)
	return re.ReplaceAllString(text, fmt.Sprintf(`$1%d`, newValue))
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