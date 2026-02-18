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

// Global templates directory (can be overridden)
var templatesDir = getTemplatesDir()

// SetTemplatesDir allows overriding the templates directory from outside the package
func SetTemplatesDir(dir string) {
	templatesDir = dir
}

func getTemplatesDir() string {
	if dir := os.Getenv("FLAVOUR_TEMPLATES"); dir != "" {
		return dir
	}
	cwd, _ := os.Getwd()
	return filepath.Join(cwd, "templates")
}

// ThemeFiles holds paths to generated theme files
type ThemeFiles struct {
	ThemeSampleDir string
	GradleFile     string
}

// DuplicateTheme copies the theme sample folder and creates the Gradle flavor file
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
		if err := os.MkdirAll(dstSrcDir, 0755); err != nil {
			return files, err
		}
		if err := os.MkdirAll(filepath.Dir(dstGradle), 0755); err != nil {
			return files, err
		}

		// Copy theme folder with placeholder replacement
		if err := copyDir(srcThemeDir, dstSrcDir, cd); err != nil {
			return files, fmt.Errorf("copy theme folder: %w", err)
		}

		// Read and modify gradle file
		gradleContent, err := os.ReadFile(srcGradle)
		if err != nil {
			return files, fmt.Errorf("read theme gradle: %w", err)
		}

		text := processGradle(string(gradleContent), themeID, cd)

		// Update strings.xml
		if err := updateStringsXML(dstSrcDir, cd); err != nil {
			return files, fmt.Errorf("update strings.xml: %w", err)
		}

		if err := os.WriteFile(dstGradle, []byte(text), 0644); err != nil {
			return files, fmt.Errorf("write gradle: %w", err)
		}
	}

	files.ThemeSampleDir = dstSrcDir
	files.GradleFile = dstGradle
	return files, nil
}

// processGradle replaces all placeholder values in the gradle file
func processGradle(text string, themeID int, cd *config.ClientData) string {
	// Education comment
	if cd.EducationNumber > 0 {
		text = regexp.MustCompile(`//Education \d+`).
			ReplaceAllString(text, fmt.Sprintf("//Education %d", cd.EducationNumber))
	}

	// Replace flavor name (with _km suffix)
	oldFlavorName := fmt.Sprintf("appx_theme%d_km", themeID)
	text = strings.ReplaceAll(text, oldFlavorName, cd.ArchiveBasename)
	// Also replace without suffix
	oldFlavorName = fmt.Sprintf("appx_theme%d", themeID)
	text = strings.ReplaceAll(text, oldFlavorName, cd.ArchiveBasename)

	// Replace simple key-value pairs
	text = replaceGradleValue(text, "applicationId", cd.PackageName)
	text = replaceGradleValue(text, "versionName", cd.VersionName)
	text = replaceGradleValueInt(text, "versionCode", cd.VersionCode)
	text = replaceResValue(text, "app_name", cd.AppName)

	// Replace buildConfigField String values
	text = replaceBuildConfig(text, "String", "BASE_URL", cd.BaseURL)
	text = replaceBuildConfig(text, "String", "TEST_BASE_URL", cd.TestBaseURL)
	text = replaceBuildConfig(text, "String", "FIREBASE_URL", cd.FirebaseURL)
	text = replaceBuildConfig(text, "String", "APP_NAME", cd.AppName)
	text = replaceBuildConfig(text, "String", "ALT_APP_NAME", cd.AltAppName)
	text = replaceBuildConfig(text, "String", "DOWNLOAD_FOLDER_NAME", cd.DownloadFolder)
	text = replaceBuildConfig(text, "String", "IDENTITY", cd.Identity)
	text = replaceBuildConfig(text, "String", "DYNAMIC_LINK_DOMAIN", cd.DynamicLinkDomain)
	text = replaceBuildConfig(text, "String", "DYNAMIC_LINK_PREFIX", cd.DynamicLinkPrefix)

	// Replace buildConfigField int values
	text = replaceBuildConfigInt(text, "DOT_COUNT", cd.DotCount)

	// Add DOT_COUNT if missing
	if cd.DotCount > 0 && !strings.Contains(text, "DOT_COUNT") {
		text = addBuildConfigInt(text, "DOT_COUNT", cd.DotCount)
	}

	return text
}

