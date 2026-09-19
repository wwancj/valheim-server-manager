package app

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"valheim-server-manager/app/backup"
	"valheim-server-manager/app/config"
	"valheim-server-manager/app/events"
	"valheim-server-manager/app/mods"
	"valheim-server-manager/app/models"
	"valheim-server-manager/app/process"
	"valheim-server-manager/app/services"
	"valheim-server-manager/app/steamcmd"
	"valheim-server-manager/app/system"
	"valheim-server-manager/app/thunderstore"
)

// App is the main application struct.
type App struct {
	ctx    context.Context

	ServerService      *services.ServerService
	LogService         *services.LogService
	BepInExService     *services.BepInExService
	SafetyService      *services.SafetyService
	ProfileService     *services.ProfileService
	SteamCMDManager    *steamcmd.Manager
	ProcessManager     *process.Manager
	ConfigManager      *config.Manager
	AdminManager       *config.AdminManager
	BackupManager      *backup.Manager
	Monitor            *system.Monitor
	ThunderstoreClient *thunderstore.Client
	ModInstaller       *mods.Installer
}

// NewApp creates a new App instance.
func NewApp() *App {
	return &App{
		ServerService:      services.NewServerService(),
		LogService:         services.NewLogService(),
		BepInExService:     services.NewBepInExService(),
		SafetyService:      services.NewSafetyService(),
		SteamCMDManager:    steamcmd.NewManager(""),
		ProcessManager:     process.NewManager(),
		ConfigManager:      config.NewManager(""),
		AdminManager:       &config.AdminManager{},
		BackupManager:      backup.NewManager(""),
		Monitor:            system.NewMonitor(),
		ThunderstoreClient: thunderstore.NewClient(),
	}
}

// Startup is called when the app starts.
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	a.LogService.SetContext(ctx)
}

// SetEventEmitter wires the event emitter into all services that emit real-time events.
func (a *App) SetEventEmitter(emitter events.EventEmitter) {
	a.LogService.SetEventEmitter(emitter)
	a.BepInExService.SetEventEmitter(emitter)
	a.ProcessManager.SetEventEmitter(emitter)
	a.SteamCMDManager.SetEventEmitter(emitter)
}

// --- Phase 0 bindings ---

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

// --- Phase 1: Server management ---

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

// --- Phase 2: SteamCMD ---

func (a *App) IsSteamCMDInstalled() bool {
	return a.SteamCMDManager.IsInstalled()
}

func (a *App) InstallSteamCMD() error {
	return a.SteamCMDManager.Install(a.ctx)
}

func (a *App) InstallValheimServer(serverDir string) error {
	return a.SteamCMDManager.InstallValheimServer(a.ctx, serverDir)
}

// --- Phase 3: Server process management ---

func (a *App) StartServer(serverDir, name, world, password string, port int, public bool, preset string) error {
	cfg := config.DefaultConfig()
	cfg.Server.Name = name
	cfg.Server.World = world
	cfg.Server.Password = password
	cfg.Server.Port = port
	cfg.Server.Public = public
	cfg.World.Preset = preset
	builder := &process.CommandBuilder{Config: cfg}
	return a.ProcessManager.Start(a.ctx, serverDir, builder)
}

// StartServerWithConfig starts the server with a full configuration.
func (a *App) StartServerWithConfig(serverDir string, cfg config.ServerConfig) error {
	builder := &process.CommandBuilder{Config: cfg}
	return a.ProcessManager.Start(a.ctx, serverDir, builder)
}

func (a *App) StopServer() error {
	return a.ProcessManager.Stop(a.ctx)
}

func (a *App) IsServerRunning() bool {
	return a.ProcessManager.IsRunning()
}

func (a *App) GetServerPID() int {
	return a.ProcessManager.GetPID()
}

// --- Phase 4: Log management ---

func (a *App) GetLogs() []string {
	return a.LogService.GetLogs()
}

func (a *App) ClearLogs() {
	a.LogService.ClearLogs()
}

func (a *App) SetLogPaused(paused bool) {
	a.LogService.SetPaused(paused)
}

func (a *App) ExportLogs(dir string) (string, error) {
	return a.LogService.ExportLogs(dir)
}

// --- Phase 5: Config management ---

func (a *App) LoadServerConfig(serverID string) (config.ServerConfig, error) {
	return a.ConfigManager.LoadConfig(serverID)
}

func (a *App) SaveServerConfig(serverID string, cfg config.ServerConfig) error {
	return a.ConfigManager.SaveConfig(serverID, cfg)
}

// TestBinding is a simple test to verify bindings work.
func (a *App) TestBinding() string {
	return "binding works!"
}

// TestHTTP tests HTTP connectivity to Thunderstore.
func (a *App) TestHTTP() string {
	client := a.ThunderstoreClient
	pkgs, count, err := client.GetPopularPackages(1)
	if err != nil {
		return "ERROR: " + err.Error()
	}
	return fmt.Sprintf("OK: count=%d, len=%d", count, len(pkgs))
}

func (a *App) GetDefaultConfig() config.ServerConfig {
	return config.DefaultConfig()
}

func (a *App) GetConfigTemplates() []config.ConfigTemplate {
	return config.GetTemplates()
}

func (a *App) GetCustomPresets() ([]config.ConfigTemplate, error) {
	return a.ConfigManager.GetCustomPresets()
}

func (a *App) SaveCustomPreset(preset config.ConfigTemplate) error {
	return a.ConfigManager.SaveCustomPreset(preset)
}

func (a *App) DeleteCustomPreset(id string) error {
	return a.ConfigManager.DeleteCustomPreset(id)
}

// --- Admin/Whitelist/Blacklist management ---

