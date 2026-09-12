package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"valheim-server-manager/app/models"
)

// ServerService manages Valheim server instances.
type ServerService struct {
	dataDir  string
	filePath string
	servers  []models.Server
}

// NewServerService creates a new ServerService.
// Data is stored in %APPDATA%/ValheimServerManager/servers.json
func NewServerService() *ServerService {
	appData := os.Getenv("APPDATA")
	if appData == "" {
		home, _ := os.UserHomeDir()
		appData = filepath.Join(home, ".config")
	}
	dataDir := filepath.Join(appData, "ValheimServerManager")
	svc := &ServerService{
		dataDir:  dataDir,
		filePath: filepath.Join(dataDir, "servers.json"),
	}
	svc.load()
	return svc
}

// load reads servers from disk.
func (s *ServerService) load() {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		s.servers = []models.Server{}
		return
	}
	_ = json.Unmarshal(data, &s.servers)
}

// save writes servers to disk.
func (s *ServerService) save() error {
	if err := os.MkdirAll(s.dataDir, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}
	data, err := json.MarshalIndent(s.servers, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal servers: %w", err)
	}
	return os.WriteFile(s.filePath, data, 0644)
}

// CreateServer adds a new server instance.
func (s *ServerService) CreateServer(name, installPath, worldName string, port int) (models.Server, error) {
	if name == "" {
		return models.Server{}, fmt.Errorf("server name is required")
	}
	if installPath == "" {
		return models.Server{}, fmt.Errorf("install path is required")
	}
	if port <= 0 {
		port = 2456
	}
	if worldName == "" {
		worldName = "World"
	}

	server := models.Server{
		ID:          fmt.Sprintf("srv_%d", time.Now().UnixNano()),
		Name:        name,
		InstallPath: installPath,
		WorldName:   worldName,
		Port:        port,
		Status:      models.StatusNotInstalled,
	}

	// Check if server files exist
	if s.ValidateServerPath(installPath) {
		server.Status = models.StatusInstalled
	}

	s.servers = append(s.servers, server)
	if err := s.save(); err != nil {
		return models.Server{}, err
	}
	return server, nil
}

// GetServers returns all servers.
func (s *ServerService) GetServers() []models.Server {
	// Update status based on file existence
	for i := range s.servers {
		if s.servers[i].Status == models.StatusInstalled || s.servers[i].Status == models.StatusNotInstalled {
			if s.ValidateServerPath(s.servers[i].InstallPath) {
				s.servers[i].Status = models.StatusInstalled
			} else {
				s.servers[i].Status = models.StatusNotInstalled
			}
		}
	}
	return s.servers
}

// DeleteServer removes a server by ID.
func (s *ServerService) DeleteServer(id string) error {
	for i, srv := range s.servers {
		if srv.ID == id {
			s.servers = append(s.servers[:i], s.servers[i+1:]...)
			return s.save()
		}
	}
	return fmt.Errorf("server not found: %s", id)
}

// UpdateServer updates an existing server.
func (s *ServerService) UpdateServer(id, name, installPath, worldName string, port int) (models.Server, error) {
	for i, srv := range s.servers {
		if srv.ID == id {
			if name != "" {
				s.servers[i].Name = name
			}
			if installPath != "" {
				s.servers[i].InstallPath = installPath
			}
			if worldName != "" {
				s.servers[i].WorldName = worldName
			}
			if port > 0 {
				s.servers[i].Port = port
			}
			if s.ValidateServerPath(s.servers[i].InstallPath) {
				s.servers[i].Status = models.StatusInstalled
			} else {
				s.servers[i].Status = models.StatusNotInstalled
			}
			if err := s.save(); err != nil {
				return models.Server{}, err
			}
			return s.servers[i], nil
		}
	}
	return models.Server{}, fmt.Errorf("server not found: %s", id)
}

// ValidateServerPath checks if the given path contains Valheim server files.
func (s *ServerService) ValidateServerPath(path string) bool {
	if path == "" {
		return false
	}

	// Check for valheim_server.exe (Windows) or valheim_server.x86_64 (Linux)
	var serverExe string
	if runtime.GOOS == "windows" {
		serverExe = filepath.Join(path, "valheim_server.exe")
	} else {
		serverExe = filepath.Join(path, "valheim_server.x86_64")
	}

	_, err := os.Stat(serverExe)
	return err == nil
}

// ValidateServerPathResult returns detailed validation info for the frontend.
func (s *ServerService) ValidateServerPathResult(path string) map[string]interface{} {
	result := map[string]interface{}{
		"valid":   false,
		"path":    path,
		"message": "",
	}

	if path == "" {
		result["message"] = "路径不能为空"
		return result
	}

	info, err := os.Stat(path)
	if err != nil {
		result["message"] = "路径不存在"
		return result
	}
	if !info.IsDir() {
		result["message"] = "路径不是目录"
		return result
	}

	if s.ValidateServerPath(path) {
		result["valid"] = true
		result["message"] = "已找到服务器文件"
	} else {
		result["message"] = "目录存在，但未找到 valheim_server 可执行文件（未安装或路径不正确）"
	}

	return result
}
