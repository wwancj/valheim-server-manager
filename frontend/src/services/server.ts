import {CreateServer, GetServers, DeleteServer, UpdateServer, ValidateServerPath} from '../../wailsjs/go/app/App'
import type {Server, ServerPathValidation} from '../types/server'

export const serverService = {
    async create(name: string, installPath: string, worldName: string, port: number): Promise<Server> {
        const result = await CreateServer(name, installPath, worldName, port)
        return result as unknown as Server
    },

    async getAll(): Promise<Server[]> {
        const result = await GetServers()
        return (result || []) as unknown as Server[]
    },

    async remove(id: string): Promise<void> {
        await DeleteServer(id)
    },

    async update(id: string, name: string, installPath: string, worldName: string, port: number): Promise<Server> {
        const result = await UpdateServer(id, name, installPath, worldName, port)
        return result as unknown as Server
    },

    async validatePath(path: string): Promise<ServerPathValidation> {
        const result = await ValidateServerPath(path)
        return result as unknown as ServerPathValidation
    },
}
