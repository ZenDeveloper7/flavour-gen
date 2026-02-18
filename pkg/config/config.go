package config

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
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

// ClientDataArray represents multiple clients (for batch processing)
type ClientDataArray []ClientData

// UnmarshalJSON handles both single client and array of clients
func (cda *ClientDataArray) UnmarshalJSON(data []byte) error {
	// Try to parse as array first
	var clients []ClientData
	if err := json.Unmarshal(data, &clients); err == nil {
		*cda = clients
		return nil
	}

	// If array fails, try single client
	var client ClientData
	if err := json.Unmarshal(data, &client); err != nil {
		return err
	}
	*cda = []ClientData{client}
	return nil
}

// LoadClientData loads client data from JSON file (supports single or array)
func LoadClientData(path string) ([]ClientData, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var clients ClientDataArray
	if err := json.Unmarshal(data, &clients); err != nil {
		return nil, err
	}

	// Compute derived fields for each client
	for i := range clients {
		clients[i].Compute()
	}

	return clients, nil
}

// Compute calculates derived fields from the input data
func (cd *ClientData) Compute() {
	// Default version_name
	if cd.VersionName == "" {
		cd.VersionName = "1.0.0"
	}

	// Default version_code
	if cd.VersionCode == 0 {
		cd.VersionCode = 1
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
	Client      []Client    `json:"client"`
}

type ProjectInfo struct {
	ProjectNumber string `json:"project_number"`
	ProjectID    string `json:"project_id"`
}

type Client struct {
	OAuthInfo []OAuthInfo `json:"oauth_client"`
}

type OAuthInfo struct {
	ClientID string `json:"client_id"`
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

	// Try to get from project_id first (e.g., "education-303-xxx" -> 303)
	if gs.ProjectInfo.ProjectID != "" {
		// Extract number from project_id (e.g., "education-303-default-rtdb" -> 303)
		parts := strings.Split(gs.ProjectInfo.ProjectID, "-")
		for _, part := range parts {
			if len(part) == 3 {
				if num, err := strconv.Atoi(part); err == nil && num > 0 {
					return num, nil
				}
			}
		}
	}

	// Fallback: use last 3 digits of project_number
	if gs.ProjectInfo.ProjectNumber != "" {
		num, _ := strconv.ParseInt(gs.ProjectInfo.ProjectNumber, 10, 64)
		return int(num % 1000), nil
	}

	return 0, fmt.Errorf("project_id or project_number not found in google-services.json")
}
