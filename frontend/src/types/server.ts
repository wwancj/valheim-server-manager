export type ServerStatus =
    | 'Stopped' | 'Starting' | 'Running' | 'Stopping'
    | 'Error' | 'Unknown' | 'Installed' | 'NotInstalled'

export interface Server {
    id: string
    name: string
    installPath: string
    worldName: string
    port: number
    status: ServerStatus
}

export interface ServerPathValidation {
    valid: boolean
    path: string
    message: string
}

export interface SystemInfo {
    cpuPercent: number
    memUsed: number
    memTotal: number
    memPercent: number
    goVersion: string
    os: string
    arch: string
}

export interface ModPackage {
    name: string
    full_name: string
    owner: string
    description: string
    version_number: string
    downloads: number
    rating_score: number
    categories: string[]
    package_url: string
    icon: string
    versions: ModVersion[]
}

export interface ModVersion {
    version_number: string
    download_url: string
    dependencies: string[]
    file_size: number
    date_created: string
}

export interface InstalledMod {
    name: string
    fullName: string
    version: string
    description: string
    enabled: boolean
    dependencies: string[]
    installedAt: string
}

export interface ServerConfig {
    serverName: string
    worldName: string
    password: string
    port: number
    public: boolean
    preset: string
    modifiers: Record<string, string>
}

export interface SafetyCheck {
    name: string
    status: 'ok' | 'warning' | 'error'
    message: string
}

export interface Backup {
    id: string
    serverId: string
    name: string
    path: string
    size: number
    createdAt: string
    files: string[]
}
