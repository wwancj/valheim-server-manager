package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ServerConfig represents the Valheim server configuration.
type ServerConfig struct {
	ServerName string            `json:"serverName"`
	WorldName  string            `json:"worldName"`
	Password   string            `json:"password"`
	Port       int               `json:"port"`
	Public     bool              `json:"public"`
	Preset     string            `json:"preset"`
	Modifiers  map[string]string `json:"modifiers"`
}

// DefaultConfig returns a default server configuration.
func DefaultConfig() ServerConfig {
	return ServerConfig{
		ServerName: "Valheim Server",
		WorldName:  "World",
		Password:   "",
		Port:       2456,
		Public:     false,
		Preset:     "Normal",
		Modifiers: map[string]string{
			"combat":      "Normal",
			"deathpenalty": "Normal",
			"resources":   "Normal",
			"raids":       "Normal",
			"portals":     "Normal",
		},
	}
}

// Manager handles server configuration.
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

// BuildArgs converts a ServerConfig into command-line arguments.
func BuildArgs(cfg ServerConfig) []string {
	args := []string{
		"-name", cfg.ServerName,
		"-world", cfg.WorldName,
		"-port", fmt.Sprintf("%d", cfg.Port),
		"-password", cfg.Password,
	}
	if cfg.Public {
		args = append(args, "-public", "1")
	} else {
		args = append(args, "-public", "0")
	}
	if cfg.Preset != "" {
		args = append(args, "-preset", cfg.Preset)
	}
	return args
}

// SaveConfig saves a server config to a file.
func (m *Manager) SaveConfig(serverID string, cfg ServerConfig) error {
	if err := os.MkdirAll(m.configDir, 0755); err != nil {
		return err
	}
	path := filepath.Join(m.configDir, serverID+".cfg")
	data := fmt.Sprintf("serverName=%s\nworldName=%s\npassword=%s\nport=%d\npublic=%v\npreset=%s\n",
		cfg.ServerName, cfg.WorldName, cfg.Password, cfg.Port, cfg.Public, cfg.Preset)
	for k, v := range cfg.Modifiers {
		data += fmt.Sprintf("modifier.%s=%s\n", k, v)
	}
	return os.WriteFile(path, []byte(data), 0644)
}

// LoadConfig loads a server config from a file.
func (m *Manager) LoadConfig(serverID string) (ServerConfig, error) {
	path := filepath.Join(m.configDir, serverID+".cfg")
	data, err := os.ReadFile(path)
	if err != nil {
		return DefaultConfig(), nil // Return default if not found
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
			cfg.ServerName = val
		case "worldName":
			cfg.WorldName = val
		case "password":
			cfg.Password = val
		case "port":
			fmt.Sscanf(val, "%d", &cfg.Port)
		case "public":
			cfg.Public = val == "true"
		case "preset":
			cfg.Preset = val
		default:
			if strings.HasPrefix(key, "modifier.") {
				modKey := strings.TrimPrefix(key, "modifier.")
				cfg.Modifiers[modKey] = val
			}
		}
	}
	return cfg, nil
}
