/**
 * HTTP API Client for Valheim Server Manager
 * Replaces Wails bindings with standard fetch calls
 */

const BASE_URL = '';  // Same origin (Go server serves both API and static files)

interface ApiResponse<T = any> {
    code: number;
    message: string;
    data: T;
}

async function request<T>(url: string, options?: RequestInit): Promise<T> {
    const response = await fetch(`${BASE_URL}${url}`, {
        headers: {
            'Content-Type': 'application/json',
            ...options?.headers,
        },
        ...options,
    });

    if (!response.ok) {
        const error = await response.json().catch(() => ({ message: response.statusText }));
        throw new Error(error.message || `HTTP ${response.status}`);
    }

    const result: ApiResponse<T> = await response.json();

    if (result.code !== 0) {
        throw new Error(result.message);
    }

    return result.data;
}

// ============ App ============

export function GetAppInfo(): Promise<Record<string, string>> {
    return request('/api/app/info');
}

export function TestBinding(): Promise<string> {
    return request('/api/app/test');
}

// ============ Servers ============

export function GetServers(): Promise<any[]> {
    return request('/api/servers');
}

export function CreateServer(
    name: string,
    installPath: string,
    worldName: string,
    port: number
): Promise<any> {
    return request('/api/servers', {
        method: 'POST',
        body: JSON.stringify({ name, installPath, worldName, port }),
    });
}

export function UpdateServer(
    id: string,
    name: string,
    installPath: string,
    worldName: string,
    port: number
): Promise<any> {
    return request(`/api/servers/${id}`, {
        method: 'PUT',
        body: JSON.stringify({ name, installPath, worldName, port }),
    });
}

export function DeleteServer(id: string): Promise<void> {
    return request(`/api/servers/${id}`, {
        method: 'DELETE',
    });
}

export function ValidateServerPath(path: string): Promise<Record<string, any>> {
    return request(`/api/servers/validate`, {
        method: 'POST',
        body: JSON.stringify({ path }),
    });
}

// ============ Process ============

export function StartServer(
    serverDir: string,
    name: string,
    world: string,
    password: string,
    port: number,
    publicServer: boolean,
    preset: string
): Promise<void> {
    return request('/api/servers/start', {
        method: 'POST',
        body: JSON.stringify({
            serverDir,
            name,
            world,
            password,
            port,
            public: publicServer,
            preset,
        }),
    });
}

export function StopServer(): Promise<void> {
    return request('/api/servers/stop', { method: 'POST' });
}

export function IsServerRunning(): Promise<boolean> {
    return request('/api/servers/status');
}

export function GetServerPID(): Promise<number> {
    return request('/api/servers/pid');
}

// ============ SteamCMD ============

export function IsSteamCMDInstalled(): Promise<boolean> {
    return request('/api/steamcmd/status');
}

export function InstallSteamCMD(): Promise<void> {
    return request('/api/steamcmd/install', { method: 'POST' });
}

export function InstallValheimServer(serverDir: string): Promise<void> {
    return request('/api/valheim/install', {
        method: 'POST',
        body: JSON.stringify({ serverDir }),
    });
}

// ============ Config ============

export function LoadServerConfig(serverID: string): Promise<any> {
    return request(`/api/config/${serverID}`);
}

export function SaveServerConfig(serverID: string, config: any): Promise<void> {
    return request(`/api/config/${serverID}`, {
        method: 'PUT',
        body: JSON.stringify(config),
    });
}

export function GetDefaultConfig(): Promise<any> {
    return request('/api/config/default');
}

export function GetConfigTemplates(): Promise<any[]> {
    return request('/api/config/templates');
}

export function GetCustomPresets(): Promise<any[]> {
    return request('/api/config/custom-presets');
}

export function SaveCustomPreset(preset: any): Promise<void> {
    return request('/api/config/custom-presets', {
        method: 'POST',
        body: JSON.stringify(preset),
    });
}

export function DeleteCustomPreset(id: string): Promise<void> {
    return request(`/api/config/custom-presets/${id}`, { method: 'DELETE' });
}

export function StartServerWithConfig(serverDir: string, config: any): Promise<void> {
    return request('/api/servers/start-with-config', {
        method: 'POST',
        body: JSON.stringify({ serverDir, config }),
    });
}

export function ExportStartScript(serverDir: string, config: any): Promise<string> {
    return request('/api/servers/export-script', {
        method: 'POST',
        body: JSON.stringify({ serverDir, config }),
    });
}

// ============ Admin Lists ============

export function GetAdminList(serverDir: string): Promise<string[]> {
    return request(`/api/admin/list?serverDir=${encodeURIComponent(serverDir)}`);
}

export function SaveAdminList(serverDir: string, list: string[]): Promise<void> {
    return request('/api/admin/list', {
        method: 'PUT',
        body: JSON.stringify({ serverDir, list }),
    });
}

export function GetWhitelist(serverDir: string): Promise<string[]> {
    return request(`/api/admin/whitelist?serverDir=${encodeURIComponent(serverDir)}`);
}

