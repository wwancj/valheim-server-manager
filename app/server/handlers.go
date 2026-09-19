package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"valheim-server-manager/app"
	"valheim-server-manager/app/config"
	"valheim-server-manager/app/thunderstore"
)

// Handlers holds references to the App and WebSocket hub for handling HTTP requests.
type Handlers struct {
	app *app.App
	hub *Hub
}

// NewHandlers creates a new Handlers instance.
func NewHandlers(a *app.App, hub *Hub) *Handlers {
	return &Handlers{app: a, hub: hub}
}

// Response is the unified API response format.
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func jsonResponse(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{
		Code:    code,
		Message: "success",
		Data:    data,
	})
}

func jsonError(w http.ResponseWriter, httpCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	json.NewEncoder(w).Encode(Response{
		Code:    1,
		Message: message,
		Data:    nil,
	})
}

// --- App info ---

func (h *Handlers) GetAppInfo(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, 0, h.app.GetAppInfo())
}

func (h *Handlers) TestBinding(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, 0, h.app.TestBinding())
}

// --- Servers ---

func (h *Handlers) GetServers(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, 0, h.app.GetServers())
}

func (h *Handlers) CreateServer(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name        string `json:"name"`
		InstallPath string `json:"installPath"`
		WorldName   string `json:"worldName"`
		Port        int    `json:"port"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, 400, "Invalid request body")
		return
	}
	server, err := h.app.CreateServer(req.Name, req.InstallPath, req.WorldName, req.Port)
	if err != nil {
		jsonError(w, 500, err.Error())
		return
	}
	jsonResponse(w, 0, server)
}

func (h *Handlers) UpdateServer(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var req struct {
		Name        string `json:"name"`
		InstallPath string `json:"installPath"`
		WorldName   string `json:"worldName"`
		Port        int    `json:"port"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, 400, "Invalid request body")
		return
	}
	server, err := h.app.UpdateServer(id, req.Name, req.InstallPath, req.WorldName, req.Port)
	if err != nil {
		jsonError(w, 500, err.Error())
		return
	}
	jsonResponse(w, 0, server)
}

func (h *Handlers) DeleteServer(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := h.app.DeleteServer(id); err != nil {
		jsonError(w, 500, err.Error())
		return
	}
	jsonResponse(w, 0, nil)
}

func (h *Handlers) ValidateServerPath(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Path string `json:"path"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, 400, "Invalid request body")
		return
	}
	jsonResponse(w, 0, h.app.ValidateServerPath(req.Path))
}

// --- Server process ---

func (h *Handlers) StartServer(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ServerDir string `json:"serverDir"`
		Name      string `json:"name"`
		World     string `json:"world"`
		Password  string `json:"password"`
		Port      int    `json:"port"`
		Public    bool   `json:"public"`
		Preset    string `json:"preset"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, 400, "Invalid request body")
		return
	}
	if err := h.app.StartServer(req.ServerDir, req.Name, req.World, req.Password, req.Port, req.Public, req.Preset); err != nil {
		jsonError(w, 500, err.Error())
		return
	}
	jsonResponse(w, 0, nil)
}

func (h *Handlers) StopServer(w http.ResponseWriter, r *http.Request) {
	if err := h.app.StopServer(); err != nil {
		jsonError(w, 500, err.Error())
		return
	}
	jsonResponse(w, 0, nil)
}

func (h *Handlers) IsServerRunning(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, 0, h.app.IsServerRunning())
}

func (h *Handlers) GetServerPID(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, 0, h.app.GetServerPID())
}

// --- SteamCMD ---

func (h *Handlers) IsSteamCMDInstalled(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, 0, h.app.IsSteamCMDInstalled())
}

func (h *Handlers) InstallSteamCMD(w http.ResponseWriter, r *http.Request) {
	if err := h.app.InstallSteamCMD(); err != nil {
		jsonError(w, 500, err.Error())
		return
	}
	jsonResponse(w, 0, nil)
}

