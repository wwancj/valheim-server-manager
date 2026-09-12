import {Card, Typography, Form, Switch, Select, Button, Divider, Space, Tag} from 'antd'
import {SaveOutlined} from '@ant-design/icons'

const {Title, Text} = Typography

function Settings() {
    return (
        <div>
            <Title level={2}>设置</Title>
            <Card title="应用设置">
                <Form layout="vertical">
                    <Form.Item label="深色模式" valuePropName="checked">
                        <Switch defaultChecked/>
                    </Form.Item>
                    <Form.Item label="语言">
                        <Select defaultValue="zh-CN" options={[
                            {value: 'zh-CN', label: '简体中文'},
                            {value: 'en', label: 'English'},
                        ]} style={{width: 200}}/>
                    </Form.Item>
                </Form>
            </Card>
            <Card title="SteamCMD" style={{marginTop: 16}}>
                <Space direction="vertical">
                    <Text>路径：<Tag>自动检测</Tag></Text>
                    <Text type="secondary">SteamCMD 会在首次安装服务器时自动下载</Text>
                </Space>
            </Card>
            <Card title="关于" style={{marginTop: 16}}>
                <Space direction="vertical">
                    <Text strong>Valheim Server Manager</Text>
                    <Text>版本：<Tag color="blue">0.1.0</Tag></Text>
                    <Text type="secondary">用于管理 Valheim Dedicated Server、BepInEx 和 Thunderstore Mod 的桌面工具</Text>
                </Space>
            </Card>
        </div>
    )
}

export default Settings