export function SaveWhitelist(serverDir: string, list: string[]): Promise<void> {
    return request('/api/admin/whitelist', {
        method: 'PUT',
        body: JSON.stringify({ serverDir, list }),
    });
}

export function GetBlacklist(serverDir: string): Promise<string[]> {
    return request(`/api/admin/blacklist?serverDir=${encodeURIComponent(serverDir)}`);
}

export function SaveBlacklist(serverDir: string, list: string[]): Promise<void> {
    return request('/api/admin/blacklist', {
        method: 'PUT',
        body: JSON.stringify({ serverDir, list }),
    });
}

// ============ Logs ============

export function GetLogs(): Promise<string[]> {
    return request('/api/logs');
}

export function ClearLogs(): Promise<void> {
    return request('/api/logs', { method: 'DELETE' });
}

export function SetLogPaused(paused: boolean): Promise<void> {
    return request('/api/logs/pause', {
        method: 'PUT',
        body: JSON.stringify({ paused }),
    });
}

export function ExportLogs(dir: string): Promise<string> {
    return request('/api/logs/export', {
        method: 'POST',
        body: JSON.stringify({ dir }),
    });
}

// ============ BepInEx ============

export function IsBepInExInstalled(serverDir: string): Promise<boolean> {
    return request(`/api/bepinex/status?serverDir=${encodeURIComponent(serverDir)}`);
}

export function InstallBepInEx(serverDir: string): Promise<void> {
    return request('/api/bepinex/install', {
        method: 'POST',
        body: JSON.stringify({ serverDir }),
    });
}

export function GetBepInExPlugins(serverDir: string): Promise<string[]> {
    return request(`/api/bepinex/plugins?serverDir=${encodeURIComponent(serverDir)}`);
}

// ============ Mods ============

export function SearchMods(query: string, page: number): Promise<{ packages: any[]; total: number }> {
    return request(`/api/mods/search?query=${encodeURIComponent(query)}&page=${page}`);
}

export function GetModDetail(namespace: string, name: string): Promise<any> {
    return request(`/api/mods/detail/${namespace}/${name}`);
}

export function GetInstalledMods(serverDir: string): Promise<any[]> {
    return request(`/api/mods/installed?serverDir=${encodeURIComponent(serverDir)}`);
}

export function InstallMod(serverDir: string, pkg: any, version: any): Promise<void> {
    return request('/api/mods/install', {
        method: 'POST',
        body: JSON.stringify({ serverDir, package: pkg, version }),
    });
}

export function UninstallMod(serverDir: string, fullName: string): Promise<void> {
    return request('/api/mods/installed', {
        method: 'DELETE',
        body: JSON.stringify({ serverDir, fullName }),
    });
}

export function SetModEnabled(serverDir: string, fullName: string, enabled: boolean): Promise<void> {
    return request('/api/mods/toggle', {
        method: 'PUT',
        body: JSON.stringify({ serverDir, fullName, enabled }),
    });
}

// ============ System ============

export function GetSystemInfo(): Promise<any> {
    return request('/api/system');
}

export function GetProcessInfo(pid: number): Promise<any> {
    return request(`/api/system/process/${pid}`);
}

export function GetAppUptime(): Promise<number> {
    return request('/api/system/uptime');
}

// ============ Backups ============

export function CreateBackup(serverID: string, serverDir: string): Promise<any> {
    return request('/api/backups', {
        method: 'POST',
        body: JSON.stringify({ serverId: serverID, serverDir }),
    });
}

export function ListBackups(serverID: string): Promise<any[]> {
    return request(`/api/backups/${serverID}`);
}

export function RestoreBackup(backupID: string, serverDir: string): Promise<void> {
    return request(`/api/backups/${backupID}/restore`, {
        method: 'POST',
        body: JSON.stringify({ serverDir }),
    });
}

export function DeleteBackup(backupID: string): Promise<void> {
    return request(`/api/backups/${backupID}`, {
        method: 'DELETE',
    });
}

// ============ Profiles ============

export function GetProfiles(serverDir: string): Promise<any[]> {
    return request(`/api/profiles?serverDir=${encodeURIComponent(serverDir)}`);
}

export function CreateProfile(
    serverDir: string,
    name: string,
    description: string,
    modList: string[]
): Promise<any> {
    return request('/api/profiles', {
        method: 'POST',
        body: JSON.stringify({ serverDir, name, description, modList }),
    });
}

export function DeleteProfile(serverDir: string, id: string): Promise<void> {
    return request(`/api/profiles/${id}?serverDir=${encodeURIComponent(serverDir)}`, {
        method: 'DELETE',
    });
}

// ============ Safety ============

export function RunSafetyChecks(
    serverDir: string,
    password: string,
    port: number,
    bepInExInstalled: boolean
): Promise<any[]> {
    return request('/api/safety/check', {
        method: 'POST',
        body: JSON.stringify({ serverDir, password, port, bepInExInstalled }),
    });
}
