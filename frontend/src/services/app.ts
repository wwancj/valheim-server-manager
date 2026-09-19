import * as API from '../api/client'
import type {
    Server, ServerConfig, SystemInfo, ModPackage, ModVersion,
    InstalledMod, SafetyCheck, Backup,
} from '../types/server'

export const serverService = {
    create: (name: string, installPath: string, worldName: string, port: number) =>
        API.CreateServer(name, installPath, worldName, port) as Promise<Server>,
    getAll: () => API.GetServers() as Promise<Server[]>,
    remove: (id: string) => API.DeleteServer(id),
    update: (id: string, name: string, path: string, world: string, port: number) =>
        API.UpdateServer(id, name, path, world, port) as Promise<Server>,
    validatePath: (path: string) => API.ValidateServerPath(path),
}

export const steamcmdService = {
    isInstalled: () => API.IsSteamCMDInstalled(),
    install: () => API.InstallSteamCMD(),
    installValheim: (dir: string) => API.InstallValheimServer(dir),
}

export const processService = {
    start: (dir: string, name: string, world: string, pwd: string, port: number, pub: boolean, preset: string) =>
        API.StartServer(dir, name, world, pwd, port, pub, preset),
    stop: () => API.StopServer(),
    isRunning: () => API.IsServerRunning(),
    getPID: () => API.GetServerPID(),
}

export const logService = {
    getLogs: () => API.GetLogs(),
    clear: () => API.ClearLogs(),
    setPaused: (p: boolean) => API.SetLogPaused(p),
    exportLogs: (dir: string) => API.ExportLogs(dir),
}

export const configService = {
    load: (id: string) => API.LoadServerConfig(id) as Promise<ServerConfig>,
    save: (id: string, cfg: ServerConfig) => API.SaveServerConfig(id, cfg),
    getDefault: () => API.GetDefaultConfig() as Promise<ServerConfig>,
}

export const bepInExService = {
    isInstalled: (dir: string) => API.IsBepInExInstalled(dir),
    install: (dir: string) => API.InstallBepInEx(dir),
    getPlugins: (dir: string) => API.GetBepInExPlugins(dir),
}

export const thunderstoreService = {
    search: (query: string, page: number) => API.SearchMods(query, page) as Promise<{ packages: ModPackage[]; total: number }>,
    getDetail: (ns: string, name: string) => API.GetModDetail(ns, name) as Promise<ModPackage>,
}

export const modService = {
    getInstalled: (dir: string) => API.GetInstalledMods(dir) as Promise<InstalledMod[]>,
    install: (dir: string, pkg: ModPackage, ver: ModVersion) => API.InstallMod(dir, pkg as any, ver as any),
    uninstall: (dir: string, fullName: string) => API.UninstallMod(dir, fullName),
    setEnabled: (dir: string, fullName: string, enabled: boolean) => API.SetModEnabled(dir, fullName, enabled),
}

export const monitorService = {
    getSystem: () => API.GetSystemInfo() as Promise<SystemInfo>,
    getProcess: (pid: number) => API.GetProcessInfo(pid),
    getUptime: () => API.GetAppUptime(),
}

export const backupService = {
    create: (serverId: string, dir: string) => API.CreateBackup(serverId, dir) as Promise<Backup>,
    list: (serverId: string) => API.ListBackups(serverId) as Promise<Backup[]>,
    restore: (id: string, dir: string) => API.RestoreBackup(id, dir),
    remove: (id: string) => API.DeleteBackup(id),
}

export const safetyService = {
    run: (dir: string, pwd: string, port: number, bepInEx: boolean) =>
        API.RunSafetyChecks(dir, pwd, port, bepInEx) as Promise<SafetyCheck[]>,
}
