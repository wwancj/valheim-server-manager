package thunderstore

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	baseAPIURL = "https://thunderstore.io/api/experimental/package/"
)

// TSPackage represents a Thunderstore mod package.
type TSPackage struct {
	Name          string   `json:"name"`
	FullName      string   `json:"full_name"`
	Owner         string   `json:"owner"`
	Description   string   `json:"description"`
	VersionNumber string   `json:"version_number"`
	TotalDownloads int     `json:"total_downloads"`
	RatingScore   float64  `json:"rating_score"`
	PackageURL    string   `json:"package_url"`
	Icon          string   `json:"icon"`
	Latest        *Latest  `json:"latest"`
}

// Latest contains the latest version info.
type Latest struct {
	Namespace      string   `json:"namespace"`
	Name           string   `json:"name"`
	VersionNumber  string   `json:"version_number"`
	FullName       string   `json:"full_name"`
	Description    string   `json:"description"`
	Icon           string   `json:"icon"`
	Dependencies   []string `json:"dependencies"`
	DownloadURL    string   `json:"download_url"`
	Downloads      int      `json:"downloads"`
	DateCreated    string   `json:"date_created"`
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
	Next     string     `json:"next"`
	Previous string     `json:"previous"`
	Results  []TSPackage `json:"results"`
}

// Client is a Thunderstore API client.
type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{httpClient: &http.Client{Timeout: 15 * time.Second}}
}

// SearchPackages searches for mods by query.
func (c *Client) SearchPackages(query string, page int) ([]TSPackage, int, error) {
	url := fmt.Sprintf("%s?community=valheim&search=%s&page=%d", baseAPIURL, query, page)
	return c.fetchPackages(url)
}

// GetPopularPackages gets popular mods.
func (c *Client) GetPopularPackages(page int) ([]TSPackage, int, error) {
	url := fmt.Sprintf("%s?community=valheim&page=%d", baseAPIURL, page)
	return c.fetchPackages(url)
}

func (c *Client) fetchPackages(url string) ([]TSPackage, int, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("User-Agent", "ValheimServerManager/1.0")
	req.Header.Set("Accept", "application/json")
	resp, err := c.httpClient.Do(req)
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

	// Fill top-level fields from latest if empty
	for i := range apiResp.Results {
		p := &apiResp.Results[i]
		if p.Latest != nil {
			if p.Description == "" {
				p.Description = p.Latest.Description
			}
			if p.VersionNumber == "" {
				p.VersionNumber = p.Latest.VersionNumber
			}
			if p.Icon == "" {
				p.Icon = p.Latest.Icon
			}
		}
	}

	return apiResp.Results, len(apiResp.Results), nil
}

// GetPackage gets details for a specific package.
func (c *Client) GetPackage(namespace, name string) (*TSPackage, error) {
	url := fmt.Sprintf("%s?community=valheim&namespace=%s&name=%s", baseAPIURL, namespace, name)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, err
	}

	if len(apiResp.Results) == 0 {
		return nil, fmt.Errorf("not found")
	}

	pkg := &apiResp.Results[0]
	if pkg.Latest != nil {
		if pkg.Description == "" {
			pkg.Description = pkg.Latest.Description
		}
		if pkg.VersionNumber == "" {
			pkg.VersionNumber = pkg.Latest.VersionNumber
		}
		if pkg.Icon == "" {
			pkg.Icon = pkg.Latest.Icon
		}
	}
	return pkg, nil
}
