import {useState, useEffect} from 'react'
import {Card, Col, Row, Statistic, Typography, Progress, Space, Tag} from 'antd'
import {CloudServerOutlined, AppstoreOutlined, SaveOutlined, DashboardOutlined, DatabaseOutlined} from '@ant-design/icons'
import {GetAppInfo, GetSystemInfo, GetAppUptime, IsServerRunning} from '../../api/client'

const {Title, Text} = Typography

function Dashboard() {
    const [appInfo, setAppInfo] = useState<Record<string, string>>({})
    const [sysInfo, setSysInfo] = useState<any>({})
    const [uptime, setUptime] = useState(0)
    const [serverRunning, setServerRunning] = useState(false)

    useEffect(() => {
        GetAppInfo().then(setAppInfo)
        GetSystemInfo().then(setSysInfo)
        GetAppUptime().then(setUptime)
        IsServerRunning().then(setServerRunning)

        const timer = setInterval(() => {
            GetSystemInfo().then(setSysInfo)
            GetAppUptime().then(setUptime)
            IsServerRunning().then(setServerRunning)
        }, 3000)

        return () => clearInterval(timer)
    }, [])

    const formatBytes = (bytes: number) => {
        if (!bytes) return '0 B'
        const gb = bytes / (1024 * 1024 * 1024)
        return gb.toFixed(1) + ' GB'
    }

    const formatUptime = (seconds: number) => {
        const h = Math.floor(seconds / 3600)
        const m = Math.floor((seconds % 3600) / 60)
        return `${h}h ${m}m`
    }

    return (
        <div>
            <Title level={2}>Dashboard</Title>
            <Row gutter={[16, 16]}>
                <Col span={6}>
                    <Card>
                        <Statistic
                            title="服务器状态"
                            value={serverRunning ? '运行中' : '未运行'}
                            prefix={<CloudServerOutlined/>}
                            valueStyle={{color: serverRunning ? '#52c41a' : '#8c8c8c'}}
                        />
                    </Card>
                </Col>
                <Col span={6}>
                    <Card>
                        <Statistic
                            title="CPU 使用率"
                            value={sysInfo.cpuPercent?.toFixed(1) || 0}
                            suffix="%"
                            prefix={<DashboardOutlined/>}
                        />
                        <Progress percent={sysInfo.cpuPercent || 0} showInfo={false} status="active"/>
                    </Card>
                </Col>
                <Col span={6}>
                    <Card>
                        <Statistic
                            title="内存使用"
                            value={formatBytes(sysInfo.memUsed)}
                            prefix={<DatabaseOutlined/>}
                        />
                        <Progress percent={sysInfo.memPercent || 0} showInfo={false} status="active"/>
                        <Text type="secondary" style={{fontSize: 12}}>共 {formatBytes(sysInfo.memTotal)}</Text>
                    </Card>
                </Col>
                <Col span={6}>
                    <Card>
                        <Statistic
                            title="运行时间"
                            value={formatUptime(uptime)}
                        />
                    </Card>
                </Col>
            </Row>
            <Card style={{marginTop: 16}}>
                <Title level={4}>应用信息</Title>
                <Space direction="vertical">
                    <Text>名称：{appInfo.name || '-'}</Text>
                    <Text>版本：<Tag color="blue">{appInfo.version || '-'}</Tag></Text>
                    <Text>Go：{appInfo.go || '-'}</Text>
                    <Text>系统：{appInfo.os || '-'} / {appInfo.arch || '-'}</Text>
                </Space>
            </Card>
        </div>
    )
}

export default Dashboard
