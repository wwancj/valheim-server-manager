import {useState, useEffect} from 'react'
import {Card, Typography, Button, Table, Space, Input, Popconfirm, message, Empty} from 'antd'
import {SaveOutlined, UndoOutlined, DeleteOutlined} from '@ant-design/icons'
import {CreateBackup, ListBackups, RestoreBackup, DeleteBackup} from '../../../wailsjs/go/app/App'

const {Title} = Typography

function Backup() {
    const [backups, setBackups] = useState<any[]>([])
    const [loading, setLoading] = useState(false)
    const [serverId, setServerId] = useState('')
    const [serverDir, setServerDir] = useState('')

    const loadBackups = async () => {
        if (!serverId) return
        setLoading(true)
        try {
            const data = await ListBackups(serverId)
            setBackups(data || [])
        } catch {
            message.error('加载备份列表失败')
        } finally {
            setLoading(false)
        }
    }

    useEffect(() => {
        if (serverId) loadBackups()
    }, [serverId])

    const handleCreate = async () => {
        if (!serverId || !serverDir) {
            message.warning('请输入服务器 ID 和服务器目录')
            return
        }
        try {
            await CreateBackup(serverId, serverDir)
            message.success('备份创建成功')
            loadBackups()
        } catch (err: any) {
            message.error(err?.message || '备份失败')
        }
    }

    const handleRestore = async (id: string) => {
        if (!serverDir) {
            message.warning('请输入服务器目录')
            return
        }
        try {
            await RestoreBackup(id, serverDir)
            message.success('备份已恢复')
        } catch (err: any) {
            message.error(err?.message || '恢复失败')
        }
    }

    const handleDelete = async (id: string) => {
        try {
            await DeleteBackup(id)
            message.success('已删除')
            loadBackups()
        } catch (err: any) {
            message.error(err?.message || '删除失败')
        }
    }

    const formatSize = (bytes: number) => {
        if (!bytes) return '0 B'
        if (bytes < 1024) return bytes + ' B'
        if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
        return (bytes / (1024 * 1024)).toFixed(1) + ' MB'
    }

    const columns = [
        {title: '备份名称', dataIndex: 'name', key: 'name'},
        {title: '创建时间', dataIndex: 'createdAt', key: 'createdAt', render: (t: string) => new Date(t).toLocaleString()},
        {title: '大小', dataIndex: 'size', key: 'size', render: (s: number) => formatSize(s)},
        {title: '包含文件', dataIndex: 'files', key: 'files', render: (f: string[]) => f?.join(', ') || '-'},
        {
            title: '操作', key: 'action',
            render: (_: any, record: any) => (
                <Space>
                    <Button icon={<UndoOutlined/>} size="small" onClick={() => handleRestore(record.id)}>恢复</Button>
                    <Popconfirm title="确定删除？" onConfirm={() => handleDelete(record.id)}>
                        <Button danger icon={<DeleteOutlined/>} size="small">删除</Button>
                    </Popconfirm>
                </Space>
            ),
        },
    ]

    return (
        <div>
            <div style={{display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16}}>
                <Title level={2} style={{margin: 0}}>备份管理</Title>
                <Button type="primary" icon={<SaveOutlined/>} onClick={handleCreate}>创建备份</Button>
            </div>
            <Space style={{marginBottom: 16}}>
                <Input placeholder="服务器 ID" value={serverId} onChange={e => setServerId(e.target.value)} style={{width: 200}}/>
                <Input placeholder="服务器目录" value={serverDir} onChange={e => setServerDir(e.target.value)} style={{width: 300}}/>
                <Button onClick={loadBackups}>刷新</Button>
            </Space>
            <Card>
                <Table
                    dataSource={backups} columns={columns} rowKey="id"
                    loading={loading} pagination={false}
                    locale={{emptyText: <Empty description="暂无备份"/>}}
                />
            </Card>
        </div>
    )
}

export default Backup
