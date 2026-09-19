import {useState, useEffect} from 'react'
import {
    Card, Typography, Form, Input, InputNumber, Switch, Select, Button,
    Space, message, Divider, Radio, List, Popconfirm, Tag, Modal, Empty,
} from 'antd'
import {
    PlusOutlined, DeleteOutlined, SaveOutlined, EditOutlined,
    CopyOutlined, PlayCircleOutlined,
} from '@ant-design/icons'
import {
    GetDefaultConfig, GetConfigTemplates, GetServers,
    GetCustomPresets, SaveCustomPreset, DeleteCustomPreset,
    StartServerWithConfig, SaveServerConfig, LoadServerConfig,
} from '../../api/client'

const {Title, Text} = Typography

const COMBAT_OPTIONS = [
    {value: 'veryeasy', label: '非常简单'}, {value: 'easy', label: '简单'},
    {value: 'normal', label: '正常'}, {value: 'hard', label: '困难'}, {value: 'veryhard', label: '非常困难'},
]
const DEATH_OPTIONS = [
    {value: 'casual', label: '休闲'}, {value: 'veryeasy', label: '非常简单'},
    {value: 'easy', label: '简单'}, {value: 'normal', label: '正常'},
    {value: 'hard', label: '困难'}, {value: 'hardcore', label: '极限'},
]
const RESOURCE_OPTIONS = [
    {value: 'muchless', label: '少很多'}, {value: 'less', label: '少'},
    {value: 'normal', label: '正常'}, {value: 'more', label: '多'},
    {value: 'muchmore', label: '多很多'}, {value: 'most', label: '最大'},
]
const RAID_OPTIONS = [
    {value: 'none', label: '无'}, {value: 'muchless', label: '少很多'},
    {value: 'less', label: '少'}, {value: 'normal', label: '正常'},
    {value: 'more', label: '多'}, {value: 'muchmore', label: '多很多'},
]
const PORTAL_OPTIONS = [
    {value: 'casual', label: '休闲'}, {value: 'normal', label: '正常'},
    {value: 'hard', label: '困难'}, {value: 'veryhard', label: '非常困难'},
]

