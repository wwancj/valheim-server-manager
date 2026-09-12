import {useState, useEffect} from 'react'
import {Card, Typography, Input, Tabs, List, Button, Tag, Space, message, Switch, Popconfirm, Empty, Avatar} from 'antd'
import {SearchOutlined, DownloadOutlined, DeleteOutlined} from '@ant-design/icons'
import {SearchMods, GetInstalledMods, UninstallMod, SetModEnabled} from '../../../wailsjs/go/app/App'

const {Title, Paragraph, Text} = Typography

function Mods() {
    const [searchResults, setSearchResults] = useState<any[]>([])
    const [installedMods, setInstalledMods] = useState<any[]>([])
    const [searching, setSearching] = useState(false)
    const [serverDir, setServerDir] = useState('')

    // Load popular mods on mount
    useEffect(() => {
        loadPopular()
    }, [])

    const loadPopular = async () => {
        setSearching(true)
        try {
            const results = await SearchMods('', 1)
            setSearchResults(Array.isArray(results) ? results : [])
        } catch (err: any) {
            console.error('Load popular failed:', err)
        } finally {
            setSearching(false)
        }
    }

    const handleSearch = async (query: string) => {
        if (!query.trim()) return
        setSearching(true)
        try {
            const results = await SearchMods(query, 1)
            setSearchResults(Array.isArray(results) ? results : [])
        } catch (err: any) {
            message.error('搜索失败: ' + (err?.message || '未知错误'))
        } finally {
            setSearching(false)
        }
    }

    const loadInstalled = async () => {
        if (!serverDir) {
            message.warning('请先输入服务器路径')
            return
        }
        try {
            const mods = await GetInstalledMods(serverDir)
            setInstalledMods(mods || [])
        } catch {
            message.error('加载失败')
        }
    }

    const handleUninstall = async (fullName: string) => {
        try {
            await UninstallMod(serverDir, fullName)
            message.success('已卸载')
            loadInstalled()
        } catch (err: any) {
            message.error(err?.message || '卸载失败')
        }
    }

    const handleToggle = async (fullName: string, enabled: boolean) => {
        try {
            await SetModEnabled(serverDir, fullName, enabled)
            message.success(enabled ? '已启用' : '已禁用')
            loadInstalled()
        } catch (err: any) {
            message.error(err?.message || '操作失败')
        }
    }

    const getModIcon = (item: any) => {
        const icon = item.icon || item.latest?.icon
        return icon || undefined
    }

    const getModVersion = (item: any) => {
        return item.version_number || item.latest?.version_number || item.version || '-'
    }

    const getModDesc = (item: any) => {
        return item.description || item.latest?.description || '暂无描述'
    }

    const getModDeps = (item: any) => {
        return item.latest?.dependencies || []
    }

    return (
        <div>
            <Title level={2}>Mod 管理</Title>
            <Tabs
                items={[
                    {
                        key: 'store',
                        label: 'Mod 商店',
                        children: (
                            <div>
                                <Input.Search
                                    placeholder="搜索 Mod（例如：CreatureManager、BepInEx）"
                                    enterButton={<><SearchOutlined/> 搜索</>}
                                    size="large"
                                    loading={searching}
                                    onSearch={handleSearch}
                                    style={{marginBottom: 16}}
                                />
                                <List
                                    grid={{gutter: 16, column: 2}}
                                    dataSource={searchResults}
                                    locale={{emptyText: <Empty description="输入关键词搜索 Mod"/>}}
                                    renderItem={(item: any) => (
                                        <List.Item>
                                            <Card
                                                hoverable
                                                title={
                                                    <Space>
                                                        {getModIcon(item) && (
                                                            <Avatar src={getModIcon(item)} size={32} shape="square"/>
                                                        )}
                                                        <Text strong>{item.name}</Text>
                                                        <Tag>v{getModVersion(item)}</Tag>
                                                    </Space>
                                                }
                                                extra={<Tag color="blue">{item.owner}</Tag>}
                                            >
                                                <Paragraph ellipsis={{rows: 2}}>{getModDesc(item)}</Paragraph>
                                                {getModDeps(item).length > 0 && (
                                                    <div style={{marginBottom: 8}}>
                                                        <Text type="secondary" style={{fontSize: 12}}>依赖: </Text>
                                                        {getModDeps(item).slice(0, 3).map((dep: string) => (
                                                            <Tag key={dep} style={{fontSize: 11}}>{dep.split('-').pop()}</Tag>
                                                        ))}
                                                    </div>
                                                )}
                                                <Space>
                                                    <Text type="secondary">⬇ {(item.total_downloads || item.downloads || 0).toLocaleString()}</Text>
                                                    <Button
                                                        type="primary"
                                                        icon={<DownloadOutlined/>}
                                                        size="small"
                                                        onClick={() => message.info('请先在服务器页面选择一个已安装的服务器')}
                                                    >
                                                        安装
                                                    </Button>
                                                </Space>
                                            </Card>
                                        </List.Item>
                                    )}
                                />
                            </div>
                        ),
                    },
                    {
                        key: 'installed',
                        label: '已安装',
                        children: (
                            <div>
                                <Space style={{marginBottom: 16}}>
                                    <Input
                                        placeholder="服务器安装路径（例如 D:\ValheimServer）"
                                        value={serverDir}
                                        onChange={e => setServerDir(e.target.value)}
                                        style={{width: 400}}
                                    />
                                    <Button onClick={loadInstalled}>加载</Button>
                                </Space>
                                <List
                                    dataSource={installedMods}
                                    locale={{emptyText: <Empty description="暂无已安装的 Mod"/>}}
                                    renderItem={(mod: any) => (
                                        <List.Item
                                            actions={[
                                                <Switch
                                                    checked={mod.enabled}
                                                    checkedChildren="启用"
                                                    unCheckedChildren="禁用"
                                                    onChange={(checked) => handleToggle(mod.fullName, checked)}
                                                />,
                                                <Popconfirm title="确定卸载？" onConfirm={() => handleUninstall(mod.fullName)}>
                                                    <Button danger icon={<DeleteOutlined/>} size="small">卸载</Button>
                                                </Popconfirm>,
                                            ]}
                                        >
                                            <List.Item.Meta
                                                title={<Space>
                                                    <Text strong>{mod.name}</Text>
                                                    <Tag>v{mod.version}</Tag>
                                                    {mod.enabled ? <Tag color="green">启用</Tag> : <Tag color="red">禁用</Tag>}
                                                </Space>}
                                                description={mod.description || mod.fullName}
                                            />
                                        </List.Item>
                                    )}
                                />
                            </div>
                        ),
                    },
                ]}
            />
        </div>
    )
}

export default Mods
