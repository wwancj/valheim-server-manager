export type ServerStatus =
    | 'Stopped'
    | 'Starting'
    | 'Running'
    | 'Stopping'
    | 'Error'
    | 'Unknown'
    | 'Installed'
    | 'NotInstalled'

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