function Config() {
    const [form] = Form.useForm()
    const [builtInPresets, setBuiltInPresets] = useState<any[]>([])
    const [customPresets, setCustomPresets] = useState<any[]>([])
    const [selectedPreset, setSelectedPreset] = useState<any>(null)
    const [editing, setEditing] = useState(false) // 是否在编辑模式
    const [isNew, setIsNew] = useState(false) // 是否新建
    const [servers, setServers] = useState<any[]>([])
    const [launchModal, setLaunchModal] = useState(false)
    const [launchServerId, setLaunchServerId] = useState('')
    const [launchPassword, setLaunchPassword] = useState('')

    useEffect(() => {
        loadPresets()
        GetServers().then(d => setServers(d || []))
    }, [])

    const loadPresets = () => {
        GetConfigTemplates().then(setBuiltInPresets)
        GetCustomPresets().then(setCustomPresets).catch(() => setCustomPresets([]))
    }

    // 选中一个预设，展示详情
    const handleSelect = (preset: any, isCustom: boolean) => {
        setSelectedPreset({...preset, isCustom})
        setEditing(false)
        setIsNew(false)
        fillForm(preset.config)
    }

    // 用表单填充配置
    const fillForm = (cfg: any) => {
        form.setFieldsValue({
            name: cfg.server?.name || '',
            world: cfg.server?.world || 'Dedicated',
            password: cfg.server?.password || '',
            port: cfg.server?.port || 2456,
            public: cfg.server?.public || false,
            crossplay: cfg.server?.crossplay ?? true,
            saveDir: cfg.server?.saveDir || '',
            logFile: cfg.server?.logFile || '',
            saveInterval: cfg.backup?.saveInterval ? cfg.backup.saveInterval / 60 : 30,
            backupCount: cfg.backup?.count || 4,
            preset: cfg.world?.preset || 'normal',
            combat: cfg.world?.combat || 'normal',
            deathPenalty: cfg.world?.deathPenalty || 'normal',
            resources: cfg.world?.resources || 'normal',
            raids: cfg.world?.raids || 'normal',
            portals: cfg.world?.portals || 'normal',
            noBuildCost: cfg.world?.keys?.noBuildCost || false,
            playerEvents: cfg.world?.keys?.playerEvents || false,
            passiveMobs: cfg.world?.keys?.passiveMobs || false,
            noMap: cfg.world?.keys?.noMap || false,
        })
    }

    // 从表单读取配置
    const readForm = () => {
        const v = form.getFieldsValue()
        return {
            server: {
                name: v.name, world: v.world, password: v.password,
                port: v.port, public: v.public, crossplay: v.crossplay,
                saveDir: v.saveDir, logFile: v.logFile,
            },
            backup: {
                saveInterval: (v.saveInterval || 30) * 60,
                count: v.backupCount || 4,
                shortInterval: 7200, longInterval: 43200,
            },
            world: {
                preset: v.preset, combat: v.combat, deathPenalty: v.deathPenalty,
                resources: v.resources, raids: v.raids, portals: v.portals,
                keys: {
                    noBuildCost: v.noBuildCost, playerEvents: v.playerEvents,
                    passiveMobs: v.passiveMobs, noMap: v.noMap,
                },
            },
        }
    }

    // 新建预设
    const handleNew = () => {
        setSelectedPreset(null)
        setIsNew(true)
        setEditing(true)
        const cfg = readForm()
        fillForm({
            server: {name: '新预设', world: 'Dedicated', port: 2456, crossplay: true},
            backup: {saveInterval: 1800, count: 4, shortInterval: 7200, longInterval: 43200},
            world: {preset: 'normal', combat: 'normal', deathPenalty: 'normal', resources: 'normal', raids: 'normal', portals: 'normal'},
        })
    }

    // 保存预设
    const handleSave = async () => {
        const v = form.getFieldsValue()
        const presetName = v.name?.trim()
        if (!presetName) { message.warning('请填写服务器名称（作为预设名称）'); return }

        const id = isNew
            ? 'custom_' + presetName.replace(/\s+/g, '_').toLowerCase() + '_' + Date.now()
            : selectedPreset?.id || 'custom_' + presetName.replace(/\s+/g, '_').toLowerCase()

        const preset = {
            id,
            name: presetName,
            description: `战斗:${v.combat} 资源:${v.resources} 袭击:${v.raids}`,
            config: readForm(),
        }

        try {
            await SaveCustomPreset(preset)
            message.success('预设已保存')
            setEditing(false)
            setIsNew(false)
            setSelectedPreset({...preset, isCustom: true})
            loadPresets()
        } catch (err: any) {
            message.error(err?.message || '保存失败')
        }
    }

    // 删除自定义预设
    const handleDelete = async (id: string) => {
        try {
            await DeleteCustomPreset(id)
            message.success('已删除')
            if (selectedPreset?.id === id) {
                setSelectedPreset(null)
                setEditing(false)
            }
            loadPresets()
        } catch {
            message.error('删除失败')
        }
    }

    // 复制预设（基于内置预设创建自定义副本）
    const handleCopy = (preset: any) => {
        const id = 'custom_' + preset.name + '_' + Date.now()
        setSelectedPreset({...preset, id, isCustom: true})
        setIsNew(true)
        setEditing(true)
        fillForm(preset.config)
        message.info('已复制，请修改后保存')
    }

    // 用此预设启动服务器
    const handleLaunch = () => {
        if (servers.length === 0) { message.warning('请先添加服务器'); return }
        setLaunchServerId(servers[0]?.id || '')
        setLaunchPassword('')
        setLaunchModal(true)
    }

    const handleConfirmLaunch = async () => {
        if (!launchServerId) { message.warning('请选择服务器'); return }
        const server = servers.find(s => s.id === launchServerId)
        if (!server) return
        try {
            const cfg = readForm()
            if (launchPassword) cfg.server.password = launchPassword
            await StartServerWithConfig(server.installPath, cfg)
            message.success('服务器启动中...')
            setLaunchModal(false)
        } catch (err: any) {
            message.error(err?.message || '启动失败')
        }
    }

    // 同时保存配置到指定服务器
    const handleSaveToServer = async () => {
        if (!launchServerId) return
        try {
            await SaveServerConfig(launchServerId, readForm())
            message.success('配置已保存到服务器')
        } catch (err: any) {
            message.error(err?.message || '保存失败')
        }
    }

    const renderPresetItem = (preset: any, isCustom: boolean) => (
        <List.Item
            key={preset.id}
            onClick={() => handleSelect(preset, isCustom)}
            style={{
                cursor: 'pointer',
                background: selectedPreset?.id === preset.id ? 'rgba(212,160,23,0.15)' : 'transparent',
                borderLeft: selectedPreset?.id === preset.id ? '3px solid #d4a017' : '3px solid transparent',
                padding: '8px 12px',
                transition: 'all 0.2s',
            }}
            actions={isCustom ? [
                <Popconfirm key="del" title="确定删除？" onConfirm={(e) => { e?.stopPropagation(); handleDelete(preset.id) }}>
                    <Button danger size="small" icon={<DeleteOutlined/>} onClick={e => e.stopPropagation()}/>
                </Popconfirm>,
            ] : [
                <Button key="copy" size="small" icon={<CopyOutlined/>} onClick={(e) => { e.stopPropagation(); handleCopy(preset) }}>
                    复制
                </Button>,
            ]}
        >
            <List.Item.Meta
                title={<Space>
                    {isCustom && <Tag color="gold">自定义</Tag>}
                    <Text strong>{preset.name}</Text>
                </Space>}
                description={<Text type="secondary" ellipsis>{preset.description}</Text>}
            />
        </List.Item>
    )

    return (
        <div style={{display: 'flex', gap: 16, height: 'calc(100vh - 120px)'}}>
            {/* 左侧：预设列表 */}
            <Card
                title="预设列表"
                size="small"
                style={{width: 320, flexShrink: 0, overflow: 'auto'}}
                extra={<Button type="primary" size="small" icon={<PlusOutlined/>} onClick={handleNew}>新建</Button>}
            >
                <Text type="secondary" style={{fontSize: 12, marginBottom: 8, display: 'block'}}>内置预设</Text>
                <List
                    size="small"
                    dataSource={builtInPresets}
                    renderItem={item => renderPresetItem(item, false)}
                    style={{marginBottom: 12}}
                />
                {customPresets.length > 0 && (
                    <>
                        <Divider style={{margin: '8px 0'}}/>
                        <Text type="secondary" style={{fontSize: 12, marginBottom: 8, display: 'block'}}>自定义预设</Text>
                        <List
                            size="small"
                            dataSource={customPresets}
                            renderItem={item => renderPresetItem(item, true)}
                        />
                    </>
                )}
                {customPresets.length === 0 && (
                    <>
                        <Divider style={{margin: '8px 0'}}/>
                        <Text type="secondary" style={{fontSize: 12}}>暂无自定义预设，点击"新建"或复制内置预设</Text>
                    </>
                )}
            </Card>

            {/* 右侧：预设详情/编辑 */}
            <div style={{flex: 1, overflow: 'auto'}}>
                {!selectedPreset && !isNew ? (
                    <Card style={{height: '100%', display: 'flex', alignItems: 'center', justifyContent: 'center'}}>
                        <Empty description={'选择左侧预设查看详情，或点击"新建"创建新预设'}/>
                    </Card>
                ) : (
                    <>
                        <div style={{display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16}}>
                            <Title level={3} style={{margin: 0}}>
                                {isNew ? '新建预设' : selectedPreset?.name}
                                {selectedPreset?.isCustom && !editing && <Tag color="gold" style={{marginLeft: 8}}>自定义</Tag>}
                                {!selectedPreset?.isCustom && selectedPreset && <Tag style={{marginLeft: 8}}>内置</Tag>}
                            </Title>
                            <Space>
                                {!editing && selectedPreset?.isCustom && (
                                    <Button icon={<EditOutlined/>} onClick={() => setEditing(true)}>编辑</Button>
                                )}
                                {!editing && !selectedPreset?.isCustom && (
                                    <Button icon={<CopyOutlined/>} onClick={() => handleCopy(selectedPreset)}>复制并编辑</Button>
                                )}
                                {editing && (
                                    <Button type="primary" icon={<SaveOutlined/>} onClick={handleSave}>保存预设</Button>
                                )}
                                {!editing && (
                                    <Button type="primary" icon={<PlayCircleOutlined/>} onClick={handleLaunch}>
                                        用此预设启动
                                    </Button>
                                )}
                            </Space>
                        </div>

                        <Form form={form} layout="vertical" disabled={!editing && !isNew}>
                            {/* 服务器基础 */}
                            <Card title="服务器设置" size="small" style={{marginBottom: 16}}>
                                <div style={{display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '0 16px'}}>
                                    <Form.Item name="name" label="服务器名称" rules={[{required: true}]}>
                                        <Input placeholder="我的英灵神殿"/>
                                    </Form.Item>
                                    <Form.Item name="world" label="世界名称">
                                        <Input placeholder="Dedicated"/>
                                    </Form.Item>
                                    <Form.Item name="password" label="默认密码">
                                        <Input.Password placeholder="可留空，启动时再填"/>
                                    </Form.Item>
                                    <Form.Item name="port" label="端口">
                                        <InputNumber min={1} max={65535} style={{width: '100%'}}/>
                                    </Form.Item>
                                    <Form.Item name="public" label="服务器公开" valuePropName="checked">
                                        <Switch/>
                                    </Form.Item>
                                    <Form.Item name="crossplay" label="跨平台" valuePropName="checked">
                                        <Switch/>
                                    </Form.Item>
                                    <Form.Item name="saveDir" label="存档目录">
                                        <Input placeholder="留空使用默认"/>
                                    </Form.Item>
                                    <Form.Item name="logFile" label="日志目录">
                                        <Input placeholder="留空使用默认"/>
                                    </Form.Item>
                                    <Form.Item name="saveInterval" label="自动保存间隔（分钟）">
                                        <InputNumber min={1} max={1440} style={{width: '100%'}}/>
                                    </Form.Item>
                                    <Form.Item name="backupCount" label="保留备份数量">
                                        <InputNumber min={1} max={20} style={{width: '100%'}}/>
                                    </Form.Item>
                                </div>
                            </Card>

                            {/* 世界规则 */}
                            <Card title="世界规则" size="small" style={{marginBottom: 16}}>
                                <Form.Item name="preset" label="难度预设">
                                    <Select options={[
                                        {value: 'normal', label: '正常'}, {value: 'casual', label: '休闲'},
                                        {value: 'easy', label: '简单'}, {value: 'hard', label: '困难'},
                                        {value: 'hardcore', label: '极限'}, {value: 'immersive', label: '沉浸'},
                                        {value: 'hammer', label: '锤子模式'},
                                    ]}/>
                                </Form.Item>
                                <Divider/>
                                <div style={{display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '0 24px'}}>
                                    <Form.Item name="combat" label="战斗难度">
                                        <Radio.Group>
                                            <Space direction="vertical">
                                                {COMBAT_OPTIONS.map(o => <Radio key={o.value} value={o.value}>{o.label}</Radio>)}
                                            </Space>
                                        </Radio.Group>
                                    </Form.Item>
                                    <Form.Item name="deathPenalty" label="死亡惩罚">
                                        <Radio.Group>
                                            <Space direction="vertical">
                                                {DEATH_OPTIONS.map(o => <Radio key={o.value} value={o.value}>{o.label}</Radio>)}
                                            </Space>
                                        </Radio.Group>
                                    </Form.Item>
                                    <Form.Item name="resources" label="资源数量">
                                        <Radio.Group>
                                            <Space direction="vertical">
                                                {RESOURCE_OPTIONS.map(o => <Radio key={o.value} value={o.value}>{o.label}</Radio>)}
                                            </Space>
                                        </Radio.Group>
                                    </Form.Item>
                                    <Form.Item name="raids" label="袭击频率">
                                        <Radio.Group>
                                            <Space direction="vertical">
                                                {RAID_OPTIONS.map(o => <Radio key={o.value} value={o.value}>{o.label}</Radio>)}
                                            </Space>
                                        </Radio.Group>
                                    </Form.Item>
                                    <Form.Item name="portals" label="传送门">
                                        <Radio.Group>
                                            <Space direction="vertical">
                                                {PORTAL_OPTIONS.map(o => <Radio key={o.value} value={o.value}>{o.label}</Radio>)}
                                            </Space>
                                        </Radio.Group>
                                    </Form.Item>
                                </div>
                            </Card>

                            {/* 高级规则 */}
                            <Card title="高级规则" size="small" style={{marginBottom: 16}}>
                                <div style={{display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '0 16px'}}>
                                    <Form.Item name="noBuildCost" label="无建造材料" valuePropName="checked">
                                        <Switch/>
                                    </Form.Item>
                                    <Form.Item name="playerEvents" label="玩家事件" valuePropName="checked">
                                        <Switch/>
                                    </Form.Item>
                                    <Form.Item name="passiveMobs" label="被动生物" valuePropName="checked">
                                        <Switch/>
                                    </Form.Item>
                                    <Form.Item name="noMap" label="无地图" valuePropName="checked">
                                        <Switch/>
                                    </Form.Item>
                                </div>
                            </Card>
                        </Form>
                    </>
                )}
            </div>

            {/* 启动弹窗 */}
            <Modal
                title="选择服务器并启动"
                open={launchModal}
                onOk={handleConfirmLaunch}
                onCancel={() => setLaunchModal(false)}
                okText="启动" cancelText="取消"
            >
                <Space direction="vertical" style={{width: '100%'}} size="middle">
                    <div>
                        <Text>服务器</Text>
                        <Select
                            value={launchServerId || undefined}
                            onChange={setLaunchServerId}
                            style={{width: '100%', marginTop: 4}}
                            placeholder="选择要启动的服务器"
                            options={servers.map(s => ({value: s.id, label: `${s.name} (${s.installPath})`}))}
                        />
                    </div>
                    <div>
                        <Text>密码（留空使用预设默认密码）</Text>
                        <Input.Password
                            value={launchPassword}
                            onChange={e => setLaunchPassword(e.target.value)}
                            style={{marginTop: 4}}
                        />
                    </div>
                </Space>
            </Modal>
        </div>
    )
}

export default Config
