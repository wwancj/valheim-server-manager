package app

import (
	"context"
	"fmt"
	"runtime"

	"valheim-server-manager/app/backup"
	"valheim-server-manager/app/config"
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

	ServerService   *services.ServerService
	LogService      *services.LogService
	BepInExService  *services.BepInExService
	SafetyService   *services.SafetyService
	ProfileService  *services.ProfileService
	SteamCMDManager *steamcmd.Manager
	ProcessManager  *process.Manager
	ConfigManager   *config.Manager
	BackupManager   *backup.Manager
	Monitor         *system.Monitor
	ThunderstoreClient *thunderstore.Client
	ModInstaller    *mods.Installer
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
	builder := &process.CommandBuilder{
		Name:     name,
		Port:     port,
		World:    world,
		Password: password,
		Public:   public,
		Preset:   preset,
	}
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

func (a *App) GetDefaultConfig() config.ServerConfig {
	return config.DefaultConfig()
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

// --- Phase 14: Safety checks ---

func (a *App) RunSafetyChecks(serverDir, password string, port int, bepInExInstalled bool) []services.SafetyCheck {
	return a.SafetyService.RunChecks(serverDir, password, port, bepInExInstalled, nil)
}
