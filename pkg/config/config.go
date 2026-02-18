package config

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

// ClientData represents the JSON input from data.json
type ClientData struct {
	AppName           string `json:"app_name"`
	ArchiveBasename   string `json:"archivebasename"`
	PackageName       string `json:"package_name"`
	VersionName       string `json:"version_name"`
	VersionCode       int    `json:"version_code"`
	ThemeID           int    `json:"theme_id"`
	AppLogo           string `json:"app_logo"`         // Path to logo file (relative to input folder)
	EducationNumber   int    `json:"education_number"` // For gradle comment
	BaseURL           string `json:"base_url"`
	TestBaseURL       string `json:"test_base_url"`
	FirebaseURL       string `json:"firebase_url"`
	DynamicLinkDomain string `json:"dynamic_link_domain"`
	DynamicLinkPrefix string `json:"dynamic_link_prefix"`

	// Computed fields (generated automatically)
	Identity       string `json:"identity,omitempty"`
	DotCount       int    `json:"dot_count,omitempty"`
	AltAppName     string `json:"alt_app_name,omitempty"`
	DownloadFolder string `json:"download_folder,omitempty"`
}

// Compute calculates derived fields from the input data
func (cd *ClientData) Compute() {
	// IDENTITY = base64(package_name)
	cd.Identity = base64.StdEncoding.EncodeToString([]byte(cd.PackageName))

	// DOT_COUNT = count of dots in package_name
	cd.DotCount = strings.Count(cd.PackageName, ".")

	// ALT_APP_NAME = app_name with spaces replaced by underscores, lowercase
	cd.AltAppName = strings.ToLower(strings.ReplaceAll(cd.AppName, " ", "_"))

	// DOWNLOAD_FOLDER_NAME = archivebasename with underscores
	cd.DownloadFolder = cd.ArchiveBasename
}

// UnmarshalJSON parses JSON data into ClientData
func (cd *ClientData) UnmarshalJSON(data []byte) error {
	// Use standard json unmarshaling
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
	return nil
}
