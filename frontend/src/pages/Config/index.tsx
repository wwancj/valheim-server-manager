import {useState} from 'react'
import {Card, Typography, Form, Input, InputNumber, Switch, Select, Button, Space, message, Divider} from 'antd'
import {SaveOutlined} from '@ant-design/icons'
import {LoadServerConfig, SaveServerConfig, GetDefaultConfig} from '../../../wailsjs/go/app/App'

const {Title} = Typography

function Config() {
    const [form] = Form.useForm()
    const [serverId, setServerId] = useState('')

    const handleLoad = async () => {
        if (!serverId) {
            message.warning('请输入服务器 ID')
            return
        }
        try {
            const cfg = await LoadServerConfig(serverId)
            form.setFieldsValue(cfg)
            message.success('配置已加载')
        } catch (err: any) {
            message.error(err?.message || '加载失败')
        }
    }

    const handleSave = async () => {
        if (!serverId) {
            message.warning('请输入服务器 ID')
            return
        }
        try {
            const values = form.getFieldsValue()
            await SaveServerConfig(serverId, values)
            message.success('配置已保存')
        } catch (err: any) {
            message.error(err?.message || '保存失败')
        }
    }

    const handleLoadDefault = async () => {
        try {
            const cfg = await GetDefaultConfig()
            form.setFieldsValue(cfg)
            message.success('已加载默认配置')
        } catch (err: any) {
            message.error(err?.message || '加载失败')
        }
    }

    return (
        <div>
            <Title level={2}>服务器配置</Title>
            <Space style={{marginBottom: 16}}>
                <Input
                    placeholder="服务器 ID"
                    value={serverId}
                    onChange={e => setServerId(e.target.value)}
                    style={{width: 300}}
                />
                <Button onClick={handleLoad}>加载配置</Button>
                <Button onClick={handleLoadDefault}>加载默认</Button>
            </Space>

            <Form form={form} layout="vertical" initialValues={{
                serverName: 'Valheim Server', worldName: 'World', port: 2456,
                public: false, preset: 'Normal',
                modifiers: {combat: 'Normal', deathPenalty: 'Normal', resources: 'Normal', raids: 'Normal', portals: 'Normal'},
            }}>
                <Card title="基础配置">
                    <Form.Item name="serverName" label="服务器名称">
                        <Input placeholder="Valheim Server"/>
                    </Form.Item>
                    <Form.Item name="worldName" label="世界名称">
                        <Input placeholder="World"/>
                    </Form.Item>
                    <Form.Item name="password" label="密码">
                        <Input.Password placeholder="服务器密码"/>
                    </Form.Item>
                    <Form.Item name="port" label="端口">
                        <InputNumber min={1} max={65535} style={{width: '100%'}}/>
                    </Form.Item>
                    <Form.Item name="public" label="公开服务器" valuePropName="checked">
                        <Switch/>
                    </Form.Item>
                </Card>

                <Card title="游戏难度" style={{marginTop: 16}}>
                    <Form.Item name="preset" label="难度预设">
                        <Select options={[
                            {value: 'Normal', label: 'Normal'},
                            {value: 'Easy', label: 'Easy'},
                            {value: 'Hard', label: 'Hard'},
                            {value: 'Hardcore', label: 'Hardcore'},
                            {value: 'Casual', label: 'Casual'},
                            {value: 'Immersive', label: 'Immersive'},
                        ]}/>
                    </Form.Item>
                    <Divider>自定义修改器</Divider>
                    {['combat', 'deathPenalty', 'resources', 'raids', 'portals'].map(key => (
                        <Form.Item key={key} name={['modifiers', key]} label={key}>
                            <Select options={[
                                {value: 'Normal', label: 'Normal'},
                                {value: 'Easy', label: 'Easy'},
                                {value: 'Hard', label: 'Hard'},
                                {value: 'VeryHard', label: 'Very Hard'},
                            ]}/>
                        </Form.Item>
                    ))}
                </Card>

                <Button type="primary" icon={<SaveOutlined/>} size="large" onClick={handleSave} style={{marginTop: 16}}>
                    保存配置
                </Button>
            </Form>
        </div>
    )
}

export default Config