// replaceGradleValue replaces a simple key "value" line
func replaceGradleValue(text, key, newValue string) string {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), key+" ") {
			old := extractQuotedValue(line, key)
			if old == "" {
				continue
			}
			lines[i] = strings.ReplaceAll(line, fmt.Sprintf(`%s "%s"`, key, old), fmt.Sprintf(`%s "%s"`, key, newValue))
			break
		}
	}
	return strings.Join(lines, "\n")
}

// replaceGradleValueInt replaces a key value (integer, no quotes)
func replaceGradleValueInt(text, key string, newValue int) string {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), key+" ") {
			oldVal := extractIntValue(line, key)
			if oldVal == 0 {
				continue
			}
			lines[i] = strings.ReplaceAll(line, fmt.Sprintf(`%s %d`, key, oldVal), fmt.Sprintf(`%s %d`, key, newValue))
			break
		}
	}
	return strings.Join(lines, "\n")
}

// replaceResValue replaces resValue "string", "key", "value"
func replaceResValue(text, key, newValue string) string {
	old := extractResValue(text, key)
	if old == "" {
		return text
	}
	return strings.ReplaceAll(text, fmt.Sprintf(`resValue "string", "%s", "%s"`, key, old),
		fmt.Sprintf(`resValue "string", "%s", "%s"`, key, newValue))
}

// replaceBuildConfig replaces buildConfigField("TYPE", "KEY", "\"value\"")
func replaceBuildConfig(text, fieldType, key, newValue string) string {
	old := extractBuildConfig(text, fieldType, key)
	if old == "" {
		return text
	}
	oldPattern := fmt.Sprintf(`buildConfigField("%s", "%s", "\"%s\"")`, fieldType, key, old)
	newPattern := fmt.Sprintf(`buildConfigField("%s", "%s", "\"%s\")`, fieldType, key, newValue)
	return strings.ReplaceAll(text, oldPattern, newPattern)
}

// replaceBuildConfigInt replaces buildConfigField("int", "KEY", value)
func replaceBuildConfigInt(text, key string, newValue int) string {
	old := extractBuildConfigInt(text, key)
	if old == 0 {
		return text
	}
	oldPattern := fmt.Sprintf(`buildConfigField("int", "%s", %d)`, key, old)
	newPattern := fmt.Sprintf(`buildConfigField("int", "%s", %d)`, key, newValue)
	return strings.ReplaceAll(text, oldPattern, newPattern)
}

// addBuildConfigInt adds a new buildConfigField int if it doesn't exist
func addBuildConfigInt(text, key string, value int) string {
	newEntry := fmt.Sprintf(`            buildConfigField("int", "%s", %d)`, key, value)
	lines := strings.Split(text, "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		if strings.Contains(lines[i], "buildConfigField") {
			var newLines []string
			newLines = append(newLines, lines[:i+1]...)
			newLines = append(newLines, newEntry)
			newLines = append(newLines, lines[i+1:]...)
			return strings.Join(newLines, "\n")
		}
	}
	return text + "\n" + newEntry
}

// Helper functions to extract values from gradle lines
func extractQuotedValue(line, key string) string {
	pattern := fmt.Sprintf(`%s "`, key)
	idx := strings.Index(line, pattern)
	if idx >= 0 {
		start := idx + len(pattern)
		endIdx := strings.Index(line[start:], `"`)
		if endIdx >= 0 {
			return line[start : start+endIdx]
		}
	}
	return ""
}

