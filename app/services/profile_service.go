package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ModProfile represents a collection of mods for quick switching.
type ModProfile struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Mods        []string `json:"mods"` // list of full names
	CreatedAt   string   `json:"createdAt"`
}

// ProfileService manages mod profiles.
type ProfileService struct {
	profileDir string
}

// NewProfileService creates a new profile service.
func NewProfileService(serverDir string) *ProfileService {
	profileDir := filepath.Join(serverDir, "BepInEx", "profiles")
	return &ProfileService{profileDir: profileDir}
}

// CreateProfile creates a new mod profile.
func (s *ProfileService) CreateProfile(name, description string, modList []string) (ModProfile, error) {
	if err := os.MkdirAll(s.profileDir, 0755); err != nil {
		return ModProfile{}, err
	}

	profile := ModProfile{
		ID:          fmt.Sprintf("profile_%d", time.Now().UnixNano()),
		Name:        name,
		Description: description,
		Mods:        modList,
		CreatedAt:   time.Now().Format(time.RFC3339),
	}

	data, _ := json.MarshalIndent(profile, "", "  ")
	path := filepath.Join(s.profileDir, profile.ID+".json")
	return profile, os.WriteFile(path, data, 0644)
}

// ListProfiles returns all saved profiles.
func (s *ProfileService) ListProfiles() []ModProfile {
	entries, err := os.ReadDir(s.profileDir)
	if err != nil {
		return []ModProfile{}
	}

	var profiles []ModProfile
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		data, err := os.ReadFile(filepath.Join(s.profileDir, entry.Name()))
		if err != nil {
			continue
		}
		var p ModProfile
		if json.Unmarshal(data, &p) == nil {
			profiles = append(profiles, p)
		}
	}
	return profiles
}

// DeleteProfile deletes a profile.
func (s *ProfileService) DeleteProfile(id string) error {
	path := filepath.Join(s.profileDir, id+".json")
	return os.Remove(path)
}

// ExportProfile exports a profile to a JSON file.
func (s *ProfileService) ExportProfile(id, destPath string) error {
	srcPath := filepath.Join(s.profileDir, id+".json")
	data, err := os.ReadFile(srcPath)
	if err != nil {
		return fmt.Errorf("profile not found: %w", err)
	}
	return os.WriteFile(destPath, data, 0644)
}

// ImportProfile imports a profile from a JSON file.
func (s *ProfileService) ImportProfile(filePath string) (ModProfile, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return ModProfile{}, err
	}

	var profile ModProfile
	if err := json.Unmarshal(data, &profile); err != nil {
		return ModProfile{}, fmt.Errorf("invalid profile file: %w", err)
	}

	// Generate new ID to avoid conflicts
	profile.ID = fmt.Sprintf("profile_%d", time.Now().UnixNano())
	profile.CreatedAt = time.Now().Format(time.RFC3339)

	newData, _ := json.MarshalIndent(profile, "", "  ")
	path := filepath.Join(s.profileDir, profile.ID+".json")
	return profile, os.WriteFile(path, newData, 0644)
}
