import {useState, useEffect, useCallback} from 'react'
import {
    Card, Typography, Button, Table, Tag, Space, Modal, Form, Input,
    InputNumber, Popconfirm, message, Progress, Alert,
} from 'antd'
import {PlusOutlined, DeleteOutlined, ReloadOutlined, PlayCircleOutlined, StopOutlined, DownloadOutlined} from '@ant-design/icons'
import {
    CreateServer, GetServers, DeleteServer, IsSteamCMDInstalled,
    InstallSteamCMD, InstallValheimServer, StartServer, StopServer,
    IsServerRunning, RunSafetyChecks,
} from '../../../wailsjs/go/app/App'

const {Title} = Typography

function Server() {
    const [servers, setServers] = useState<any[]>([])
    const [loading, setLoading] = useState(false)
    const [modalOpen, setModalOpen] = useState(false)
    const [installing, setInstalling] = useState(false)
    const [installProgress, setInstallProgress] = useState(0)
    const [installMsg, setInstallMsg] = useState('')
    const [serverRunning, setServerRunning] = useState(false)
    const [steamcmdInstalled, setSteamcmdInstalled] = useState(false)
    const [form] = Form.useForm()

    const loadServers = useCallback(async () => {
        setLoading(true)
        try {
            const data = await GetServers()
            setServers(data || [])
        } catch {
            message.error('加载失败')
        } finally {
            setLoading(false)
        }
    }, [])

    useEffect(() => {
        loadServers()
        IsSteamCMDInstalled().then(setSteamcmdInstalled)
        IsServerRunning().then(setServerRunning)
    }, [loadServers])

    const handleCreate = async () => {
        try {
            const values = await form.validateFields()
            await CreateServer(values.name, values.installPath, values.worldName || 'World', values.port || 2456)
            message.success('添加成功')
            setModalOpen(false)
            form.resetFields()
            loadServers()
        } catch (err: any) {
            if (err.errorFields) return
            message.error(err?.message || '添加失败')
        }
    }

    const handleDelete = async (id: string) => {
        try {
            await DeleteServer(id)
            message.success('已删除')
            loadServers()
        } catch (err: any) {
            message.error(err?.message || '删除失败')
        }
    }

    const handleInstallSteamCMD = async () => {
        setInstalling(true)
        setInstallMsg('正在安装 SteamCMD...')
        try {
            await InstallSteamCMD()
            message.success('SteamCMD 安装完成')
            setSteamcmdInstalled(true)
        } catch (err: any) {
            message.error(err?.message || '安装失败')
        } finally {
            setInstalling(false)
        }
    }

    const handleInstallValheim = async (serverDir: string) => {
        setInstalling(true)
        setInstallProgress(0)
        setInstallMsg('正在安装 Valheim Server...')
        try {
            if (!steamcmdInstalled) {
                await InstallSteamCMD()
                setSteamcmdInstalled(true)
            }
            await InstallValheimServer(serverDir)
            message.success('Valheim Server 安装完成')
            loadServers()
        } catch (err: any) {
            message.error(err?.message || '安装失败')
        } finally {
            setInstalling(false)
        }
    }

    const handleStart = async (server: any) => {
        try {
            const checks = await RunSafetyChecks(server.installPath, '', server.port, false)
            const errors = checks.filter((c: any) => c.status === 'error')
            if (errors.length > 0) {
                Modal.error({
                    title: '无法启动',
                    content: errors.map((e: any) => `${e.name}: ${e.message}`).join('\n'),
                })
                return
            }
            await StartServer(server.installPath, server.name, server.worldName, '', server.port, false, 'Normal')
            message.success('服务器启动中...')
            setServerRunning(true)
        } catch (err: any) {
            message.error(err?.message || '启动失败')
        }
    }

    const handleStop = async () => {
        try {
            await StopServer()
            message.success('服务器停止中...')
            setServerRunning(false)
        } catch (err: any) {
            message.error(err?.message || '停止失败')
        }
    }

    const columns = [
        {title: '名称', dataIndex: 'name', key: 'name', render: (t: string) => <strong>{t}</strong>},
        {title: '安装路径', dataIndex: 'installPath', key: 'installPath', ellipsis: true},
        {title: '世界名称', dataIndex: 'worldName', key: 'worldName'},
        {title: '端口', dataIndex: 'port', key: 'port'},
        {
            title: '状态', dataIndex: 'status', key: 'status',
            render: (s: string) => {
                const colorMap: Record<string, string> = {
                    Installed: 'blue', Running: 'green', Stopped: 'default',
                    NotInstalled: 'orange', Error: 'red',
                }
                const labelMap: Record<string, string> = {
                    Installed: '已安装', Running: '运行中', Stopped: '已停止',
                    NotInstalled: '未安装', Error: '错误',
                }
                return <Tag color={colorMap[s] || 'default'}>{labelMap[s] || s}</Tag>
            },
        },
        {
            title: '操作', key: 'action',
            render: (_: any, record: any) => (
                <Space>
                    {record.status === 'NotInstalled' && (
                        <Button
                            type="primary"
                            icon={<DownloadOutlined/>}
                            size="small"
                            loading={installing}
                            onClick={() => handleInstallValheim(record.installPath)}
                        >
                            安装
                        </Button>
                    )}
                    {record.status === 'Installed' && !serverRunning && (
                        <Button
                            type="primary"
                            icon={<PlayCircleOutlined/>}
                            size="small"
                            onClick={() => handleStart(record)}
                        >
                            启动
                        </Button>
                    )}
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
                <Title level={2} style={{margin: 0}}>服务器管理</Title>
                <Space>
                    {serverRunning && (
                        <Button danger icon={<StopOutlined/>} onClick={handleStop}>停止服务器</Button>
                    )}
                    {!steamcmdInstalled && (
                        <Button icon={<DownloadOutlined/>} loading={installing} onClick={handleInstallSteamCMD}>
                            安装 SteamCMD
                        </Button>
                    )}
                    <Button icon={<ReloadOutlined/>} onClick={loadServers}>刷新</Button>
                    <Button type="primary" icon={<PlusOutlined/>} onClick={() => setModalOpen(true)}>添加服务器</Button>
                </Space>
            </div>

            {steamcmdInstalled && (
                <Alert message="SteamCMD 已安装" type="success" showIcon style={{marginBottom: 16}}/>
            )}

            {installing && (
                <Card style={{marginBottom: 16}}>
                    <Progress percent={installProgress} status="active"/>
                    <p>{installMsg}</p>
                </Card>
            )}

            <Card>
                <Table
                    dataSource={servers} columns={columns} rowKey="id"
                    loading={loading} pagination={false}
                    locale={{emptyText: '暂无服务器，点击上方按钮添加'}}
                />
            </Card>

            <Modal
                title="添加服务器"
                open={modalOpen}
                onOk={handleCreate}
                onCancel={() => { setModalOpen(false); form.resetFields() }}
                okText="添加" cancelText="取消"
            >
                <Form form={form} layout="vertical" style={{marginTop: 16}}>
                    <Form.Item name="name" label="服务器名称" rules={[{required: true, message: '请输入名称'}]}>
                        <Input placeholder="My Valheim Server"/>
                    </Form.Item>
                    <Form.Item name="installPath" label="安装路径" rules={[{required: true, message: '请输入路径'}]}
                               extra="Valheim Dedicated Server 安装目录">
                        <Input placeholder="D:\ValheimServer"/>
                    </Form.Item>
                    <Form.Item name="worldName" label="世界名称" initialValue="World">
                        <Input/>
                    </Form.Item>
                    <Form.Item name="port" label="端口" initialValue={2456}>
                        <InputNumber min={1} max={65535} style={{width: '100%'}}/>
                    </Form.Item>
                </Form>
            </Modal>
        </div>
    )
}

export default Server
