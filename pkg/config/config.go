package config

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"strings"
)

type ClientData struct {
	AppName           string `json:"app_name"`
	ArchiveBasename   string `json:"archivebasename"`
	PackageName       string `json:"package_name"`
	VersionName       string `json:"version_name"`
	VersionCode       int    `json:"version_code"`
	BaseURL           string `json:"base_url"`
	TestBaseURL       string `json:"test_base_url"`
	FirebaseURL       string `json:"firebase_url"`
	DynamicLinkDomain string `json:"dynamic_link_domain"`
	DynamicLinkPrefix string `json:"dynamic_link_prefix"`

	// Computed fields (internal)
	Identity          string `json:"identity,omitempty"`
	DotCount          int    `json:"dot_count,omitempty"`
	AltAppName        string `json:"alt_app_name,omitempty"`
	DownloadFolder    string `json:"download_folder,omitempty"`
}

// Compute fills derived fields from package name and app name
func (cd *ClientData) Compute() {
	// IDENTITY = base64(package_name)
	Identity = base64.StdEncoding.EncodeToString([]byte(cd.PackageName))

	// DOT_COUNT = number of '.' in package name (excluding last segment)
	parts := strings.Split(cd.PackageName, ".")
	cd.DotCount = len(parts) - 1

	// ALT_APP_NAME: first char uppercase, rest lowercase
	if cd.AppName != "" {
		runes := []rune(strings.ToLower(cd.AppName))
		if len(runes) > 0 {
			runes[0] = []rune(strings.ToUpper(string(runes[0])))[0]
		}
		cd.AltAppName = string(runes)
	}

	// DOWNLOAD_FOLDER_NAME: first 64 chars of archivebasename, sanitized
	cd.DownloadFolder = sanitizeFolderName(cd.ArchiveBasename)
}

func sanitizeFolderName(s string) string {
	// Lowercase, replace spaces/unsafe chars with underscore, trim to 64
	s = strings.ToLower(s)
	s = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			return r
		}
		return '_'
	}, s)
	if len(s) > 64 {
		s = s[:64]
	}
	return s
}