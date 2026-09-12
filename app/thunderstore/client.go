package thunderstore

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	baseAPIURL = "https://thunderstore.io/c/valheim"
)

// Package represents a Thunderstore mod package.
type Package struct {
	Name        string   `json:"name"`
	FullName    string   `json:"full_name"`
	Owner       string   `json:"owner"`
	Description string   `json:"description"`
	Version     string   `json:"version_number"`
	Downloads   int      `json:"downloads"`
	Rating      float64  `json:"rating_score"`
	Categories  []string `json:"categories"`
	WebsiteURL  string   `json:"package_url"`
	Versions    []Version `json:"versions"`
	Tags        []string `json:"tags"`
	Icon        string   `json:"icon"`
}

// Version represents a specific version of a mod.
type Version struct {
	VersionNumber  string   `json:"version_number"`
	DownloadURL    string   `json:"download_url"`
	Dependencies   []string `json:"dependencies"`
	FileSize       int      `json:"file_size"`
	DateCreated    string   `json:"date_created"`
}

// APIResponse represents the Thunderstore API response.
type APIResponse struct {
	Count    int       `json:"count"`
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
	url := fmt.Sprintf("%s/api/v1/package/?q=%s&page=%d", baseAPIURL, query, page)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, 0, fmt.Errorf("API returned HTTP %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, 0, fmt.Errorf("failed to parse response: %w", err)
	}

	return apiResp.Results, apiResp.Count, nil
}

// GetPackage gets details for a specific package.
func (c *Client) GetPackage(namespace, name string) (*Package, error) {
	url := fmt.Sprintf("%s/api/v1/package/%s/%s/", baseAPIURL, namespace, name)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get package: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("package not found: HTTP %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var pkg Package
	if err := json.Unmarshal(body, &pkg); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &pkg, nil
}

// GetPopularPackages gets popular mods for the Valheim community.
func (c *Client) GetPopularPackages(page int) ([]Package, int, error) {
	url := fmt.Sprintf("%s/api/v1/package/?ordering=downloads&page=%d", baseAPIURL, page)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch packages: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, 0, fmt.Errorf("API returned HTTP %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, 0, fmt.Errorf("failed to parse response: %w", err)
	}

	return apiResp.Results, apiResp.Count, nil
}

// ParseFullName parses "owner-name" format.
func ParseFullName(fullName string) (namespace, name string) {
	parts := strings.SplitN(fullName, "-", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return "", fullName
}