func (h *Handlers) InstallValheimServer(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ServerDir string `json:"serverDir"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, 400, "Invalid request body")
		return
	}
	if err := h.app.InstallValheimServer(req.ServerDir); err != nil {
		jsonError(w, 500, err.Error())
		return
	}
	jsonResponse(w, 0, nil)
}

// --- Config ---

func (h *Handlers) LoadServerConfig(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	cfg, err := h.app.LoadServerConfig(id)
	if err != nil {
		jsonError(w, 500, err.Error())
		return
	}
	jsonResponse(w, 0, cfg)
}

func (h *Handlers) SaveServerConfig(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var cfg config.ServerConfig
	if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
		jsonError(w, 400, "Invalid request body")
		return
	}
	if err := h.app.SaveServerConfig(id, cfg); err != nil {
		jsonError(w, 500, err.Error())
		return
	}
	jsonResponse(w, 0, nil)
}

func (h *Handlers) GetDefaultConfig(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, 0, h.app.GetDefaultConfig())
}

func (h *Handlers) GetConfigTemplates(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, 0, h.app.GetConfigTemplates())
}

func (h *Handlers) GetCustomPresets(w http.ResponseWriter, r *http.Request) {
	presets, err := h.app.GetCustomPresets()
	if err != nil {
		jsonError(w, 500, err.Error())
		return
	}
	jsonResponse(w, 0, presets)
}

func (h *Handlers) SaveCustomPreset(w http.ResponseWriter, r *http.Request) {
	var preset config.ConfigTemplate
	if err := json.NewDecoder(r.Body).Decode(&preset); err != nil {
		jsonError(w, 400, "Invalid request body")
		return
	}
	if err := h.app.SaveCustomPreset(preset); err != nil {
		jsonError(w, 500, err.Error())
		return
	}
	jsonResponse(w, 0, nil)
}

func (h *Handlers) DeleteCustomPreset(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := h.app.DeleteCustomPreset(id); err != nil {
		jsonError(w, 500, err.Error())
		return
	}
	jsonResponse(w, 0, nil)
}

func (h *Handlers) StartServerWithConfig(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ServerDir string             `json:"serverDir"`
		Config    config.ServerConfig `json:"config"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, 400, "Invalid request body")
		return
	}
	if err := h.app.StartServerWithConfig(req.ServerDir, req.Config); err != nil {
		jsonError(w, 500, err.Error())
		return
	}
	jsonResponse(w, 0, nil)
}

func (h *Handlers) ExportStartScript(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ServerDir string             `json:"serverDir"`
		Config    config.ServerConfig `json:"config"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, 400, "Invalid request body")
		return
	}
	path, err := h.app.ExportStartScript(req.ServerDir, req.Config)
	if err != nil {
		jsonError(w, 500, err.Error())
		return
	}
	jsonResponse(w, 0, path)
}

// --- Admin lists ---

func (h *Handlers) GetAdminList(w http.ResponseWriter, r *http.Request) {
	dir := r.URL.Query().Get("serverDir")
	list, err := h.app.GetAdminList(dir)
	if err != nil {
		jsonError(w, 500, err.Error())
		return
	}
	jsonResponse(w, 0, list)
}

func (h *Handlers) SaveAdminList(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ServerDir string   `json:"serverDir"`
		List      []string `json:"list"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, 400, "Invalid request body")
		return
	}
	if err := h.app.SaveAdminList(req.ServerDir, req.List); err != nil {
		jsonError(w, 500, err.Error())
		return
	}
	jsonResponse(w, 0, nil)
}

func (h *Handlers) GetWhitelist(w http.ResponseWriter, r *http.Request) {
	dir := r.URL.Query().Get("serverDir")
	list, err := h.app.GetWhitelist(dir)
	if err != nil {
		jsonError(w, 500, err.Error())
		return
	}
	jsonResponse(w, 0, list)
}

