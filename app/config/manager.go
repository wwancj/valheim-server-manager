package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ServerConfig represents the full Valheim server configuration.
type ServerConfig struct {
	Server ServerSettings `json:"server"`
	Backup BackupSettings `json:"backup"`
	World  WorldSettings  `json:"world"`
}

// ServerSettings contains basic server parameters.
type ServerSettings struct {
	Name      string `json:"name"`
	World     string `json:"world"`
	Password  string `json:"password"`
	Port      int    `json:"port"`
	Public    bool   `json:"public"`
	Crossplay bool   `json:"crossplay"`
	SaveDir   string `json:"saveDir"`
	LogFile   string `json:"logFile"`
}

// BackupSettings contains backup-related parameters.
type BackupSettings struct {
	SaveInterval  int `json:"saveInterval"`  // seconds
	Count         int `json:"count"`
	ShortInterval int `json:"shortInterval"` // seconds
	LongInterval  int `json:"longInterval"`  // seconds
}

// WorldSettings contains world difficulty and rule settings.
type WorldSettings struct {
	Preset       string   `json:"preset"`
	Combat       string   `json:"combat"`
	DeathPenalty string   `json:"deathPenalty"`
	Resources    string   `json:"resources"`
	Raids        string   `json:"raids"`
	Portals      string   `json:"portals"`
	Keys         WorldKeys `json:"keys"`
}

// WorldKeys contains advanced world rule toggles.
type WorldKeys struct {
	NoBuildCost  bool `json:"noBuildCost"`
	PlayerEvents bool `json:"playerEvents"`
	PassiveMobs  bool `json:"passiveMobs"`
	NoMap        bool `json:"noMap"`
}

// ConfigTemplate represents a preset configuration template.
type ConfigTemplate struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Config      ServerConfig `json:"config"`
}

// DefaultConfig returns the default server configuration.
func DefaultConfig() ServerConfig {
	return ServerConfig{
		Server: ServerSettings{
			Name:      "Valheim Server",
			World:     "Dedicated",
			Password:  "",
			Port:      2456,
			Public:    false,
			Crossplay: true,
		},
		Backup: BackupSettings{
			SaveInterval:  1800,
			Count:         4,
			ShortInterval: 7200,
			LongInterval:  43200,
		},
		World: WorldSettings{
			Preset:       "normal",
			Combat:       "normal",
			DeathPenalty: "normal",
			Resources:    "normal",
			Raids:        "normal",
			Portals:      "normal",
			Keys:         WorldKeys{},
		},
	}
}

// GetTemplates returns built-in configuration templates.
func GetTemplates() []ConfigTemplate {
	return []ConfigTemplate{
		{
			ID: "normal", Name: "原版生存", Description: "标准游戏体验",
			Config: ServerConfig{
				Server: ServerSettings{Name: "Valheim Server", World: "Dedicated", Port: 2456, Crossplay: true},
				Backup: BackupSettings{SaveInterval: 1800, Count: 4, ShortInterval: 7200, LongInterval: 43200},
				World:  WorldSettings{Preset: "normal", Combat: "normal", DeathPenalty: "normal", Resources: "normal", Raids: "normal", Portals: "normal"},
			},
		},
		{
			ID: "casual", Name: "轻松休闲", Description: "低难度，适合休闲玩家",
			Config: ServerConfig{
				Server: ServerSettings{Name: "休闲服务器", World: "Dedicated", Port: 2456, Crossplay: true},
				Backup: BackupSettings{SaveInterval: 900, Count: 6, ShortInterval: 3600, LongInterval: 21600},
				World:  WorldSettings{Preset: "casual", Combat: "veryeasy", DeathPenalty: "casual", Resources: "more", Raids: "muchless", Portals: "casual"},
			},
		},
		{
			ID: "easy", Name: "简单模式", Description: "降低难度，适合新手",
			Config: ServerConfig{
				Server: ServerSettings{Name: "简单服务器", World: "Dedicated", Port: 2456, Crossplay: true},
				Backup: BackupSettings{SaveInterval: 1200, Count: 5, ShortInterval: 5400, LongInterval: 32400},
				World:  WorldSettings{Preset: "easy", Combat: "easy", DeathPenalty: "easy", Resources: "more", Raids: "less", Portals: "normal"},
			},
		},
		{
			ID: "hard", Name: "困难模式", Description: "提高难度，适合老玩家",
			Config: ServerConfig{
				Server: ServerSettings{Name: "困难服务器", World: "Dedicated", Port: 2456, Crossplay: true},
				Backup: BackupSettings{SaveInterval: 1800, Count: 4, ShortInterval: 7200, LongInterval: 43200},
				World:  WorldSettings{Preset: "hard", Combat: "hard", DeathPenalty: "hard", Resources: "less", Raids: "more", Portals: "hard"},
			},
		},
		{
			ID: "hardcore", Name: "极限模式", Description: "极限挑战，一命通关",
			Config: ServerConfig{
				Server: ServerSettings{Name: "极限服务器", World: "Dedicated", Port: 2456, Crossplay: true},
				Backup: BackupSettings{SaveInterval: 1800, Count: 3, ShortInterval: 7200, LongInterval: 43200},
				World:  WorldSettings{Preset: "hardcore", Combat: "veryhard", DeathPenalty: "hardcore", Resources: "muchless", Raids: "muchmore", Portals: "veryhard"},
			},
		},
		{
			ID: "immersive", Name: "沉浸模式", Description: "无地图传送门，纯探索体验",
			Config: ServerConfig{
				Server: ServerSettings{Name: "沉浸服务器", World: "Dedicated", Port: 2456, Crossplay: true},
				Backup: BackupSettings{SaveInterval: 1800, Count: 4, ShortInterval: 7200, LongInterval: 43200},
				World: WorldSettings{
					Preset: "immersive", Combat: "normal", DeathPenalty: "normal", Resources: "normal", Raids: "normal", Portals: "hard",
					Keys: WorldKeys{NoMap: true},
				},
			},
		},
		{
			ID: "builder", Name: "建筑服务器", Description: "无建造材料消耗，自由建造",
			Config: ServerConfig{
				Server: ServerSettings{Name: "建筑服务器", World: "Dedicated", Port: 2456, Crossplay: true},
				Backup: BackupSettings{SaveInterval: 900, Count: 6, ShortInterval: 3600, LongInterval: 21600},
				World: WorldSettings{
					Preset: "normal", Combat: "veryeasy", DeathPenalty: "casual", Resources: "most", Raids: "none", Portals: "casual",
					Keys: WorldKeys{NoBuildCost: true, PassiveMobs: true},
				},
			},
		},
	}
}

