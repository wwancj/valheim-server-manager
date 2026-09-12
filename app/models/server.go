package models

// ServerStatus represents the current state of a Valheim server.
type ServerStatus string

const (
	StatusStopped  ServerStatus = "Stopped"
	StatusStarting ServerStatus = "Starting"
	StatusRunning  ServerStatus = "Running"
	StatusStopping ServerStatus = "Stopping"
	StatusError    ServerStatus = "Error"
	StatusUnknown  ServerStatus = "Unknown"
	StatusInstalled ServerStatus = "Installed"
	StatusNotInstalled ServerStatus = "NotInstalled"
)

// Server represents a Valheim dedicated server instance.
type Server struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	InstallPath string       `json:"installPath"`
	WorldName   string       `json:"worldName"`
	Port        int          `json:"port"`
	Status      ServerStatus `json:"status"`
}