func extractIntValue(line, key string) int {
	var val int
	_, err := fmt.Sscanf(strings.TrimSpace(line), key+" %d", &val)
	if err == nil {
		return val
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

func extractBuildConfig(text, fieldType, key string) string {
	pattern := fmt.Sprintf(`buildConfigField("%s", "%s", "\"`, fieldType, key)
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

// copyDir copies a directory recursively, replacing placeholders in text files
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

// replacePlaceholders replaces ${VAR} style placeholders in text
func replacePlaceholders(text string, cd *config.ClientData) string {
	replacements := map[string]string{
		"${APP_NAME}":            cd.AppName,
		"${ARCHIVE_BASENAME}":    cd.ArchiveBasename,
		"${PACKAGE_NAME}":         cd.PackageName,
		"${VERSION_NAME}":        cd.VersionName,
		"${VERSION_CODE}":         fmt.Sprintf("%d", cd.VersionCode),
		"${BASE_URL}":             cd.BaseURL,
		"${TEST_BASE_URL}":        cd.TestBaseURL,
		"${FIREBASE_URL}":         cd.FirebaseURL,
		"${DYNAMIC_LINK_DOMAIN}":  cd.DynamicLinkDomain,
		"${DYNAMIC_LINK_PREFIX}":  cd.DynamicLinkPrefix,
		"${IDENTITY}":             cd.Identity,
		"${DOT_COUNT}":            fmt.Sprintf("%d", cd.DotCount),
		"${ALT_APP_NAME}":         cd.AltAppName,
		"${DOWNLOAD_FOLDER_NAME}": cd.DownloadFolder,
	}
	for placeholder, value := range replacements {
		text = strings.ReplaceAll(text, placeholder, value)
	}
	return text
}

// ListThemes scans the project's flavours folder for theme IDs
func ListThemes() ([]int, error) {
	var themes []int
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

// AddFlavorToBuildType adds the new flavor signing config to build_type.gradle
func AddFlavorToBuildType(archiveBasename string, projectRoot string, dryRun bool) error {
	buildTypePath := filepath.Join(projectRoot, "app", "build_type.gradle")
	if _, err := os.Stat(buildTypePath); err != nil {
		return fmt.Errorf("build_type.gradle not found: %w", err)
	}

	content, err := os.ReadFile(buildTypePath)
	if err != nil {
		return err
	}

	// Check if already exists
	if strings.Contains(string(content), fmt.Sprintf("productFlavors.%s.signingConfig", archiveBasename)) {
		return nil
	}

	text := string(content)
	newEntry := fmt.Sprintf(`            productFlavors.%s.signingConfig signingConfigs.%s`, archiveBasename, archiveBasename)

	lines := strings.Split(text, "\n")
	releaseStarted := false
	inserted := false

	for i, line := range lines {
		if strings.Contains(line, "release {") {
			releaseStarted = true
			continue
		}
		if releaseStarted && strings.TrimSpace(line) == "}" {
			var newLines []string
			newLines = append(newLines, lines[:i]...)
			newLines = append(newLines, newEntry)
			newLines = append(newLines, lines[i:]...)
			text = strings.Join(newLines, "\n")
			inserted = true
			break
		}
	}

	if !inserted {
		text = strings.TrimRight(text, "\n") + "\n" + newEntry + "\n"
	}

	if !dryRun {
		return os.WriteFile(buildTypePath, []byte(text), 0644)
	}
	return nil
}

// AddFlavorToFlavours adds the new flavor to flavours.gradle
func AddFlavorToFlavours(archiveBasename string, projectRoot string, dryRun bool) error {
	flavoursPath := filepath.Join(projectRoot, "app", "flavours.gradle")
	if _, err := os.Stat(flavoursPath); err != nil {
		return fmt.Errorf("flavours.gradle not found: %w", err)
	}

	content, err := os.ReadFile(flavoursPath)
	if err != nil {
		return err
	}

	if strings.Contains(string(content), fmt.Sprintf("flavours/%s.gradle", archiveBasename)) {
		return nil
	}

	newEntry := fmt.Sprintf(`apply from:  './flavours/%s.gradle'`, archiveBasename)
	text := strings.TrimRight(string(content), "\n") + "\n" + newEntry + "\n"

	if !dryRun {
		return os.WriteFile(flavoursPath, []byte(text), 0644)
	}
	return nil
}

// updateStringsXML updates strings.xml with client data
func updateStringsXML(themeDir string, cd *config.ClientData) error {
	stringsPath := filepath.Join(themeDir, "res/values/strings.xml")
	if _, err := os.Stat(stringsPath); err != nil {
		return nil // Skip if doesn't exist
	}

	content, err := os.ReadFile(stringsPath)
	if err != nil {
		return err
	}

	text := string(content)

	// Update dynamic_link_host from dynamic_link_prefix
	if cd.DynamicLinkPrefix != "" {
		host := strings.TrimPrefix(cd.DynamicLinkPrefix, "https://")
		host = strings.TrimSuffix(host, "/")
		text = strings.Replace(text, ">appxcore.com<", ">"+host+"<", 1)
	}

	// Update app_name
	if cd.AppName != "" {
		text = strings.Replace(text, ">appx_theme1<", ">"+cd.AppName+"<", 1)
	}

	return os.WriteFile(stringsPath, []byte(text), 0644)
}
