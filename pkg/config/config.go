package config

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// ClientData represents the JSON input from data.json
type ClientData struct {
	AppName           string `json:"app_name"`
	ArchiveBasename   string `json:"archivebasename"`
	PackageName       string `json:"package_name"`
	VersionName       string `json:"version_name,omitempty"`   // Default: "1.0.0"
	VersionCode       int    `json:"version_code,omitempty"`   // Default: 0
	ThemeID           int    `json:"theme_id"`
	AppLogo           string `json:"app_logo"` // Path to logo file (relative to input folder)
	BaseURL           string `json:"base_url"`
	FirebaseURL       string `json:"firebase_url,omitempty"`
	DynamicLinkDomain string `json:"dynamic_link_domain,omitempty"`
	DynamicLinkPrefix string `json:"dynamic_link_prefix,omitempty"`

	// Computed fields (generated automatically)
	Identity        string `json:"identity,omitempty"`
	DotCount        int    `json:"dot_count,omitempty"`
	AltAppName      string `json:"alt_app_name,omitempty"`
	DownloadFolder  string `json:"download_folder,omitempty"`
	TestBaseURL     string `json:"test_base_url,omitempty"`     // Computed from base_url
	EducationNumber int    `json:"education_number,omitempty"`   // From google-services.json
}

// Compute calculates derived fields from the input data
func (cd *ClientData) Compute() {
	// Default version_name
	if cd.VersionName == "" {
		cd.VersionName = "1.0.0"
	}

	// Default version_code
	if cd.VersionCode == 0 {
		cd.VersionCode = 0
	}

	// TestBaseURL defaults to BaseURL
	if cd.TestBaseURL == "" {
		cd.TestBaseURL = cd.BaseURL
	}

	// IDENTITY = base64(package_name)
	cd.Identity = base64.StdEncoding.EncodeToString([]byte(cd.PackageName))

	// DOT_COUNT = count of dots in package_name
	cd.DotCount = strings.Count(cd.PackageName, ".")

	// ALT_APP_NAME = app_name with spaces replaced by underscores, lowercase
	cd.AltAppName = strings.ToLower(strings.ReplaceAll(cd.AppName, " ", "_"))

	// DOWNLOAD_FOLDER_NAME = archivebasename
	cd.DownloadFolder = cd.ArchiveBasename
}

// UnmarshalJSON parses JSON data into ClientData
func (cd *ClientData) UnmarshalJSON(data []byte) error {
	type Alias ClientData
	var aux Alias
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	*cd = ClientData(aux)
	return nil
}

// Validate checks if required fields are present
func (cd *ClientData) Validate() error {
	if cd.AppName == "" {
		return fmt.Errorf("app_name is required")
	}
	if cd.ArchiveBasename == "" {
		return fmt.Errorf("archivebasename is required")
	}
	if cd.PackageName == "" {
		return fmt.Errorf("package_name is required")
	}
	if cd.ThemeID <= 0 {
		return fmt.Errorf("theme_id is required")
	}
	if cd.AppLogo == "" {
		return fmt.Errorf("app_logo is required")
	}
	if cd.BaseURL == "" {
		return fmt.Errorf("base_url is required")
	}
	return nil
}

// GoogleServicesJSON represents the structure of google-services.json
type GoogleServicesJSON struct {
	ProjectInfo ProjectInfo `json:"project_info"`
}

type ProjectInfo struct {
	ProjectNumber string `json:"project_number"`
}

// ExtractEducationNumber extracts education_number from google-services.json
func ExtractEducationNumber(gsJSONPath string) (int, error) {
	data, err := os.ReadFile(gsJSONPath)
	if err != nil {
		return 0, err
	}

	var gs GoogleServicesJSON
	if err := json.Unmarshal(data, &gs); err != nil {
		return 0, err
	}

	// Project number is the education number
	if gs.ProjectInfo.ProjectNumber == "" {
		return 0, fmt.Errorf("project_number not found in google-services.json")
	}

	// Convert project number to int
	var projectNum int
	fmt.Sscanf(gs.ProjectInfo.ProjectNumber, "%d", &projectNum)

	return projectNum, nil
}
