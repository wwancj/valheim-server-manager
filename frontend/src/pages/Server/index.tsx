import {useState, useEffect, useCallback} from 'react'
import {
    Card, Typography, Button, Table, Tag, Space, Modal, Form, Input,
    InputNumber, Popconfirm, message, Alert, Select, Divider, Radio, Switch,
} from 'antd'
import {
    PlusOutlined, DeleteOutlined, ReloadOutlined, PlayCircleOutlined,
    StopOutlined, DownloadOutlined, SettingOutlined, ExportOutlined,
} from '@ant-design/icons'
import {
    CreateServer, GetServers, DeleteServer, IsSteamCMDInstalled,
    InstallSteamCMD, InstallValheimServer, StopServer,
    IsServerRunning, GetConfigTemplates, GetCustomPresets, LoadServerConfig, StartServerWithConfig,
    ExportStartScript,
} from '../../api/client'
import {useProgress} from '../../contexts/ProgressContext'

const {Title, Text} = Typography

function Server() {
    const [servers, setServers] = useState<any[]>([])
    const [loading, setLoading] = useState(false)
    const [modalOpen, setModalOpen] = useState(false)
    const [startModalOpen, setStartModalOpen] = useState(false)
    const [selectedServer, setSelectedServer] = useState<any>(null)
    const [serverRunning, setServerRunning] = useState(false)
    const [steamcmdInstalled, setSteamcmdInstalled] = useState(false)
    const [templates, setTemplates] = useState<any[]>([])
    const [customPresets, setCustomPresets] = useState<any[]>([])
    const [selectedTemplate, setSelectedTemplate] = useState('normal')
    const [savedConfig, setSavedConfig] = useState<any>(null)
    const [useCustomConfig, setUseCustomConfig] = useState(false)
    const [crossplay, setCrossplay] = useState(false)
    const [form] = Form.useForm()
    const [startForm] = Form.useForm()
    const {startProgress} = useProgress()

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
        GetConfigTemplates().then(setTemplates)
        GetCustomPresets().then(setCustomPresets).catch(() => {})
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
        startProgress('安装 SteamCMD')
        try {
            await InstallSteamCMD()
            message.success('SteamCMD 安装完成')
            setSteamcmdInstalled(true)
        } catch (err: any) {
            message.error(err?.message || '安装失败')
        }
    }

    const handleInstallValheim = async (serverDir: string) => {
        startProgress('安装 Valheim 服务器')
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
        }
    }

    // 打开启动配置弹窗
    const handleOpenStart = async (server: any) => {
        setSelectedServer(server)
        setUseCustomConfig(false)
        setSelectedTemplate('normal')
        setCrossplay(false)

        // 尝试加载已保存的配置
        try {
            const cfg = await LoadServerConfig(server.id)
            setSavedConfig(cfg)
            setCrossplay(cfg.server?.crossplay ?? false)
            startForm.setFieldsValue({
                password: cfg.server?.password || '',
                preset: cfg.world?.preset || 'normal',
            })
        } catch {
            setSavedConfig(null)
            startForm.setFieldsValue({password: '', preset: 'normal'})
        }

        setStartModalOpen(true)
    }

    // 确认启动
    const handleConfirmStart = async () => {
        if (!selectedServer) return
        const values = await startForm.validateFields()

        let cfg: any

        if (useCustomConfig && savedConfig) {
            // 使用已保存的配置
            cfg = JSON.parse(JSON.stringify(savedConfig))
            if (values.password) cfg.server.password = values.password
            cfg.server.crossplay = crossplay
        } else {
            // 使用选择的模板（内置或自定义）
            const tpl = templates.find(t => t.id === selectedTemplate) || customPresets.find(t => t.id === selectedTemplate)
            cfg = tpl ? JSON.parse(JSON.stringify(tpl.config)) : {
                server: {name: selectedServer.name, world: selectedServer.worldName || 'Dedicated', port: selectedServer.port || 2456, crossplay: false},
                backup: {saveInterval: 1800, count: 4, shortInterval: 7200, longInterval: 43200},
                world: {preset: selectedTemplate, combat: 'normal', deathPenalty: 'normal', resources: 'normal', raids: 'normal', portals: 'normal', keys: {}},
            }
            cfg.server.password = values.password || ''
            cfg.server.port = selectedServer.port || 2456
            cfg.server.crossplay = crossplay
        }

        try {
            await StartServerWithConfig(selectedServer.installPath, cfg)
            message.success('服务器启动中...')
            setServerRunning(true)
            setStartModalOpen(false)
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

    const handleExportScript = async () => {
        if (!selectedServer) return
        try {
            const values = await startForm.validateFields()
            let cfg: any
            if (useCustomConfig && savedConfig) {
                cfg = JSON.parse(JSON.stringify(savedConfig))
                if (values.password) cfg.server.password = values.password
                cfg.server.crossplay = crossplay
            } else {
                const tpl = templates.find(t => t.id === selectedTemplate) || customPresets.find(t => t.id === selectedTemplate)
                cfg = tpl ? JSON.parse(JSON.stringify(tpl.config)) : {
                    server: {name: selectedServer.name, world: selectedServer.worldName || 'Dedicated', port: selectedServer.port || 2456, crossplay: false},
                    backup: {saveInterval: 1800, count: 4, shortInterval: 7200, longInterval: 43200},
                    world: {preset: selectedTemplate, combat: 'normal', deathPenalty: 'normal', resources: 'normal', raids: 'normal', portals: 'normal', keys: {}},
                }
                cfg.server.password = values.password || ''
                cfg.server.port = selectedServer.port || 2456
                cfg.server.crossplay = crossplay
            }
            const path = await ExportStartScript(selectedServer.installPath, cfg)
            message.success(`启动脚本已导出: ${path}`)
        } catch (err: any) {
            if (err.errorFields) return
            message.error(err?.message || '导出失败')
        }
    }

    const getTemplateDesc = (id: string) => {
        const map: Record<string, string> = {
            normal: '标准游戏体验',
            casual: '低难度，适合休闲',
            easy: '降低难度，适合新手',
            hard: '提高难度，适合老玩家',
            hardcore: '极限挑战，一命通关',
            immersive: '无地图传送门，纯探索',
            builder: '无建造消耗，自由建造',
        }
        return map[id] || ''
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
                            onClick={() => handleOpenStart(record)}
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
                        <Button icon={<DownloadOutlined/>} onClick={handleInstallSteamCMD}>
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

            <Card>
                <Table
                    dataSource={servers} columns={columns} rowKey="id"
                    loading={loading} pagination={false}
                    locale={{emptyText: '暂无服务器，点击上方按钮添加'}}
                />
            </Card>

            {/* 添加服务器弹窗 */}
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

            {/* 启动配置弹窗 */}
            <Modal
                title={<Space><SettingOutlined/> 启动服务器配置</Space>}
                open={startModalOpen}
                onOk={handleConfirmStart}
                onCancel={() => setStartModalOpen(false)}
                okText="启动服务器" cancelText="取消"
                width={600}
                footer={(_, {OkBtn, CancelBtn}) => (
                    <div style={{display: 'flex', justifyContent: 'space-between'}}>
                        <Button icon={<ExportOutlined/>} onClick={handleExportScript}>
                            导出启动脚本
                        </Button>
                        <Space>
                            <CancelBtn/>
                            <OkBtn/>
                        </Space>
                    </div>
                )}
            >
                {selectedServer && (
                    <div style={{marginBottom: 16}}>
                        <Text type="secondary">服务器: </Text>
                        <Text strong>{selectedServer.name}</Text>
                        <Text type="secondary"> | 路径: </Text>
                        <Text>{selectedServer.installPath}</Text>
                    </div>
                )}

                <Form form={startForm} layout="vertical">
                    <Form.Item name="password" label="服务器密码">
                        <Input.Password placeholder="留空则无密码"/>
                    </Form.Item>

                    <Form.Item label="跨平台联机 (Crossplay)">
                        <div style={{display: 'flex', alignItems: 'center', gap: 12}}>
                            <Switch checked={crossplay} onChange={setCrossplay}/>
                            <Text type="secondary">
                                {crossplay
                                    ? '已开启：通过 PlayFab 中继，支持 Steam/Xbox 跨平台，需用好友邀请加入'
                                    : '已关闭：传统 UDP 直连，支持 IP 直接连接（如 127.0.0.1:2456）'}
                            </Text>
                        </div>
                    </Form.Item>

                    {!crossplay && (
                        <Alert
                            type="info"
                            showIcon
                            style={{marginBottom: 16}}
                            message="连接方式：在 Valheim 客户端选择「加入游戏」→ 输入 IP 和端口"
                            description={
                                <div>
                                    <div>地址: <Text code>127.0.0.1:2456</Text>（本机）或 <Text code>你的公网IP:2456</Text>（远程）</div>
                                    <div style={{marginTop: 4}}>注意：远程连接需要在路由器/防火墙开放 UDP 2456-2457 端口</div>
                                </div>
                            }
                        />
                    )}

                    <Divider>选择游戏模式</Divider>

                    {savedConfig && (
                        <div style={{marginBottom: 12}}>
                            <Radio checked={useCustomConfig} onChange={() => setUseCustomConfig(true)}>
                                <Tag color="blue">已保存配置</Tag>
                                <Text type="secondary">使用之前在配置页面保存的设置</Text>
                            </Radio>
                        </div>
                    )}

                    <Radio checked={!useCustomConfig} onChange={() => setUseCustomConfig(false)}>
                        <Text strong>使用预设模板</Text>
                    </Radio>

                    {!useCustomConfig && (
                        <div style={{marginLeft: 24, marginTop: 8}}>
                            <div style={{display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: 8}}>
                                {templates.map(tpl => (
                                    <div
                                        key={tpl.id}
                                        onClick={() => {
                                            setSelectedTemplate(tpl.id)
                                            startForm.setFieldsValue({preset: tpl.id})
                                        }}
                                        style={{
                                            padding: '8px 12px',
                                            border: `2px solid ${selectedTemplate === tpl.id ? '#d4a017' : '#333'}`,
                                            borderRadius: 8,
                                            cursor: 'pointer',
                                            background: selectedTemplate === tpl.id ? 'rgba(212,160,23,0.1)' : 'transparent',
                                            transition: 'all 0.2s',
                                        }}
                                    >
                                        <div style={{fontWeight: 600}}>{tpl.name}</div>
                                        <div style={{fontSize: 12, color: '#888'}}>{tpl.description}</div>
                                    </div>
                                ))}
                                {customPresets.map(tpl => (
                                    <div
                                        key={tpl.id}
                                        onClick={() => {
                                            setSelectedTemplate(tpl.id)
                                            startForm.setFieldsValue({preset: tpl.id})
                                        }}
                                        style={{
                                            padding: '8px 12px',
                                            border: `2px solid ${selectedTemplate === tpl.id ? '#d4a017' : '#333'}`,
                                            borderRadius: 8,
                                            cursor: 'pointer',
                                            background: selectedTemplate === tpl.id ? 'rgba(212,160,23,0.1)' : 'transparent',
                                            transition: 'all 0.2s',
                                        }}
                                    >
                                        <div style={{fontWeight: 600}}><Tag color="gold">自定义</Tag> {tpl.name}</div>
                                        <div style={{fontSize: 12, color: '#888'}}>{tpl.description}</div>
                                    </div>
                                ))}
                            </div>
                        </div>
                    )}

                    {useCustomConfig && savedConfig && (
                        <Card size="small" style={{marginLeft: 24, marginTop: 8, background: '#1a1a2e'}}>
                            <Space direction="vertical" size={4}>
                                <Text>预设: <Tag>{savedConfig.world?.preset}</Tag></Text>
                                <Text>战斗: <Tag>{savedConfig.world?.combat}</Tag> 资源: <Tag>{savedConfig.world?.resources}</Tag></Text>
                                <Text>袭击: <Tag>{savedConfig.world?.raids}</Tag> 传送门: <Tag>{savedConfig.world?.portals}</Tag></Text>
                            </Space>
                        </Card>
                    )}
                </Form>
            </Modal>
        </div>
    )
}

export default Server