func (a *App) GetAdminList(serverDir string) ([]string, error) {
	return a.AdminManager.GetList(serverDir, "adminlist.txt")
}

func (a *App) SaveAdminList(serverDir string, list []string) error {
	return a.AdminManager.SaveList(serverDir, "adminlist.txt", list)
}

func (a *App) GetWhitelist(serverDir string) ([]string, error) {
	return a.AdminManager.GetList(serverDir, "permittedlist.txt")
}

func (a *App) SaveWhitelist(serverDir string, list []string) error {
	return a.AdminManager.SaveList(serverDir, "permittedlist.txt", list)
}

func (a *App) GetBlacklist(serverDir string) ([]string, error) {
	return a.AdminManager.GetList(serverDir, "bannedlist.txt")
}

func (a *App) SaveBlacklist(serverDir string, list []string) error {
	return a.AdminManager.SaveList(serverDir, "bannedlist.txt", list)
}

// --- Phase 6: BepInEx ---

func (a *App) IsBepInExInstalled(serverDir string) bool {
	return a.BepInExService.IsInstalled(serverDir)
}

func (a *App) InstallBepInEx(serverDir string) error {
	return a.BepInExService.Install(a.ctx, serverDir)
}

func (a *App) GetBepInExPlugins(serverDir string) []string {
	return a.BepInExService.GetPlugins(serverDir)
}

// --- Phase 7: Thunderstore ---

func (a *App) SearchMods(query string, page int) ([]thunderstore.TSPackage, int, error) {
	if query == "" {
		return a.ThunderstoreClient.GetPopularPackages(page)
	}
	return a.ThunderstoreClient.SearchPackages(query, page)
}

func (a *App) GetModDetail(namespace, name string) (*thunderstore.TSPackage, error) {
	return a.ThunderstoreClient.GetPackage(namespace, name)
}

// --- Phase 8-9: Mod installation ---

func (a *App) GetInstalledMods(serverDir string) []mods.InstalledMod {
	installer := mods.NewInstaller(serverDir)
	return installer.GetInstalled()
}

func (a *App) InstallMod(serverDir string, pkg thunderstore.TSPackage, version thunderstore.Version) error {
	installer := mods.NewInstaller(serverDir)
	return installer.InstallMod(a.ctx, pkg, version)
}

func (a *App) UninstallMod(serverDir, fullName string) error {
	installer := mods.NewInstaller(serverDir)
	return installer.UninstallMod(fullName)
}

func (a *App) SetModEnabled(serverDir, fullName string, enabled bool) error {
	installer := mods.NewInstaller(serverDir)
	return installer.SetEnabled(fullName, enabled)
}

// --- Phase 11: System monitoring ---

func (a *App) GetSystemInfo() system.SystemInfo {
	return a.Monitor.GetSystemInfo()
}

func (a *App) GetProcessInfo(pid int32) system.ProcessInfo {
	return a.Monitor.GetProcessInfo(pid)
}

func (a *App) GetAppUptime() int64 {
	return a.Monitor.GetAppUptime()
}

// --- Phase 12: Backup ---

func (a *App) CreateBackup(serverID, serverDir string) (*backup.Backup, error) {
	return a.BackupManager.CreateBackup(serverID, serverDir)
}

func (a *App) ListBackups(serverID string) []backup.Backup {
	return a.BackupManager.ListBackups(serverID)
}

func (a *App) RestoreBackup(backupID, serverDir string) error {
	return a.BackupManager.RestoreBackup(backupID, serverDir)
}

func (a *App) DeleteBackup(backupID string) error {
	return a.BackupManager.DeleteBackup(backupID)
}

// --- Phase 13: Mod profiles ---

func (a *App) GetProfiles(serverDir string) []services.ModProfile {
	ps := services.NewProfileService(serverDir)
	return ps.ListProfiles()
}

func (a *App) CreateProfile(serverDir, name, description string, modList []string) (services.ModProfile, error) {
	ps := services.NewProfileService(serverDir)
	return ps.CreateProfile(name, description, modList)
}

func (a *App) DeleteProfile(serverDir, id string) error {
	ps := services.NewProfileService(serverDir)
	return ps.DeleteProfile(id)
}

// --- Export start script ---

// ExportStartScript generates a .bat startup script in the server directory.
func (a *App) ExportStartScript(serverDir string, cfg config.ServerConfig) (string, error) {
	if serverDir == "" {
		return "", fmt.Errorf("server directory is empty")
	}
	builder := &process.CommandBuilder{Config: cfg}
	args := builder.BuildArgs()

	exe := "valheim_server.exe"
	if runtime.GOOS != "windows" {
		exe = "valheim_server.x86_64"
	}

	line := exe
	for _, a := range args {
		line += " " + a
	}

	var content string
	if runtime.GOOS == "windows" {
		content = fmt.Sprintf(`@echo off
echo Starting Valheim Server...
cd /d "%s"
%s
pause
`, serverDir, line)
	} else {
		content = fmt.Sprintf(`#!/bin/bash
echo "Starting Valheim Server..."
cd "%s"
./%s
`, serverDir, line)
	}

	scriptName := "start_server.bat"
	if runtime.GOOS != "windows" {
		scriptName = "start_server.sh"
	}
	path := filepath.Join(serverDir, scriptName)
	if err := os.WriteFile(path, []byte(content), 0755); err != nil {
		return "", fmt.Errorf("failed to write script: %w", err)
	}
	return path, nil
}

// --- Phase 14: Safety checks ---

func (a *App) RunSafetyChecks(serverDir, password string, port int, bepInExInstalled bool) []services.SafetyCheck {
	return a.SafetyService.RunChecks(serverDir, password, port, bepInExInstalled, nil)
}
