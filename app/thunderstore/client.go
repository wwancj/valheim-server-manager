package thunderstore

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	baseAPIURL = "https://thunderstore.io/api/experimental/package/"
)

// Package represents a Thunderstore mod package.
type Package struct {
	Name        string   `json:"name"`
	FullName    string   `json:"full_name"`
	Owner       string   `json:"owner"`
	Description string   `json:"description"`
	Version     string   `json:"version_number"`
	Downloads   int      `json:"total_downloads"`
	Rating      float64  `json:"rating_score"`
	WebsiteURL  string   `json:"package_url"`
	Icon        string   `json:"icon"`
	Versions    []Version `json:"-"`
	Latest      *Latest   `json:"latest"`
}

// Latest contains the latest version info.
type Latest struct {
	Namespace    string   `json:"namespace"`
	Name         string   `json:"name"`
	VersionNumber string `json:"version_number"`
	FullName     string   `json:"full_name"`
	Description  string   `json:"description"`
	Icon         string   `json:"icon"`
	Dependencies []string `json:"dependencies"`
	DownloadURL  string   `json:"download_url"`
	Downloads    int      `json:"downloads"`
	DateCreated  string   `json:"date_created"`
}

// Version represents a specific version of a mod.
type Version struct {
	VersionNumber string   `json:"version_number"`
	DownloadURL   string   `json:"download_url"`
	Dependencies  []string `json:"dependencies"`
	FileSize      int      `json:"file_size"`
	DateCreated   string   `json:"date_created"`
}

// APIResponse represents the Thunderstore API response.
type APIResponse struct {
	Next     string    `json:"next"`
	Previous string    `json:"previous"`
	Results  []Package `json:"results"`
}

// Client is a Thunderstore API client.
type Client struct {
	httpClient *http.Client
}

// NewClient creates a new Thunderstore client.
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{},
	}
}

// SearchPackages searches for mods by query.
func (c *Client) SearchPackages(query string, page int) ([]Package, int, error) {
	url := fmt.Sprintf("%s?community=valheim&search=%s&page=%d", baseAPIURL, query, page)
	return c.fetchPackages(url)
}

// GetPopularPackages gets popular mods.
func (c *Client) GetPopularPackages(page int) ([]Package, int, error) {
	url := fmt.Sprintf("%s?community=valheim&page=%d", baseAPIURL, page)
	return c.fetchPackages(url)
}

func (c *Client) fetchPackages(url string) ([]Package, int, error) {
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, 0, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, 0, fmt.Errorf("API returned HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, 0, fmt.Errorf("parse error: %w", err)
	}

	// Extract description and version from latest
	for i := range apiResp.Results {
		if apiResp.Results[i].Latest != nil {
			if apiResp.Results[i].Description == "" {
				apiResp.Results[i].Description = apiResp.Results[i].Latest.Description
			}
			if apiResp.Results[i].Version == "" {
				apiResp.Results[i].Version = apiResp.Results[i].Latest.VersionNumber
			}
			if apiResp.Results[i].Icon == "" {
				apiResp.Results[i].Icon = apiResp.Results[i].Latest.Icon
			}
		}
	}

	return apiResp.Results, len(apiResp.Results), nil
}

// GetPackage gets details for a specific package.
func (c *Client) GetPackage(namespace, name string) (*Package, error) {
	url := fmt.Sprintf("%s?community=valheim&namespace=%s&name=%s", baseAPIURL, namespace, name)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("package not found: HTTP %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, err
	}

	if len(apiResp.Results) == 0 {
		return nil, fmt.Errorf("package not found")
	}

	pkg := apiResp.Results[0]
	if pkg.Latest != nil {
		pkg.Description = pkg.Latest.Description
		pkg.Version = pkg.Latest.VersionNumber
		pkg.Icon = pkg.Latest.Icon
	}
	return &pkg, nil
}

// ParseFullName parses "owner-name" format.
func ParseFullName(fullName string) (namespace, name string) {
	parts := strings.SplitN(fullName, "-", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return "", fullName
}