func (h *Handlers) SaveWhitelist(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ServerDir string   `json:"serverDir"`
		List      []string `json:"list"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, 400, "Invalid request body")
		return
	}
	if err := h.app.SaveWhitelist(req.ServerDir, req.List); err != nil {
		jsonError(w, 500, err.Error())
		return
	}
	jsonResponse(w, 0, nil)
}

func (h *Handlers) GetBlacklist(w http.ResponseWriter, r *http.Request) {
	dir := r.URL.Query().Get("serverDir")
	list, err := h.app.GetBlacklist(dir)
	if err != nil {
		jsonError(w, 500, err.Error())
		return
	}
	jsonResponse(w, 0, list)
}

func (h *Handlers) SaveBlacklist(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ServerDir string   `json:"serverDir"`
		List      []string `json:"list"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, 400, "Invalid request body")
		return
	}
	if err := h.app.SaveBlacklist(req.ServerDir, req.List); err != nil {
		jsonError(w, 500, err.Error())
		return
	}
	jsonResponse(w, 0, nil)
}

// --- Logs ---

func (h *Handlers) GetLogs(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, 0, h.app.GetLogs())
}

func (h *Handlers) ClearLogs(w http.ResponseWriter, r *http.Request) {
	h.app.ClearLogs()
	jsonResponse(w, 0, nil)
}

func (h *Handlers) SetLogPaused(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Paused bool `json:"paused"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, 400, "Invalid request body")
		return
	}
	h.app.SetLogPaused(req.Paused)
	jsonResponse(w, 0, nil)
}

func (h *Handlers) ExportLogs(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Dir string `json:"dir"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, 400, "Invalid request body")
		return
	}
	path, err := h.app.ExportLogs(req.Dir)
	if err != nil {
		jsonError(w, 500, err.Error())
		return
	}
	jsonResponse(w, 0, path)
}

// --- BepInEx ---

func (h *Handlers) IsBepInExInstalled(w http.ResponseWriter, r *http.Request) {
	dir := r.URL.Query().Get("serverDir")
	jsonResponse(w, 0, h.app.IsBepInExInstalled(dir))
}

func (h *Handlers) InstallBepInEx(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ServerDir string `json:"serverDir"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, 400, "Invalid request body")
		return
	}
	if err := h.app.InstallBepInEx(req.ServerDir); err != nil {
		jsonError(w, 500, err.Error())
		return
	}
	jsonResponse(w, 0, nil)
}

func (h *Handlers) GetBepInExPlugins(w http.ResponseWriter, r *http.Request) {
	dir := r.URL.Query().Get("serverDir")
	jsonResponse(w, 0, h.app.GetBepInExPlugins(dir))
}

// --- Mods ---

func (h *Handlers) SearchMods(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	pageStr := r.URL.Query().Get("page")
	page := 1
	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}

	// Call the Thunderstore client directly to avoid Wails binding issues
	pkgs, total, err := h.app.ThunderstoreClient.SearchPackages(query, page)
	if query == "" {
		pkgs, total, err = h.app.ThunderstoreClient.GetPopularPackages(page)
	}
	if err != nil {
		jsonError(w, 500, err.Error())
		return
	}
	jsonResponse(w, 0, map[string]interface{}{
		"packages": pkgs,
		"total":    total,
	})
}

func (h *Handlers) GetModDetail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pkg, err := h.app.GetModDetail(vars["namespace"], vars["name"])
	if err != nil {
		jsonError(w, 500, err.Error())
		return
	}
	jsonResponse(w, 0, pkg)
}

func (h *Handlers) GetInstalledMods(w http.ResponseWriter, r *http.Request) {
	dir := r.URL.Query().Get("serverDir")
	jsonResponse(w, 0, h.app.GetInstalledMods(dir))
}

func (h *Handlers) InstallMod(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ServerDir string                `json:"serverDir"`
		Package   thunderstore.TSPackage `json:"package"`
		Version   thunderstore.Version   `json:"version"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, 400, "Invalid request body")
		return
	}
	if err := h.app.InstallMod(req.ServerDir, req.Package, req.Version); err != nil {
		jsonError(w, 500, err.Error())
		return
	}
	jsonResponse(w, 0, nil)
}

