package app

import (
	"context"
	"fmt"
	"runtime"

	"valheim-server-manager/app/models"
	"valheim-server-manager/app/services"
)

// App is the main application struct.
// All methods on this struct are automatically bound to the frontend.
type App struct {
	ctx           context.Context
	ServerService *services.ServerService
}

// NewApp creates a new App instance.
func NewApp() *App {
	return &App{
		ServerService: services.NewServerService(),
	}
}

// Startup is called when the app starts.
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

// --- Phase 0 bindings (test) ---

func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) GetAppInfo() map[string]string {
	return map[string]string{
		"name":    "Valheim Server Manager",
		"version": "0.1.0",
		"go":      runtime.Version(),
		"os":      runtime.GOOS,
		"arch":    runtime.GOARCH,
	}
}

// --- Phase 1: Server management bindings ---

func (a *App) CreateServer(name, installPath, worldName string, port int) (models.Server, error) {
	return a.ServerService.CreateServer(name, installPath, worldName, port)
}

func (a *App) GetServers() []models.Server {
	return a.ServerService.GetServers()
}

func (a *App) DeleteServer(id string) error {
	return a.ServerService.DeleteServer(id)
}

func (a *App) UpdateServer(id, name, installPath, worldName string, port int) (models.Server, error) {
	return a.ServerService.UpdateServer(id, name, installPath, worldName, port)
}

func (a *App) ValidateServerPath(path string) map[string]interface{} {
	return a.ServerService.ValidateServerPathResult(path)
}
