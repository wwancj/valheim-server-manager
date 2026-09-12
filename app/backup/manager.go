package backup

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// Backup represents a server backup.
type Backup struct {
	ID        string `json:"id"`
	ServerID  string `json:"serverId"`
	Name      string `json:"name"`
	Path      string `json:"path"`
	Size      int64  `json:"size"`
	CreatedAt string `json:"createdAt"`
	Files     []string `json:"files"`
}

// Manager handles backup operations.
type Manager struct {
	backupDir string
}

// NewManager creates a new backup manager.
func NewManager(backupDir string) *Manager {
	if backupDir == "" {
		appData := os.Getenv("APPDATA")
		if appData == "" {
			home, _ := os.UserHomeDir()
			appData = filepath.Join(home, ".config")
		}
		backupDir = filepath.Join(appData, "ValheimServerManager", "backups")
	}
	return &Manager{backupDir: backupDir}
}

// CreateBackup creates a backup of server world, config, and mod data.
func (m *Manager) CreateBackup(serverID, serverDir string) (*Backup, error) {
	timestamp := time.Now().Format("2006-01-02_150405")
	backupName := fmt.Sprintf("%s_%s", serverID, timestamp)
	backupPath := filepath.Join(m.backupDir, backupName)

	if err := os.MkdirAll(backupPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create backup directory: %w", err)
	}

	var backedUpFiles []string

	// Backup world files
	worldDir := filepath.Join(serverDir, "saves")
	if err := copyDir(worldDir, filepath.Join(backupPath, "saves")); err == nil {
		backedUpFiles = append(backedUpFiles, "saves")
	}

	// Backup server config
	configDir := filepath.Join(serverDir, "BepInEx", "config")
	if err := copyDir(configDir, filepath.Join(backupPath, "config")); err == nil {
		backedUpFiles = append(backedUpFiles, "config")
	}

	// Backup mod list
	modsFile := filepath.Join(serverDir, "BepInEx", "installed_mods.json")
	if _, err := os.Stat(modsFile); err == nil {
		copyFile(modsFile, filepath.Join(backupPath, "installed_mods.json"))
		backedUpFiles = append(backedUpFiles, "installed_mods.json")
	}

	// Calculate size
	size := dirSize(backupPath)

	backup := &Backup{
		ID:        backupName,
		ServerID:  serverID,
		Name:      fmt.Sprintf("Backup %s", timestamp),
		Path:      backupPath,
		Size:      size,
		CreatedAt: time.Now().Format(time.RFC3339),
		Files:     backedUpFiles,
	}

	// Save backup metadata
	metaPath := filepath.Join(backupPath, "backup.json")
	metaData, _ := json.MarshalIndent(backup, "", "  ")
	os.WriteFile(metaPath, metaData, 0644)

	return backup, nil
}

// ListBackups returns all backups for a server.
func (m *Manager) ListBackups(serverID string) []Backup {
	entries, err := os.ReadDir(m.backupDir)
	if err != nil {
		return []Backup{}
	}

	var backups []Backup
	prefix := serverID + "_"
	for _, entry := range entries {
		if !entry.IsDir() || len(entry.Name()) < len(prefix) {
			continue
		}
		if entry.Name()[:len(prefix)] != prefix {
			continue
		}
		metaPath := filepath.Join(m.backupDir, entry.Name(), "backup.json")
		data, err := os.ReadFile(metaPath)
		if err != nil {
			continue
		}
		var b Backup
		if json.Unmarshal(data, &b) == nil {
			backups = append(backups, b)
		}
	}
	return backups
}

// RestoreBackup restores a backup to the server directory.
func (m *Manager) RestoreBackup(backupID, serverDir string) error {
	backupPath := filepath.Join(m.backupDir, backupID)
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return fmt.Errorf("backup not found: %s", backupID)
	}

	// Restore saves
	savesDir := filepath.Join(backupPath, "saves")
	if _, err := os.Stat(savesDir); err == nil {
		copyDir(savesDir, filepath.Join(serverDir, "saves"))
	}

	// Restore config
	configDir := filepath.Join(backupPath, "config")
	if _, err := os.Stat(configDir); err == nil {
		copyDir(configDir, filepath.Join(serverDir, "BepInEx", "config"))
	}

	return nil
}

// DeleteBackup deletes a backup.
func (m *Manager) DeleteBackup(backupID string) error {
	backupPath := filepath.Join(m.backupDir, backupID)
	return os.RemoveAll(backupPath)
}

// copyDir copies a directory recursively.
func copyDir(src, dest string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !srcInfo.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	os.MkdirAll(dest, 0755)

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		destPath := filepath.Join(dest, entry.Name())

		if entry.IsDir() {
			copyDir(srcPath, destPath)
		} else {
			copyFile(srcPath, destPath)
		}
	}
	return nil
}

// copyFile copies a single file.
func copyFile(src, dest string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	os.MkdirAll(filepath.Dir(dest), 0755)
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

// dirSize calculates the total size of a directory.
func dirSize(path string) int64 {
	var size int64
	filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size
}
