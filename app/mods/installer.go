package mods

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"

	"valheim-server-manager/app/thunderstore"
	"valheim-server-manager/app/utils"
)

// InstalledMod represents a mod installed on the server.
type InstalledMod struct {
	Name         string   `json:"name"`
	FullName     string   `json:"fullName"`
	Version      string   `json:"version"`
	Description  string   `json:"description"`
	Enabled      bool     `json:"enabled"`
	Dependencies []string `json:"dependencies"`
	InstalledAt  string   `json:"installedAt"`
}

// Installer manages mod installation.
type Installer struct {
	pluginsDir string
	modsFile   string
	installed  []InstalledMod
}

// NewInstaller creates a new mod installer.
func NewInstaller(serverDir string) *Installer {
	pluginsDir := filepath.Join(serverDir, "BepInEx", "plugins")
	modsFile := filepath.Join(serverDir, "BepInEx", "installed_mods.json")
	installer := &Installer{pluginsDir: pluginsDir, modsFile: modsFile}
	installer.loadInstalled()
	return installer
}

func (i *Installer) loadInstalled() {
	data, err := os.ReadFile(i.modsFile)
	if err != nil {
		i.installed = []InstalledMod{}
		return
	}
	json.Unmarshal(data, &i.installed)
}

func (i *Installer) saveInstalled() error {
	os.MkdirAll(filepath.Dir(i.modsFile), 0755)
	data, err := json.MarshalIndent(i.installed, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(i.modsFile, data, 0644)
}

func (i *Installer) GetInstalled() []InstalledMod {
	return i.installed
}

func (i *Installer) InstallMod(ctx context.Context, pkg thunderstore.TSPackage, version thunderstore.Version) error {
	for _, mod := range i.installed {
		if mod.FullName == pkg.FullName && mod.Version == version.VersionNumber {
			return fmt.Errorf("mod %s v%s already installed", pkg.FullName, version.VersionNumber)
		}
	}

	os.MkdirAll(i.pluginsDir, 0755)

	tmpZip := filepath.Join(os.TempDir(), fmt.Sprintf("%s_%s.zip", pkg.FullName, version.VersionNumber))
	if err := utils.DownloadFile(version.DownloadURL, tmpZip); err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	defer os.Remove(tmpZip)

	if err := utils.Unzip(tmpZip, i.pluginsDir); err != nil {
		return fmt.Errorf("extract failed: %w", err)
	}

	mod := InstalledMod{
		Name:         pkg.Name,
		FullName:     pkg.FullName,
		Version:      version.VersionNumber,
		Description:  pkg.Description,
		Enabled:      true,
		Dependencies: version.Dependencies,
		InstalledAt:  time.Now().Format(time.RFC3339),
	}
	i.installed = append(i.installed, mod)

	if ctx != nil {
		wailsRuntime.EventsEmit(ctx, "mod:installed", map[string]interface{}{
			"name": mod.Name, "version": mod.Version,
		})
	}

	return i.saveInstalled()
}

func (i *Installer) UninstallMod(fullName string) error {
	for idx, mod := range i.installed {
		if mod.FullName == fullName {
			i.installed = append(i.installed[:idx], i.installed[idx+1:]...)
			return i.saveInstalled()
		}
	}
	return fmt.Errorf("mod not found: %s", fullName)
}

func (i *Installer) SetEnabled(fullName string, enabled bool) error {
	for idx, mod := range i.installed {
		if mod.FullName == fullName {
			i.installed[idx].Enabled = enabled
			dllPath := filepath.Join(i.pluginsDir, mod.Name+".dll")
			disabledPath := dllPath + ".disabled"
			if enabled {
				if _, err := os.Stat(disabledPath); err == nil {
					os.Rename(disabledPath, dllPath)
				}
			} else {
				if _, err := os.Stat(dllPath); err == nil {
					os.Rename(dllPath, disabledPath)
				}
			}
			return i.saveInstalled()
		}
	}
	return fmt.Errorf("mod not found: %s", fullName)
}

func (i *Installer) IsModInstalled(fullName string) bool {
	for _, mod := range i.installed {
		if mod.FullName == fullName {
			return true
		}
	}
	return false
}
