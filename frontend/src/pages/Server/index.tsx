import {useState, useEffect, useCallback} from 'react'
import {Card, Typography, Button, Table, Tag, Space, Modal, Form, Input, InputNumber, Popconfirm, message} from 'antd'
import {PlusOutlined, DeleteOutlined, ReloadOutlined} from '@ant-design/icons'
import {serverService} from '../../services/server'
import type {Server, ServerStatus} from '../../types/server'

const {Title} = Typography

const statusColors: Record<ServerStatus, string> = {
    Running: 'green',
    Installed: 'blue',
    Stopped: 'default',
    NotInstalled: 'orange',
    Starting: 'cyan',
    Stopping: 'gold',
    Error: 'red',
    Unknown: 'default',
}

const statusLabels: Record<ServerStatus, string> = {
    Running: '运行中',
    Installed: '已安装',
    Stopped: '已停止',
    NotInstalled: '未安装',
    Starting: '启动中',
    Stopping: '停止中',
    Error: '错误',
    Unknown: '未知',
}

function Server() {
    const [servers, setServers] = useState<Server[]>([])
    const [loading, setLoading] = useState(false)
    const [modalOpen, setModalOpen] = useState(false)
    const [form] = Form.useForm()

    const loadServers = useCallback(async () => {
        setLoading(true)
        try {
            const data = await serverService.getAll()
            setServers(data || [])
        } catch (err) {
            message.error('加载服务器列表失败')
        } finally {
            setLoading(false)
        }
    }, [])

    useEffect(() => {
        loadServers()
    }, [loadServers])

    const handleCreate = async () => {
        try {
            const values = await form.validateFields()
            await serverService.create(values.name, values.installPath, values.worldName || 'World', values.port || 2456)
            message.success('服务器添加成功')
            setModalOpen(false)
            form.resetFields()
            loadServers()
        } catch (err: any) {
            if (err.errorFields) return // form validation error
            message.error(err?.message || '添加失败')
        }
    }

    const handleDelete = async (id: string) => {
        try {
            await serverService.remove(id)
            message.success('已删除')
            loadServers()
        } catch (err: any) {
            message.error(err?.message || '删除失败')
        }
    }

    const columns = [
        {
            title: '名称',
            dataIndex: 'name',
            key: 'name',
            render: (text: string) => <strong>{text}</strong>,
        },
        {
            title: '安装路径',
            dataIndex: 'installPath',
            key: 'installPath',
            ellipsis: true,
        },
        {
            title: '世界名称',
            dataIndex: 'worldName',
            key: 'worldName',
        },
        {
            title: '端口',
            dataIndex: 'port',
            key: 'port',
        },
        {
            title: '状态',
            dataIndex: 'status',
            key: 'status',
            render: (status: ServerStatus) => (
                <Tag color={statusColors[status] || 'default'}>
                    {statusLabels[status] || status}
                </Tag>
            ),
        },
        {
            title: '操作',
            key: 'action',
            render: (_: any, record: Server) => (
                <Space>
                    <Popconfirm
                        title="确定删除这个服务器？"
                        description="此操作不会删除服务器文件，仅移除管理记录。"
                        onConfirm={() => handleDelete(record.id)}
                        okText="删除"
                        cancelText="取消"
                    >
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
                    <Button icon={<ReloadOutlined/>} onClick={loadServers}>刷新</Button>
                    <Button type="primary" icon={<PlusOutlined/>} onClick={() => setModalOpen(true)}>
                        添加服务器
                    </Button>
                </Space>
            </div>

            <Card>
                <Table
                    dataSource={servers}
                    columns={columns}
                    rowKey="id"
                    loading={loading}
                    locale={{emptyText: '暂无服务器，点击上方按钮添加'}}
                    pagination={false}
                />
            </Card>

            <Modal
                title="添加服务器"
                open={modalOpen}
                onOk={handleCreate}
                onCancel={() => {
                    setModalOpen(false)
                    form.resetFields()
                }}
                okText="添加"
                cancelText="取消"
            >
                <Form form={form} layout="vertical" style={{marginTop: 16}}>
                    <Form.Item
                        name="name"
                        label="服务器名称"
                        rules={[{required: true, message: '请输入服务器名称'}]}
                    >
                        <Input placeholder="例如：My Valheim Server"/>
                    </Form.Item>
                    <Form.Item
                        name="installPath"
                        label="安装路径"
                        rules={[{required: true, message: '请输入服务器安装路径'}]}
                        extra="Valheim Dedicated Server 的安装目录，后续可通过 SteamCMD 安装"
                    >
                        <Input placeholder="例如：D:\ValheimServer"/>
                    </Form.Item>
                    <Form.Item name="worldName" label="世界名称" initialValue="World">
                        <Input placeholder="World"/>
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
