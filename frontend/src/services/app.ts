import * as App from '../../wailsjs/go/app/App'
import type {
    Server, ServerConfig, SystemInfo, ModPackage, ModVersion,
    InstalledMod, SafetyCheck, Backup,
} from '../types/server'

export const serverService = {
    create: (name: string, installPath: string, worldName: string, port: number) =>
        App.CreateServer(name, installPath, worldName, port) as Promise<Server>,
    getAll: () => App.GetServers() as Promise<Server[]>,
    remove: (id: string) => App.DeleteServer(id),
    update: (id: string, name: string, path: string, world: string, port: number) =>
        App.UpdateServer(id, name, path, world, port) as Promise<Server>,
    validatePath: (path: string) => App.ValidateServerPath(path),
}

export const steamcmdService = {
    isInstalled: () => App.IsSteamCMDInstalled(),
    install: () => App.InstallSteamCMD(),
    installValheim: (dir: string) => App.InstallValheimServer(dir),
}

export const processService = {
    start: (dir: string, name: string, world: string, pwd: string, port: number, pub: boolean, preset: string) =>
        App.StartServer(dir, name, world, pwd, port, pub, preset),
    stop: () => App.StopServer(),
    isRunning: () => App.IsServerRunning(),
    getPID: () => App.GetServerPID(),
}

export const logService = {
    getLogs: () => App.GetLogs(),
    clear: () => App.ClearLogs(),
    setPaused: (p: boolean) => App.SetLogPaused(p),
    exportLogs: (dir: string) => App.ExportLogs(dir),
}

export const configService = {
    load: (id: string) => App.LoadServerConfig(id) as Promise<ServerConfig>,
    save: (id: string, cfg: ServerConfig) => App.SaveServerConfig(id, cfg),
    getDefault: () => App.GetDefaultConfig() as Promise<ServerConfig>,
}

export const bepInExService = {
    isInstalled: (dir: string) => App.IsBepInExInstalled(dir),
    install: (dir: string) => App.InstallBepInEx(dir),
    getPlugins: (dir: string) => App.GetBepInExPlugins(dir),
}

export const thunderstoreService = {
    search: (query: string, page: number) => App.SearchMods(query, page) as unknown as Promise<ModPackage[]>,
    getDetail: (ns: string, name: string) => App.GetModDetail(ns, name) as unknown as Promise<ModPackage>,
}

export const modService = {
    getInstalled: (dir: string) => App.GetInstalledMods(dir) as Promise<InstalledMod[]>,
    install: (dir: string, pkg: ModPackage, ver: ModVersion) => App.InstallMod(dir, pkg as any, ver as any),
    uninstall: (dir: string, fullName: string) => App.UninstallMod(dir, fullName),
    setEnabled: (dir: string, fullName: string, enabled: boolean) => App.SetModEnabled(dir, fullName, enabled),
}

export const monitorService = {
    getSystem: () => App.GetSystemInfo() as Promise<SystemInfo>,
    getProcess: (pid: number) => App.GetProcessInfo(pid),
    getUptime: () => App.GetAppUptime(),
}

export const backupService = {
    create: (serverId: string, dir: string) => App.CreateBackup(serverId, dir) as Promise<Backup>,
    list: (serverId: string) => App.ListBackups(serverId) as Promise<Backup[]>,
    restore: (id: string, dir: string) => App.RestoreBackup(id, dir),
    remove: (id: string) => App.DeleteBackup(id),
}

export const safetyService = {
    run: (dir: string, pwd: string, port: number, bepInEx: boolean) =>
        App.RunSafetyChecks(dir, pwd, port, bepInEx) as Promise<SafetyCheck[]>,
}