// Manager handles server configuration persistence.
type Manager struct {
	configDir string
}

// NewManager creates a new config manager.
func NewManager(configDir string) *Manager {
	if configDir == "" {
		appData := os.Getenv("APPDATA")
		if appData == "" {
			home, _ := os.UserHomeDir()
			appData = filepath.Join(home, ".config")
		}
		configDir = filepath.Join(appData, "ValheimServerManager", "config")
	}
	return &Manager{configDir: configDir}
}

// SaveConfig saves a server config as JSON.
func (m *Manager) SaveConfig(serverID string, cfg ServerConfig) error {
	if err := os.MkdirAll(m.configDir, 0755); err != nil {
		return err
	}
	path := filepath.Join(m.configDir, serverID+".json")
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// LoadConfig loads a server config from JSON file.
func (m *Manager) LoadConfig(serverID string) (ServerConfig, error) {
	path := filepath.Join(m.configDir, serverID+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		// Try legacy .cfg format
		cfgPath := filepath.Join(m.configDir, serverID+".cfg")
		if _, cfgErr := os.Stat(cfgPath); cfgErr == nil {
			return m.loadLegacyConfig(cfgPath)
		}
		return DefaultConfig(), nil
	}
	var cfg ServerConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return DefaultConfig(), err
	}
	return cfg, nil
}

// loadLegacyConfig loads the old .cfg format for backward compatibility.
func (m *Manager) loadLegacyConfig(path string) (ServerConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return DefaultConfig(), err
	}
	cfg := DefaultConfig()
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key, val := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		switch key {
		case "serverName":
			cfg.Server.Name = val
		case "worldName":
			cfg.Server.World = val
		case "password":
			cfg.Server.Password = val
		case "port":
			fmt.Sscanf(val, "%d", &cfg.Server.Port)
		case "public":
			cfg.Server.Public = val == "true"
		case "preset":
			cfg.World.Preset = val
		default:
			if strings.HasPrefix(key, "modifier.") {
				modKey := strings.TrimPrefix(key, "modifier.")
				switch modKey {
				case "combat":
					cfg.World.Combat = val
				case "deathpenalty":
					cfg.World.DeathPenalty = val
				case "resources":
					cfg.World.Resources = val
				case "raids":
					cfg.World.Raids = val
				case "portals":
					cfg.World.Portals = val
				}
			}
		}
	}
	return cfg, nil
}

// --- Custom Presets ---

// GetCustomPresets loads all user-saved presets.
func (m *Manager) GetCustomPresets() ([]ConfigTemplate, error) {
	dir := filepath.Join(m.configDir, "presets")
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return []ConfigTemplate{}, nil
		}
		return nil, err
	}
	var presets []ConfigTemplate
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			continue
		}
		var p ConfigTemplate
		if err := json.Unmarshal(data, &p); err != nil {
			continue
		}
		presets = append(presets, p)
	}
	return presets, nil
}

// SaveCustomPreset saves a user-created preset.
func (m *Manager) SaveCustomPreset(preset ConfigTemplate) error {
	dir := filepath.Join(m.configDir, "presets")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	if preset.ID == "" {
		preset.ID = fmt.Sprintf("custom_%d", len(preset.Name))
	}
	path := filepath.Join(dir, preset.ID+".json")
	data, err := json.MarshalIndent(preset, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// DeleteCustomPreset deletes a user-created preset.
func (m *Manager) DeleteCustomPreset(id string) error {
	path := filepath.Join(m.configDir, "presets", id+".json")
	return os.Remove(path)
}

// AdminManager manages admin/whitelist/blacklist files.
type AdminManager struct{}

// GetList reads a list file (adminlist.txt, permittedlist.txt, bannedlist.txt).
func (am *AdminManager) GetList(serverDir, filename string) ([]string, error) {
	path := filepath.Join(serverDir, filename)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}
	var list []string
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			list = append(list, line)
		}
	}
	return list, nil
}

// SaveList writes a list file.
func (am *AdminManager) SaveList(serverDir, filename string, list []string) error {
	path := filepath.Join(serverDir, filename)
	content := strings.Join(list, "\n")
	if len(list) > 0 {
		content += "\n"
	}
	return os.WriteFile(path, []byte(content), 0644)
}
