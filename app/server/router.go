package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Router wraps gorilla/mux Router with API route definitions.
type Router struct {
	*mux.Router
}

// NewRouter creates a new Router with all API routes registered.
func NewRouter(h *Handlers) *Router {
	r := mux.NewRouter()

	// API routes
	api := r.PathPrefix("/api").Subrouter()

	// App info
	api.HandleFunc("/app/info", h.GetAppInfo).Methods("GET")
	api.HandleFunc("/app/test", h.TestBinding).Methods("GET")

	// Servers
	api.HandleFunc("/servers", h.GetServers).Methods("GET")
	api.HandleFunc("/servers", h.CreateServer).Methods("POST")
	api.HandleFunc("/servers/{id}", h.UpdateServer).Methods("PUT")
	api.HandleFunc("/servers/{id}", h.DeleteServer).Methods("DELETE")
	api.HandleFunc("/servers/{id}/validate", h.ValidateServerPath).Methods("POST")

	// Server process
	api.HandleFunc("/servers/start", h.StartServer).Methods("POST")
	api.HandleFunc("/servers/stop", h.StopServer).Methods("POST")
	api.HandleFunc("/servers/status", h.IsServerRunning).Methods("GET")
	api.HandleFunc("/servers/pid", h.GetServerPID).Methods("GET")

	// SteamCMD
	api.HandleFunc("/steamcmd/status", h.IsSteamCMDInstalled).Methods("GET")
	api.HandleFunc("/steamcmd/install", h.InstallSteamCMD).Methods("POST")
	api.HandleFunc("/valheim/install", h.InstallValheimServer).Methods("POST")

	// Config (specific routes first, then parameterized)
	api.HandleFunc("/config/default", h.GetDefaultConfig).Methods("GET")
	api.HandleFunc("/config/templates", h.GetConfigTemplates).Methods("GET")
	api.HandleFunc("/config/custom-presets", h.GetCustomPresets).Methods("GET")
	api.HandleFunc("/config/custom-presets", h.SaveCustomPreset).Methods("POST")
	api.HandleFunc("/config/custom-presets/{id}", h.DeleteCustomPreset).Methods("DELETE")
	api.HandleFunc("/config/{id}", h.LoadServerConfig).Methods("GET")
	api.HandleFunc("/config/{id}", h.SaveServerConfig).Methods("PUT")
	api.HandleFunc("/servers/start-with-config", h.StartServerWithConfig).Methods("POST")
	api.HandleFunc("/servers/export-script", h.ExportStartScript).Methods("POST")

	// Admin lists
	api.HandleFunc("/admin/list", h.GetAdminList).Methods("GET")
	api.HandleFunc("/admin/list", h.SaveAdminList).Methods("PUT")
	api.HandleFunc("/admin/whitelist", h.GetWhitelist).Methods("GET")
	api.HandleFunc("/admin/whitelist", h.SaveWhitelist).Methods("PUT")
	api.HandleFunc("/admin/blacklist", h.GetBlacklist).Methods("GET")
	api.HandleFunc("/admin/blacklist", h.SaveBlacklist).Methods("PUT")

	// Logs
	api.HandleFunc("/logs", h.GetLogs).Methods("GET")
	api.HandleFunc("/logs", h.ClearLogs).Methods("DELETE")
	api.HandleFunc("/logs/pause", h.SetLogPaused).Methods("PUT")
	api.HandleFunc("/logs/export", h.ExportLogs).Methods("POST")

	// BepInEx
	api.HandleFunc("/bepinex/status", h.IsBepInExInstalled).Methods("GET")
	api.HandleFunc("/bepinex/install", h.InstallBepInEx).Methods("POST")
	api.HandleFunc("/bepinex/plugins", h.GetBepInExPlugins).Methods("GET")

	// Mods
	api.HandleFunc("/mods/search", h.SearchMods).Methods("GET")
	api.HandleFunc("/mods/detail/{namespace}/{name}", h.GetModDetail).Methods("GET")
	api.HandleFunc("/mods/installed", h.GetInstalledMods).Methods("GET")
	api.HandleFunc("/mods/install", h.InstallMod).Methods("POST")
	api.HandleFunc("/mods/installed", h.UninstallMod).Methods("DELETE")
	api.HandleFunc("/mods/toggle", h.SetModEnabled).Methods("PUT")

	// System
	api.HandleFunc("/system", h.GetSystemInfo).Methods("GET")
	api.HandleFunc("/system/process/{pid}", h.GetProcessInfo).Methods("GET")
	api.HandleFunc("/system/uptime", h.GetAppUptime).Methods("GET")

	// Backups
	api.HandleFunc("/backups", h.CreateBackup).Methods("POST")
	api.HandleFunc("/backups/{serverId}", h.ListBackups).Methods("GET")
	api.HandleFunc("/backups/{id}/restore", h.RestoreBackup).Methods("POST")
	api.HandleFunc("/backups/{id}", h.DeleteBackup).Methods("DELETE")

	// Profiles
	api.HandleFunc("/profiles", h.GetProfiles).Methods("GET")
	api.HandleFunc("/profiles", h.CreateProfile).Methods("POST")
	api.HandleFunc("/profiles/{id}", h.DeleteProfile).Methods("DELETE")

	// Safety
	api.HandleFunc("/safety/check", h.RunSafetyChecks).Methods("POST")

	// WebSocket
	api.HandleFunc("/ws/events", h.HandleWebSocket).Methods("GET")

	return &Router{Router: r}
}

// ServeSPA sets up SPA fallback routing: static files are served if they exist,
// otherwise index.html is returned for client-side routing.
func (r *Router) ServeSPA(fileServer http.Handler) {
	r.PathPrefix("/").Handler(fileServer)
}