func (h *Handlers) UninstallMod(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ServerDir string `json:"serverDir"`
		FullName  string `json:"fullName"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, 400, "Invalid request body")
		return
	}
	if err := h.app.UninstallMod(req.ServerDir, req.FullName); err != nil {
		jsonError(w, 500, err.Error())
		return
	}
	jsonResponse(w, 0, nil)
}

func (h *Handlers) SetModEnabled(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ServerDir string `json:"serverDir"`
		FullName  string `json:"fullName"`
		Enabled   bool   `json:"enabled"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, 400, "Invalid request body")
		return
	}
	if err := h.app.SetModEnabled(req.ServerDir, req.FullName, req.Enabled); err != nil {
		jsonError(w, 500, err.Error())
		return
	}
	jsonResponse(w, 0, nil)
}

// --- System ---

func (h *Handlers) GetSystemInfo(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, 0, h.app.GetSystemInfo())
}

func (h *Handlers) GetProcessInfo(w http.ResponseWriter, r *http.Request) {
	pidStr := mux.Vars(r)["pid"]
	pid, err := strconv.ParseInt(pidStr, 10, 32)
	if err != nil {
		jsonError(w, 400, "Invalid PID")
		return
	}
	jsonResponse(w, 0, h.app.GetProcessInfo(int32(pid)))
}

func (h *Handlers) GetAppUptime(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, 0, h.app.GetAppUptime())
}

// --- Backups ---

func (h *Handlers) CreateBackup(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ServerID  string `json:"serverId"`
		ServerDir string `json:"serverDir"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, 400, "Invalid request body")
		return
	}
	b, err := h.app.CreateBackup(req.ServerID, req.ServerDir)
	if err != nil {
		jsonError(w, 500, err.Error())
		return
	}
	jsonResponse(w, 0, b)
}

func (h *Handlers) ListBackups(w http.ResponseWriter, r *http.Request) {
	serverID := mux.Vars(r)["serverId"]
	jsonResponse(w, 0, h.app.ListBackups(serverID))
}

func (h *Handlers) RestoreBackup(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var req struct {
		ServerDir string `json:"serverDir"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, 400, "Invalid request body")
		return
	}
	if err := h.app.RestoreBackup(id, req.ServerDir); err != nil {
		jsonError(w, 500, err.Error())
		return
	}
	jsonResponse(w, 0, nil)
}

func (h *Handlers) DeleteBackup(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := h.app.DeleteBackup(id); err != nil {
		jsonError(w, 500, err.Error())
		return
	}
	jsonResponse(w, 0, nil)
}

// --- Profiles ---

func (h *Handlers) GetProfiles(w http.ResponseWriter, r *http.Request) {
	dir := r.URL.Query().Get("serverDir")
	jsonResponse(w, 0, h.app.GetProfiles(dir))
}

func (h *Handlers) CreateProfile(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ServerDir   string   `json:"serverDir"`
		Name        string   `json:"name"`
		Description string   `json:"description"`
		ModList     []string `json:"modList"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, 400, "Invalid request body")
		return
	}
	p, err := h.app.CreateProfile(req.ServerDir, req.Name, req.Description, req.ModList)
	if err != nil {
		jsonError(w, 500, err.Error())
		return
	}
	jsonResponse(w, 0, p)
}

func (h *Handlers) DeleteProfile(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	dir := r.URL.Query().Get("serverDir")
	if err := h.app.DeleteProfile(dir, id); err != nil {
		jsonError(w, 500, err.Error())
		return
	}
	jsonResponse(w, 0, nil)
}

// --- Safety ---

func (h *Handlers) RunSafetyChecks(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ServerDir        string `json:"serverDir"`
		Password         string `json:"password"`
		Port             int    `json:"port"`
		BepInExInstalled bool   `json:"bepInExInstalled"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, 400, "Invalid request body")
		return
	}
	jsonResponse(w, 0, h.app.RunSafetyChecks(req.ServerDir, req.Password, req.Port, req.BepInExInstalled))
}
