import {createContext, useContext, useState, useEffect, useRef, ReactNode} from 'react'
import {wsClient} from '../api/websocket'

export interface ProgressState {
    visible: boolean
    title: string
    percent: number
    message: string
    done: boolean
    error: string
}

const defaultState: ProgressState = {
    visible: false,
    title: '',
    percent: 0,
    message: '',
    done: false,
    error: '',
}

interface ProgressContextType {
    progress: ProgressState
    startProgress: (title: string) => void
    resetProgress: () => void
}

const ProgressContext = createContext<ProgressContextType>({
    progress: defaultState,
    startProgress: () => {},
    resetProgress: () => {},
})

export function ProgressProvider({children}: {children: ReactNode}) {
    const [progress, setProgress] = useState<ProgressState>(defaultState)
    const progressRef = useRef(progress)
    progressRef.current = progress

    const startProgress = (title: string) => {
        setProgress({
            visible: true,
            title,
            percent: 0,
            message: '准备中...',
            done: false,
            error: '',
        })
    }

    const resetProgress = () => {
        setProgress(defaultState)
    }

    useEffect(() => {
        // SteamCMD 安装进度
        const unsub1 = wsClient.on('steamcmd:progress', (data: any) => {
            setProgress(prev => ({
                ...prev,
                visible: true,
                title: '安装 SteamCMD',
                percent: data.percent || 0,
                message: data.message || '下载中...',
                done: data.percent >= 100,
            }))
        })

        // Valheim 服务器安装进度
        const unsub2 = wsClient.on('steamcmd:valheim:progress', (data: any) => {
            setProgress(prev => ({
                ...prev,
                visible: true,
                title: '安装 Valheim 服务器',
                percent: data.percent || 0,
                message: data.message || '下载中...',
                done: data.percent >= 100,
            }))
        })

        // SteamCMD 输出日志
        const unsub3 = wsClient.on('steamcmd:valheim:output', (data: any) => {
            if (data?.line) {
                setProgress(prev => ({
                    ...prev,
                    message: data.line,
                }))
            }
        })

        // BepInEx 安装进度
        const unsub4 = wsClient.on('bepinex:progress', (data: any) => {
            setProgress(prev => ({
                ...prev,
                visible: true,
                title: '安装 BepInEx',
                percent: data.percent || 0,
                message: data.message || '安装中...',
                done: data.percent >= 100,
            }))
        })

        // 服务器状态变更
        const unsub5 = wsClient.on('server:status', (data: any) => {
            if (data?.status === 'started') {
                setProgress(prev => ({...prev, done: true, message: '服务器已启动'}))
            }
        })

        // Mod 安装完成
        const unsub6 = wsClient.on('mod:installed', (data: any) => {
            setProgress(prev => ({
                ...prev,
                done: true,
                percent: 100,
                message: `Mod 安装完成: ${data?.name || ''}`,
            }))
        })

        return () => {
            unsub1()
            unsub2()
            unsub3()
            unsub4()
            unsub5()
            unsub6()
        }
    }, [])

    // 自动隐藏完成的进度条
    useEffect(() => {
        if (progress.done) {
            const timer = setTimeout(() => {
                setProgress(prev => ({...prev, visible: false}))
            }, 3000)
            return () => clearTimeout(timer)
        }
    }, [progress.done])

    return (
        <ProgressContext.Provider value={{progress, startProgress, resetProgress}}>
            {children}
        </ProgressContext.Provider>
    )
}

export function useProgress() {
    return useContext(ProgressContext)
}
