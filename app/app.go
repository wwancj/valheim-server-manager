package app

import (
	"context"
	"fmt"
	"runtime"
)

// App is the main application struct.
// All methods on this struct are automatically bound to the frontend.
type App struct {
	ctx context.Context
}

// NewApp creates a new App instance.
func NewApp() *App {
	return &App{}
}

// Startup is called when the app starts. The context is saved
// so we can call the runtime methods.
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name (Phase 0 binding test).
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// GetAppInfo returns basic application info (Phase 0 binding test).
func (a *App) GetAppInfo() map[string]string {
	return map[string]string{
		"name":    "Valheim Server Manager",
		"version": "0.1.0",
		"go":      runtime.Version(),
		"os":      runtime.GOOS,
		"arch":    runtime.GOARCH,
	}
}
