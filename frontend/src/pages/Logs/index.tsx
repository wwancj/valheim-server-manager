import {useState, useEffect, useRef} from 'react'
import {Card, Typography, Button, Space, Input, Tag} from 'antd'
import {ClearOutlined, PauseCircleOutlined, PlayCircleOutlined, DownloadOutlined} from '@ant-design/icons'
import {GetLogs, ClearLogs, SetLogPaused, ExportLogs} from '../../../wailsjs/go/app/App'

const {Title} = Typography

function Logs() {
    const [logs, setLogs] = useState<string[]>([])
    const [paused, setPaused] = useState(false)
    const [search, setSearch] = useState('')
    const [autoScroll, setAutoScroll] = useState(true)
    const containerRef = useRef<HTMLDivElement>(null)

    useEffect(() => {
        GetLogs().then(setLogs)
        const timer = setInterval(() => {
            if (!paused) {
                GetLogs().then(setLogs)
            }
        }, 1000)
        return () => clearInterval(timer)
    }, [paused])

    useEffect(() => {
        if (autoScroll && containerRef.current) {
            containerRef.current.scrollTop = containerRef.current.scrollHeight
        }
    }, [logs, autoScroll])

    const handleClear = async () => {
        await ClearLogs()
        setLogs([])
    }

    const handleTogglePause = () => {
        const next = !paused
        setPaused(next)
        SetLogPaused(next)
    }

    const handleExport = async () => {
        try {
            const path = await ExportLogs('')
            // message not available from wails in this context
        } catch {
            // ignore
        }
    }

    const filtered = search
        ? logs.filter(l => l.toLowerCase().includes(search.toLowerCase()))
        : logs

    const colorLine = (line: string) => {
        if (line.includes('Error') || line.includes('error') || line.includes('ERROR')) return '#ff4d4f'
        if (line.includes('Warning') || line.includes('warning')) return '#faad14'
        if (line.includes('connected') || line.includes('ready')) return '#52c41a'
        return '#d9d9d9'
    }

    return (
        <div>
            <div style={{display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16}}>
                <Title level={2} style={{margin: 0}}>实时日志</Title>
                <Space>
                    <Input.Search
                        placeholder="搜索日志..."
                        value={search}
                        onChange={e => setSearch(e.target.value)}
                        style={{width: 250}}
                        allowClear
                    />
                    <Button icon={paused ? <PlayCircleOutlined/> : <PauseCircleOutlined/>} onClick={handleTogglePause}>
                        {paused ? '继续' : '暂停'}
                    </Button>
                    <Button icon={<ClearOutlined/>} onClick={handleClear}>清空</Button>
                    <Button icon={<DownloadOutlined/>} onClick={handleExport}>导出</Button>
                </Space>
            </div>
            <Card>
                <div
                    ref={containerRef}
                    style={{
                        height: 'calc(100vh - 260px)',
                        overflow: 'auto',
                        backgroundColor: '#1a1a1a',
                        padding: 12,
                        borderRadius: 6,
                        fontFamily: 'Consolas, Monaco, monospace',
                        fontSize: 13,
                    }}
                >
                    {filtered.length === 0 ? (
                        <p style={{color: '#666'}}>等待日志输出...</p>
                    ) : (
                        filtered.map((line, i) => (
                            <div key={i} style={{color: colorLine(line), lineHeight: 1.8}}>{line}</div>
                        ))
                    )}
                </div>
                <div style={{marginTop: 8, display: 'flex', justifyContent: 'space-between'}}>
                    <Tag>{filtered.length} 条日志</Tag>
                    <Tag color={paused ? 'red' : 'green'}>{paused ? '已暂停' : '实时更新'}</Tag>
                </div>
            </Card>
        </div>
    )
}

export default Logs
